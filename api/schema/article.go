package schema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ArticleBlock struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Md         string `json:"md,omitempty"`
	Collection string `json:"collection,omitempty"`
	Filename   string `json:"filename,omitempty"`
	Caption    string `json:"caption,omitempty"`
	V          string `json:"v,omitempty"`
	Title      string `json:"title,omitempty"`
	Artist     string `json:"artist,omitempty"`
	ScenePath  string `json:"scene_path,omitempty"`
}

type ArticleBlockEnvelope struct {
	Version int            `json:"version"`
	Blocks  []ArticleBlock `json:"blocks"`
}

func (e ArticleBlockEnvelope) MarshalJSON() ([]byte, error) {
	type alias ArticleBlockEnvelope

	normalized := alias(e)
	if normalized.Blocks == nil {
		normalized.Blocks = []ArticleBlock{}
	}

	return json.Marshal(normalized)
}

func (e ArticleBlockEnvelope) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *ArticleBlockEnvelope) Scan(value interface{}) error {
	if value == nil {
		*e = ArticleBlockEnvelope{Version: 1, Blocks: []ArticleBlock{}}
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return fmt.Errorf("unsupported article block value type %T", value)
	}

	return json.Unmarshal(data, e)
}

type Article struct {
	BaseColumns

	Uri         string               `gorm:"uniqueIndex" json:"uri"`
	AuthorID    uuid.UUID            `gorm:",notnull" json:"-"`
	Author      Author               `json:"author"`
	Title       string               `json:"title"`
	SubTitle    string               `json:"subtitle"`
	Blocks      ArticleBlockEnvelope `gorm:"type:text" json:"blocks"`
	Summary     string               `json:"summary"`
	Cover       string               `json:"cover"`
	PublishedAt *time.Time           `gorm:"default:null" json:"published_at"`
}
