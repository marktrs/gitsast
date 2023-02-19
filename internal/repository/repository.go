package repository

import (
	"time"

	"github.com/uptrace/bun"
)

type Repository struct {
	bun.BaseModel `bun:"table:repositories,alias:r"`

	ID        string    `json:"id" bun:",pk,type:uuid"`
	Name      string    `json:"name" bun:",notnull"`
	RemoteURL string    `json:"remote_url" bun:",notnull"`
	CreatedAt time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
}
