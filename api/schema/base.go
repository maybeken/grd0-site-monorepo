package schema

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type BaseColumns struct {
	ID        string    `bun:",type:uuid,pk" json:"id"`
	CreatedAt time.Time `bun:"$created_at,nullzero,notnull" json:"created_at"`
	UpdatedAt time.Time `bun:"$updated_at,nullzero,notnull" json:"updated_at"`
}

func (u *BaseColumns) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		u.ID = uuid.NewString()
		u.CreatedAt = time.Now().UTC()
		u.UpdatedAt = time.Now().UTC()
	case *bun.UpdateQuery:
		u.UpdatedAt = time.Now().UTC()
	}
	return nil
}
