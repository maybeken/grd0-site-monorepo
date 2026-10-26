package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"grd0.net/api/schema"
	"grd0.net/api/storage"
	"grd0.net/api/utils"
)

const (
	ARTICLE_URI_PATTERN            = `^[a-z0-9-]+$`
	ARTICLE_SUMMARY_MAX            = 300
	DEFAULT_ARTICLE_SCENE_MAX_BYTE = 5 * 1024 * 1024
)

var (
	articleUriRegex    = regexp.MustCompile(ARTICLE_URI_PATTERN)
	markdownImageRegex = regexp.MustCompile(`!\[[^\]]*\]\(\s*<?([^)\s>]+)>?`)
)

type ArticleRequest struct {
	Uri                string                `json:"uri"`
	Title              string                `json:"title"`
	SubTitle           string                `json:"subtitle"`
	PublishedAt        *time.Time            `json:"published_at"`
	PublishedAtPresent bool                  `json:"-"`
	Blocks             []schema.ArticleBlock `json:"blocks"`
}

func (r *ArticleRequest) UnmarshalJSON(data []byte) error {
	type alias ArticleRequest

	aux := struct {
		*alias
		PublishedAt json.RawMessage `json:"published_at"`
	}{alias: (*alias)(r)}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if len(aux.PublishedAt) == 0 {
		return nil
	}

	r.PublishedAtPresent = true

	if string(aux.PublishedAt) == "null" {
		r.PublishedAt = nil
		return nil
	}

	var published_at time.Time
	if err := json.Unmarshal(aux.PublishedAt, &published_at); err != nil {
		return err
	}

	r.PublishedAt = &published_at

	return nil
}

type ArticleResponse struct {
	Uri         string                `json:"uri"`
	Author      schema.Author         `json:"author"`
	Title       string                `json:"title"`
	SubTitle    string                `json:"subtitle"`
	Blocks      []schema.ArticleBlock `json:"blocks"`
	Summary     string                `json:"summary"`
	Cover       string                `json:"cover"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	PublishedAt *time.Time            `json:"published_at"`
}

type ArticleListEntry struct {
	Uri         string        `json:"uri"`
	Author      schema.Author `json:"author"`
	Title       string        `json:"title"`
	SubTitle    string        `json:"subtitle"`
	Summary     string        `json:"summary"`
	Cover       string        `json:"cover"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	PublishedAt *time.Time    `json:"published_at"`
}

type ArticleSceneResponse struct {
	ScenePath string `json:"scene_path"`
}

func (h *Handler) resolveArticleAuthor(c echo.Context) (*schema.Author, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok || token == nil {
		return nil, fmt.Errorf("missing authentication token")
	}

	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok || claims == nil || claims.Email == "" {
		return nil, fmt.Errorf("missing author claims")
	}

	display_name := claims.Name
	if display_name == "" {
		display_name = claims.Email
	}

	db := h.DB

	var author schema.Author
	err := db.Where("email = ?", claims.Email).Take(&author).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		author = schema.Author{
			Email:       claims.Email,
			DisplayName: display_name,
		}

		if err := db.Omit(clause.Associations).Create(&author).Error; err != nil {
			return nil, err
		}

		return &author, nil
	} else if err != nil {
		return nil, err
	}

	if author.DisplayName != display_name {
		if err := db.Model(&schema.Author{}).Where("id = ?", author.ID).Update("display_name", display_name).Error; err != nil {
			return nil, err
		}
		author.DisplayName = display_name
	}

	return &author, nil
}

func validArticleUri(uri string) bool {
	return articleUriRegex.MatchString(uri)
}

func generateArticleUri(db *gorm.DB) (string, error) {
	for attempt := 0; attempt < 5; attempt++ {
		candidate := "untitled-" + uuid.NewString()[:8]

		var existing int64
		if err := db.Model(&schema.Article{}).Where("uri = ?", candidate).Count(&existing).Error; err != nil {
			return "", err
		}
		if existing == 0 {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("unable to allocate a temporary article URI")
}

func validArticleScenePath(scene_path string) bool {
	if !strings.HasPrefix(scene_path, "scenes/") {
		return false
	}

	if strings.Contains(scene_path, "\\") || strings.HasPrefix(scene_path, "/") {
		return false
	}

	for _, part := range strings.Split(scene_path, "/") {
		if part == ".." {
			return false
		}
	}

	return path.Clean(scene_path) == scene_path
}

func validateArticleBlocks(blocks []schema.ArticleBlock) error {
	seen := make(map[string]bool, len(blocks))

	for idx, block := range blocks {
		if strings.TrimSpace(block.ID) == "" {
			return fmt.Errorf("block %d: id is required", idx)
		}

		if seen[block.ID] {
			return fmt.Errorf("block %d: duplicate id %q", idx, block.ID)
		}
		seen[block.ID] = true

		switch block.Type {
		case "markdown":
		case "image":
			if strings.TrimSpace(block.Collection) == "" || strings.TrimSpace(block.Filename) == "" {
				return fmt.Errorf("block %d: image requires collection and filename", idx)
			}
		case "music":
			if strings.TrimSpace(block.V) == "" || strings.TrimSpace(block.Title) == "" {
				return fmt.Errorf("block %d: music requires v and title", idx)
			}
		case "excalidraw":
			if strings.TrimSpace(block.ScenePath) == "" {
				return fmt.Errorf("block %d: excalidraw requires scene_path", idx)
			}
			if !validArticleScenePath(block.ScenePath) {
				return fmt.Errorf("block %d: invalid scene_path", idx)
			}
		default:
			return fmt.Errorf("block %d: unknown type %q", idx, block.Type)
		}
	}

	return nil
}

func deriveArticleSummary(blocks []schema.ArticleBlock) string {
	for _, block := range blocks {
		if block.Type != "markdown" {
			continue
		}

		summary := strings.Join(strings.Fields(block.Md), " ")
		if summary == "" {
			return ""
		}

		if utf8.RuneCountInString(summary) > ARTICLE_SUMMARY_MAX {
			runes := []rune(summary)
			summary = string(runes[:ARTICLE_SUMMARY_MAX]) + "..."
		}

		return summary
	}

	return ""
}

func deriveArticleCover(blocks []schema.ArticleBlock) string {
	for _, block := range blocks {
		switch block.Type {
		case "image":
			return strings.TrimRight(block.Collection, "/") + "/" + block.Filename
		case "markdown":
			if cover := firstMarkdownImage(block.Md); cover != "" {
				return cover
			}
		}
	}

	return ""
}

func firstMarkdownImage(md string) string {
	for _, match := range markdownImageRegex.FindAllStringSubmatch(md, -1) {
		src := match[1]
		if src == "" {
			continue
		}

		return normalizeMarkdownImageSrc(src)
	}

	return ""
}

func normalizeMarkdownImageSrc(src string) string {
	if !strings.Contains(src, "//") {
		return src
	}

	parsed, err := url.Parse(src)
	if err != nil {
		return src
	}

	image_path := parsed.Path
	if idx := strings.Index(image_path, "/cdn-cgi/image/"); idx >= 0 {
		rest := image_path[idx+len("/cdn-cgi/image/"):]
		if slash := strings.Index(rest, "/"); slash >= 0 {
			return rest[slash:]
		}
	}

	return src
}

func articleBlocks(article *schema.Article) []schema.ArticleBlock {
	if article.Blocks.Blocks == nil {
		return []schema.ArticleBlock{}
	}

	return article.Blocks.Blocks
}

func articleResponse(article *schema.Article) ArticleResponse {
	return ArticleResponse{
		Uri:         article.Uri,
		Author:      article.Author,
		Title:       article.Title,
		SubTitle:    article.SubTitle,
		Blocks:      articleBlocks(article),
		Summary:     article.Summary,
		Cover:       article.Cover,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
		PublishedAt: article.PublishedAt,
	}
}

func articleListEntry(article *schema.Article) ArticleListEntry {
	cover := article.Cover
	if cover == "" {
		cover = deriveArticleCover(articleBlocks(article))
	}

	return ArticleListEntry{
		Uri:         article.Uri,
		Author:      article.Author,
		Title:       article.Title,
		SubTitle:    article.SubTitle,
		Summary:     article.Summary,
		Cover:       cover,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
		PublishedAt: article.PublishedAt,
	}
}

func (h *Handler) findArticle(uri string) (*schema.Article, error) {
	var article schema.Article
	err := h.DB.Preload("Author").Where("uri = ?", uri).Take(&article).Error

	return &article, err
}

func sceneKey(article_id uuid.UUID, scene_id string) string {
	return fmt.Sprintf("scenes/%s/%s.json", article_id.String(), scene_id)
}

func removedScenes(old_blocks, new_blocks []schema.ArticleBlock) []string {
	kept := make(map[string]bool, len(new_blocks))
	for _, block := range new_blocks {
		if block.Type == "excalidraw" {
			kept[block.ScenePath] = true
		}
	}

	removed := []string{}
	seen := make(map[string]bool, len(old_blocks))
	for _, block := range old_blocks {
		if block.Type != "excalidraw" || block.ScenePath == "" {
			continue
		}
		if kept[block.ScenePath] || seen[block.ScenePath] {
			continue
		}
		seen[block.ScenePath] = true
		removed = append(removed, block.ScenePath)
	}

	return removed
}

func (h *Handler) GetArticle(c echo.Context) error {
	uri := c.Param("uri")

	if uri != "" {
		article, err := h.findArticle(uri)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "")
		} else if err != nil {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, articleResponse(article))
	}

	var articles []schema.Article
	if err := h.DB.
		Select("id", "created_at", "updated_at", "uri", "author_id", "title", "sub_title", "summary", "cover", "blocks", "published_at").
		Where("published_at IS NOT NULL AND published_at <= ?", time.Now()).
		Preload("Author").
		Order("published_at DESC").
		Find(&articles).Error; err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	entries := make([]ArticleListEntry, 0, len(articles))
	for idx := range articles {
		entries = append(entries, articleListEntry(&articles[idx]))
	}

	return c.JSON(http.StatusOK, entries)
}

func (h *Handler) GetArticleDraft(c echo.Context) error {
	uri := c.Param("uri")

	article, err := h.findArticle(uri)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, "")
	} else if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, articleResponse(article))
}

func (h *Handler) GetArticleByAdmin(c echo.Context) error {
	page, pageSize := parsePagination(c)

	var total int64
	if err := h.DB.Model(&schema.Article{}).Count(&total).Error; err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	var articles []schema.Article
	if err := h.DB.
		Select("id", "created_at", "updated_at", "uri", "author_id", "title", "sub_title", "summary", "cover", "blocks", "published_at").
		Preload("Author").
		Order("updated_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&articles).Error; err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	entries := make([]ArticleListEntry, 0, len(articles))
	for idx := range articles {
		entries = append(entries, articleListEntry(&articles[idx]))
	}

	return c.JSON(http.StatusOK, schema.PaginatedResponse[ArticleListEntry]{
		Data:       entries,
		Total:      int(total),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(total) / float64(pageSize))),
	})
}

func (h *Handler) CreateArticle(c echo.Context) error {
	var req ArticleRequest
	if err := c.Bind(&req); err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Unable to parse input.")
	}

	db := h.DB

	requested_uri := strings.TrimSpace(req.Uri)
	uri := requested_uri
	if uri == "" {
		generated, err := generateArticleUri(db)
		if err != nil {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
		}
		uri = generated
	} else if !validArticleUri(uri) {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "URI must match [a-z0-9-]+.")
	}

	author, err := h.resolveArticleAuthor(c)
	if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusUnauthorized, err.Error())
	}

	var existing int64
	if err := db.Model(&schema.Article{}).Where("uri = ?", uri).Count(&existing).Error; err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}
	if existing > 0 {
		if requested_uri == "" {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, "Could not allocate a temporary article URI.")
		}
		return h.ErrorResponseConstructor(c, http.StatusConflict, "An article with this URI already exists.")
	}

	article := schema.Article{
		Uri:         uri,
		AuthorID:    author.ID,
		Author:      *author,
		Title:       req.Title,
		SubTitle:    req.SubTitle,
		Blocks:      schema.ArticleBlockEnvelope{Version: 1, Blocks: []schema.ArticleBlock{}},
		PublishedAt: nil,
	}

	if err := db.Omit(clause.Associations).Create(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			if requested_uri == "" {
				return h.ErrorResponseConstructor(c, http.StatusInternalServerError, "Could not allocate a temporary article URI.")
			}
			return h.ErrorResponseConstructor(c, http.StatusConflict, "An article with this URI already exists.")
		}
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, articleResponse(&article))
}

func (h *Handler) UpdateArticle(c echo.Context) error {
	uri := c.Param("uri")
	if !validArticleUri(uri) {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "URI must match [a-z0-9-]+.")
	}

	var req ArticleRequest
	if err := c.Bind(&req); err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Unable to parse input.")
	}

	article, err := h.findArticle(uri)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, "")
	} else if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	new_uri := strings.TrimSpace(req.Uri)
	if new_uri == "" {
		new_uri = uri
	}
	if !validArticleUri(new_uri) {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "URI must match [a-z0-9-]+.")
	}

	old_blocks := articleBlocks(article)

	blocks := req.Blocks
	if blocks == nil {
		blocks = old_blocks
	}
	if err := validateArticleBlocks(blocks); err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, err.Error())
	}

	author, err := h.resolveArticleAuthor(c)
	if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusUnauthorized, err.Error())
	}

	db := h.DB

	if new_uri != uri {
		var existing int64
		if err := db.Model(&schema.Article{}).
			Where("uri = ? AND id <> ?", new_uri, article.ID).
			Count(&existing).Error; err != nil {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
		}
		if existing > 0 {
			return h.ErrorResponseConstructor(c, http.StatusConflict, "An article with this URI already exists.")
		}
	}

	published_at := article.PublishedAt
	if req.PublishedAtPresent {
		published_at = req.PublishedAt
		if published_at != nil && published_at.IsZero() {
			published_at = nil
		}
	}

	envelope := schema.ArticleBlockEnvelope{Version: 1, Blocks: blocks}
	updates := map[string]interface{}{
		"uri":          new_uri,
		"author_id":    author.ID,
		"title":        req.Title,
		"sub_title":    req.SubTitle,
		"blocks":       envelope,
		"summary":      deriveArticleSummary(blocks),
		"cover":        deriveArticleCover(blocks),
		"published_at": published_at,
	}

	if err := db.Model(&schema.Article{}).Where("id = ?", article.ID).Updates(updates).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return h.ErrorResponseConstructor(c, http.StatusConflict, "An article with this URI already exists.")
		}
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	for _, scene_path := range removedScenes(old_blocks, blocks) {
		if err := h.Blobs.Delete(c.Request().Context(), scene_path); err != nil {
			h.Logger.Warnf("Failed to delete orphaned scene %s: %v", scene_path, err)
		}
	}

	updated, err := h.findArticle(new_uri)
	if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, articleResponse(updated))
}

func (h *Handler) DeleteArticle(c echo.Context) error {
	uri := c.Param("uri")

	res := h.DB.Model(&schema.Article{}).Where("uri = ?", uri).Update("published_at", nil)
	if res.Error != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, res.Error.Error())
	}
	if res.RowsAffected <= 0 {
		return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
	}

	return c.JSON(http.StatusOK, uri)
}

func (h *Handler) PurgeArticle(c echo.Context) error {
	uri := c.Param("uri")

	article, err := h.findArticle(uri)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, "")
	} else if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	prefix := fmt.Sprintf("scenes/%s", article.ID.String())
	if err := h.Blobs.DeletePrefix(c.Request().Context(), prefix); err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	if err := h.DB.Unscoped().Delete(&article).Error; err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, uri)
}

func (h *Handler) sceneLimit() int64 {
	if h.ArticleSceneMaxBytes > 0 {
		return h.ArticleSceneMaxBytes
	}

	return DEFAULT_ARTICLE_SCENE_MAX_BYTE
}

func (h *Handler) PutArticleScene(c echo.Context) error {
	uri := c.Param("uri")
	scene_id, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Scene identifier must be a UUID.")
	}

	article, err := h.findArticle(uri)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, "")
	} else if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	limit := h.sceneLimit()
	if c.Request().ContentLength > limit {
		return h.ErrorResponseConstructor(c, http.StatusRequestEntityTooLarge,
			fmt.Sprintf("Scene exceeds the %d byte limit.", limit))
	}

	body := http.MaxBytesReader(c.Response(), c.Request().Body, limit)
	key := sceneKey(article.ID, scene_id.String())

	if err := h.Blobs.Put(c.Request().Context(), key, body, "application/json"); err != nil {
		var max_err *http.MaxBytesError
		if errors.As(err, &max_err) {
			return h.ErrorResponseConstructor(c, http.StatusRequestEntityTooLarge,
				fmt.Sprintf("Scene exceeds the %d byte limit.", limit))
		}
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, ArticleSceneResponse{ScenePath: key})
}

func (h *Handler) serveArticleScene(c echo.Context) error {
	uri := c.Param("uri")
	scene_id, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Scene identifier must be a UUID.")
	}

	article, err := h.findArticle(uri)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, "")
	} else if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	reader, err := h.Blobs.Get(c.Request().Context(), sceneKey(article.ID, scene_id.String()))
	if errors.Is(err, storage.ErrBlobNotFound) {
		return c.JSON(http.StatusNotFound, "")
	} else if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}
	defer reader.Close()

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)

	if _, err := io.Copy(c.Response(), reader); err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetArticleScene(c echo.Context) error {
	return h.serveArticleScene(c)
}

func (h *Handler) GetArticleSceneDraft(c echo.Context) error {
	return h.serveArticleScene(c)
}

func (h *Handler) DeleteArticleScene(c echo.Context) error {
	uri := c.Param("uri")
	scene_id, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Scene identifier must be a UUID.")
	}

	article, err := h.findArticle(uri)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, "")
	} else if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	key := sceneKey(article.ID, scene_id.String())
	if err := h.Blobs.Delete(c.Request().Context(), key); err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, ArticleSceneResponse{ScenePath: key})
}

func (h *Handler) GenerateArticlePresignedUrl(c echo.Context) error {
	key := c.Param("key")
	if key == "" || strings.Contains(key, "/") || strings.Contains(key, "\\") || strings.Contains(key, "..") {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Invalid attachment key.")
	}

	if h.S3 == nil {
		return h.ErrorResponseConstructor(c, http.StatusNotImplemented, "Attachment uploads are unavailable in local mode.")
	}

	s3 := h.S3
	year := time.Now().UTC().Year()
	month := time.Now().UTC().Format("01")
	object_path := fmt.Sprintf("article/%d/%s/%s", year, month, key)

	presigned_url, err := s3.GeneratePresignedUrl(
		storage.S3PUT,
		utils.GetEnvWithFallback("S3_BUCKET", "assets-grd0-net"),
		object_path,
		60,
	)

	if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusAccepted, PresignedUrlResponseBody{
		Url:       presigned_url.String(),
		ObjectKey: url.PathEscape(object_path),
	})
}
