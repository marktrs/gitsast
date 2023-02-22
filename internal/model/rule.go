package model

import (
	"context"
	"fmt"
	"time"

	"github.com/marktrs/gitsast/app"
	"github.com/uptrace/bun"
)

type Rule struct {
	bun.BaseModel `bun:"table:rules,alias:rule"`

	ID          uint64    `json:"id" bun:",pk,autoincrement,notnull"`
	Name        string    `json:"name" bun:",unique,notnull"`
	Keyword     string    `json:"keyword" bun:",unique,notnull"`
	Description string    `json:"description" bun:",notnull"`
	Severity    Score     `json:"severity" bun:",notnull"`
	CreatedAt   time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
}

// Score type used by severity values
type Score int

const (
	// Low severity
	Low Score = iota + 1
	// Medium severity
	Medium
	// High severity
	High
)

// String converts a Score into a string
func (c Score) String() string {
	switch c {
	case High:
		return "HIGH"
	case Medium:
		return "MEDIUM"
	case Low:
		return "LOW"
	}
	return "UNDEFINED"
}

// GetFormattedRuleId - return a formatted rule ID
func GetFormattedRuleId(id uint64) string {
	return "G" + fmt.Sprintf("%03d", id)
}

// IRuleRepo - interface for rules repository
type IRuleRepo interface {
	GetAll(ctx context.Context) ([]*Rule, error)
	GetByID(ctx context.Context, id uint64) (*Rule, error)
	GetByKeyword(ctx context.Context, keyword string) (*Rule, error)
	Create(ctx context.Context, rule *Rule) error
	Update(ctx context.Context, rule *Rule) error
	Delete(ctx context.Context, id uint64) error
}

type RuleRepo struct {
	app *app.App
}

// NewRuleRepo - create a new rules repository instance
func NewRuleRepo(app *app.App) IRuleRepo {
	return &RuleRepo{app}
}

// GetAll - get all rules
func (r *RuleRepo) GetAll(ctx context.Context) ([]*Rule, error) {
	var rules []*Rule
	if err := r.app.DB().NewSelect().Model(&rules).Scan(ctx); err != nil {
		return nil, err
	}
	return rules, nil
}

// GetByID - get a rule by ID
func (r *RuleRepo) GetByID(ctx context.Context, id uint64) (*Rule, error) {
	var rule Rule
	if err := r.app.DB().NewSelect().Model(&rule).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return &rule, nil
}

// GetByKeyword - get a rule by keyword
func (r *RuleRepo) GetByKeyword(ctx context.Context, keyword string) (*Rule, error) {
	var rule Rule
	if err := r.app.DB().NewSelect().Model(&rule).Where("keyword = ?", keyword).Scan(ctx); err != nil {
		return nil, err
	}
	return &rule, nil
}

// Create - create a new rule
func (r *RuleRepo) Create(ctx context.Context, rule *Rule) error {
	if _, err := r.app.DB().NewInsert().Model(rule).Exec(ctx); err != nil {
		return err
	}
	return nil
}

// Update - update a rule
func (r *RuleRepo) Update(ctx context.Context, rule *Rule) error {
	if _, err := r.app.DB().NewUpdate().Model(rule).Exec(ctx); err != nil {
		return err
	}
	return nil
}

// Delete - delete a rule
func (r *RuleRepo) Delete(ctx context.Context, id uint64) error {
	if _, err := r.app.DB().NewDelete().Model(&Rule{ID: id}).Exec(ctx); err != nil {
		return err
	}
	return nil
}
