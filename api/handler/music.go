package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"grd0.net/api/schema"
)

func orderSongs(db *gorm.DB) *gorm.DB {
	return db.Order("sorting DESC").Order("created_at ASC")
}

func (h *Handler) GetMusic(c echo.Context) error {
	db := h.DB

	var music []schema.Music
	db.Order("sorting DESC").Order("created_at ASC").Find(&music)

	return c.JSON(http.StatusOK, music)
}

func (h *Handler) UpsertMusic(c echo.Context) error {
	db := h.DB

	// var music schema.Music
	var music schema.Music
	if err := c.Bind(&music); err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Unable to parse input.")
	}

	if err := db.Create(&music).Error; errors.Is(err, gorm.ErrDuplicatedKey) {
		if err := db.Unscoped().Model(schema.Music{}).Where("v = ?", music.V).Updates(&music).Update("deleted_at", nil).Error; err != nil {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
		}
	} else if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, music)
}

func (h *Handler) DeleteMusic(c echo.Context) error {
	vid := c.Param("v")

	db := h.DB

	res := db.Where("v = ?", vid).Delete(&schema.Music{})
	row_count := res.RowsAffected
	err := res.Error

	if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	} else if row_count <= 0 {
		return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
	}

	return c.JSON(http.StatusAccepted, nil)
}

func (h *Handler) GetPlaylists(c echo.Context) error {
	db := h.DB

	var playlists []schema.Playlist
	if err := db.
		Preload("Songs", orderSongs).
		Order("sorting DESC").
		Order("created_at ASC").
		Find(&playlists).Error; err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, playlists)
}

func replacePlaylistSongs(tx *gorm.DB, playlist_id uuid.UUID, songs []schema.Music) error {
	ids := make([]uuid.UUID, 0, len(songs))
	seen := make(map[uuid.UUID]bool, len(songs))

	for _, song := range songs {
		if song.ID == uuid.Nil || seen[song.ID] {
			continue
		}

		seen[song.ID] = true
		ids = append(ids, song.ID)
	}

	if err := tx.Where("playlist_id = ?", playlist_id).Delete(&schema.PlaylistSong{}).Error; err != nil {
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	rows := make([]schema.PlaylistSong, 0, len(ids))
	for _, id := range ids {
		rows = append(rows, schema.PlaylistSong{PlaylistID: playlist_id, MusicID: id})
	}

	return tx.Create(&rows).Error
}

func (h *Handler) UpsertPlaylist(c echo.Context) error {
	db := h.DB

	var playlist schema.Playlist
	if err := c.Bind(&playlist); err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Unable to parse input.")
	}

	playlist.Title = strings.TrimSpace(playlist.Title)
	if playlist.Title == "" {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Playlist title is required.")
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if playlist.ID != uuid.Nil {
			res := tx.Omit(clause.Associations).Model(&schema.Playlist{}).Where("id = ?", playlist.ID).Select("title", "description", "sorting").Updates(&playlist)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		} else {
			if err := tx.Omit(clause.Associations).Create(&playlist).Error; err != nil {
				return err
			}
		}

		return replacePlaylistSongs(tx, playlist.ID, playlist.Songs)
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
	} else if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	if err := db.Preload("Songs", orderSongs).First(&playlist, "id = ?", playlist.ID).Error; err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, playlist)
}

func (h *Handler) DeletePlaylist(c echo.Context) error {
	id := c.Param("id")

	db := h.DB

	err := db.Transaction(func(tx *gorm.DB) error {
		var playlist schema.Playlist
		if err := tx.First(&playlist, "id = ?", id).Error; err != nil {
			return err
		}

		if err := tx.Where("playlist_id = ?", playlist.ID).Delete(&schema.PlaylistSong{}).Error; err != nil {
			return err
		}

		return tx.Delete(&playlist).Error
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
	} else if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusAccepted, nil)
}
