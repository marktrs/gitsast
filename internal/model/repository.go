package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Repository struct {
	bun.BaseModel `bun:"table:repositories,alias:repo"`

	ID        string    `json:"id" bun:",pk"`
	Name      string    `json:"name" bun:",notnull"`
	RemoteURL string    `json:"remote_url" bun:",notnull"`
	CreatedAt time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`

	Report *Report `json:"report,omitempty" bun:"rel:has-one,join:id=repository_id"`
}

// IRepository defines methods for read/write repositories table.
type IRepository interface {
	GetById(ctx context.Context, id string) (*Repository, error)
	List(ctx context.Context, f *RepositoryFilter) ([]*Repository, error)
	Add(ctx context.Context, repo *Repository) (*Repository, error)
	Update(ctx context.Context, id string, repo map[string]interface{}) error
	Remove(ctx context.Context, id string) error
}

type repository struct {
	db *bun.DB
}

func NewRepositoryRepo(db *bun.DB) IRepository {
	return &repository{db}
}

func (r *repository) GetById(ctx context.Context, id string) (*Repository, error) {
	repo := &Repository{}
	err := r.db.NewSelect().Model(repo).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *repository) List(ctx context.Context, f *RepositoryFilter) ([]*Repository, error) {
	repos := []*Repository{}
	err := r.db.NewSelect().
		Model(&repos).
		Apply(f.query).
		Limit(f.Limit).
		Offset(f.Offset).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func (r *repository) Add(ctx context.Context, repo *Repository) (*Repository, error) {
	_, err := r.db.NewInsert().Model(repo).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return repo, nil

}

func (r *repository) Update(ctx context.Context, id string, repo map[string]interface{}) error {
	_, err := r.db.NewUpdate().
		Model(&repo).
		TableExpr("repositories").
		Where("id = ?", id).
		OmitZero().
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Remove(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().
		Model((*Repository)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
