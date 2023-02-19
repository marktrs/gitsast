package repository

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bunrouter"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const kb = 10

var (
	ErrInvalidParam = errors.New("error invalid parameter")
)

type handler struct {
	db        *bun.DB
	validator *validator.Validate
}

func NewRepositoryHandler(db *bun.DB, v *validator.Validate) *handler {
	return &handler{
		db:        db,
		validator: v,
	}
}

// GetById 1.parse params, 3. fetch from db, 4. write response
// TODO: Refactor this function into separate layers
func (h *handler) GetById(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	params := req.Params().Map()
	id, ok := params["id"]
	if !ok {
		log.Errorf("unable to get repo by ID : %s %v", ErrInvalidParam.Error(), req.Params().Map())
		return ErrInvalidParam
	}

	// repository layer
	repo := &Repository{}
	err := h.db.NewSelect().Model(repo).Where("id = ?", id).Scan(ctx)
	if err != nil {
		log.Error(err)
		return err
	}
	// http layer
	return bunrouter.JSON(w, bunrouter.H{
		"repositories": &repo,
	})
}

// List 1.parse req body, 2. decode and validate filter param 3. fetch from db, 4. write response
// TODO: Refactor this function into separate layers
func (h *handler) List(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	// http layer
	f, err := decodeRepositoryFilter(req)
	if err != nil {
		return err
	}

	// repository layer
	repos := []*Repository{}
	err = h.db.NewSelect().Model(&repos).
		Apply(f.query).
		Limit(f.Limit).
		Offset(f.Offset).
		Scan(ctx)
	if err != nil {
		return err
	}

	// http layer
	return bunrouter.JSON(w, bunrouter.H{
		"repositories": &repos,
		"total":        len(repos),
	})
}

type AddRepositoryRequest struct {
	Name      string `json:"name" validate:"required"`
	RemoteURL string `json:"remote_url" validate:"required"`
}

// Add 1.parse req body, 2. validate body 3. insert to db, 4. write response
// TODO: Refactor this function into separate layers
func (h *handler) Add(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	// controller layer
	var r *AddRepositoryRequest
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		return err
	}

	// service layer
	// validate request body
	if err := h.validator.Struct(r); err != nil {
		log.Errorf("Request validation failed on add repository handler : %s", err.Error())
		return err
	}

	new := &Repository{
		ID:        uuid.New().String(),
		Name:      r.Name,
		RemoteURL: r.RemoteURL,
	}

	_, err := h.db.NewInsert().Model(new).Exec(ctx)
	if err != nil {
		return err
	}

	// controller layer
	return bunrouter.JSON(w, new)
}

type UpdateRepositoryRequest struct {
	Name      string `json:"name"`
	RemoteURL string `json:"remote_url"`
}

// Update 1.parse params 2.parse req body, 3. validate body 4. update record, 5. write response
// TODO: Refactor this function into separate layers
func (h *handler) Update(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	params := req.Params().Map()
	id, ok := params["id"]
	if !ok {
		log.Errorf("unable to get repo by ID : %s %v", ErrInvalidParam.Error(), req.Params().Map())
		return ErrInvalidParam
	}

	// controller layer
	var r *UpdateRepositoryRequest
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		return err
	}

	// service layer
	// validate request body
	if err := h.validator.Struct(r); err != nil {
		log.Errorf("Request validation failed on update repository handler : %s", err.Error())
		return err
	}

	repo := map[string]interface{}{
		"name":       r.Name,
		"remote_url": r.RemoteURL,
		"updated_at": time.Now(),
	}

	_, err := h.db.NewUpdate().
		Model(&repo).
		TableExpr("repositories").
		Where("id = ?", id).
		OmitZero().
		Exec(ctx)
	if err != nil {
		return err
	}

	// controller layer
	return bunrouter.JSON(w, repo)
}

// Remove 1.parse params  2. remove record, 3. write response
// TODO: Refactor this function into separate layers
func (h *handler) Remove(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	params := req.Params().Map()
	id, ok := params["id"]
	if !ok {
		log.Errorf("unable to get repo by ID : %s %v", ErrInvalidParam.Error(), req.Params().Map())
		return ErrInvalidParam
	}

	_, err := h.db.NewDelete().
		Model((*Repository)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	// controller layer
	return nil
}
