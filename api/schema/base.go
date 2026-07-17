package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseColumns struct {
	ID        uuid.UUID      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *BaseColumns) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()

	return
}
