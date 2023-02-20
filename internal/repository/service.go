package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

var _ IService = (*service)(nil)

// IService defines methods for business logic of repository domain
// such as validate request body, enqueue analyzing task, CRUD repository or report
type IService interface {
	GetById(ctx context.Context, id string) (*Repository, error)
	List(ctx context.Context, f *RepositoryFilter) ([]*Repository, error)
	Add(ctx context.Context, req *AddRepositoryRequest) (*Repository, error)
	Update(ctx context.Context, id string, req *UpdateRepositoryRequest) error
	Remove(ctx context.Context, id string) error
	CreateReport(ctx context.Context, repoId string) (*Report, error)
}

type service struct {
	repo      IRepository
	report    IReport
	validator *validator.Validate
}

type AddRepositoryRequest struct {
	Name      string `json:"name" validate:"required,max=120"`
	RemoteURL string `json:"remote_url" validate:"required,max=120"`
}

type UpdateRepositoryRequest struct {
	Name      string `json:"name"`
	RemoteURL string `json:"remote_url"`
}

func NewService(v *validator.Validate, rs IRepository, rp IReport) IService {
	return &service{
		repo:      rs,
		report:    rp,
		validator: v,
	}
}

// GetById implements IService.GetById interface.
func (s *service) GetById(ctx context.Context, id string) (*Repository, error) {
	return s.repo.GetById(ctx, id)
}

// List implements IService.List interface.
func (s *service) List(ctx context.Context, f *RepositoryFilter) ([]*Repository, error) {
	return s.repo.List(ctx, f)
}

// Add implements IService.Add interface.
func (s *service) Add(ctx context.Context, req *AddRepositoryRequest) (*Repository, error) {
	// validate request body
	if err := s.validator.Struct(req); err != nil {
		log.Errorf("Request validation failed on add repository handler : %s", err.Error())
		return nil, err
	}

	repo := &Repository{
		ID:        uuid.New().String(),
		Name:      req.Name,
		RemoteURL: req.RemoteURL,
	}

	return s.repo.Add(ctx, repo)
}

// Update implements IService.Update interface.
func (s *service) Update(ctx context.Context, id string, req *UpdateRepositoryRequest,
) error {
	// validate request body
	if err := s.validator.Struct(req); err != nil {
		log.Errorf("Request validation failed on update repository handler : %s", err.Error())
		return err
	}

	repo := map[string]interface{}{
		"name":       req.Name,
		"remote_url": req.RemoteURL,
		"updated_at": time.Now(),
	}

	return s.repo.Update(ctx, id, repo)
}

// Remove implements IService.Remove interface.
func (s *service) Remove(ctx context.Context, id string) error {
	return s.repo.Remove(ctx, id)
}

// CreateReport implements IService.CreateReport interface.
func (s *service) CreateReport(ctx context.Context, repoId string) (*Report, error) {
	report, err := s.report.GetByRepoId(ctx, repoId)
	if err != sql.ErrNoRows {
		return nil, err
	}

	// if report exist
	if report.ID != "" {
		if report.Status != StatusSuccess && report.Status != StatusFailed {
			return nil, errors.New(
				`the report for this repository already initialized, 
				only completed/failed report can retry`)
		}
	}

	// generate a new report
	report = &Report{
		ID:           uuid.New().String(),
		Status:       StatusInitialized,
		RepositoryID: repoId,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	report, err = s.report.Add(ctx, report)
	if err != sql.ErrNoRows {
		return nil, err
	}

	// TODO: enqueue
	// ...

	report.Status = StatusEnqueued
	report.EnqueueAt = time.Now()

	report, err = s.report.Update(ctx, report)
	if err != sql.ErrNoRows {
		return nil, err
	}

	return report, nil
}
