package repository

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/marktrs/gitsast/internal/model"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bunrouter"
)

// HTTPHandler variable that does static check to make sure that httpHandler struct implements HTTPHandler interface.
var _ HTTPHandler = (*httpHandler)(nil)

var (
	ErrInvalidParam = errors.New("error invalid parameter")
)

// HTTPHandler defines methods for http handler of repository domain
// such as parse request, query and create response
type HTTPHandler interface {
	GetById(http.ResponseWriter, bunrouter.Request) error
	List(http.ResponseWriter, bunrouter.Request) error
	Add(http.ResponseWriter, bunrouter.Request) error
	Update(http.ResponseWriter, bunrouter.Request) error
	Remove(http.ResponseWriter, bunrouter.Request) error
	Scan(http.ResponseWriter, bunrouter.Request) error
	GetReport(http.ResponseWriter, bunrouter.Request) error
}

type httpHandler struct {
	service IService
}

func NewHTTPHandler(s IService) HTTPHandler {
	return &httpHandler{
		service: s,
	}
}

// GetById implements HTTPHandler.GetById interface.
func (h *httpHandler) GetById(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	params := req.Params().Map()
	id, ok := params["id"]
	if !ok {
		log.Err(ErrInvalidParam).Msg("unable to get report by repo ID")
		return ErrInvalidParam
	}

	repo, err := h.service.GetById(ctx, id)
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, &repo)
}

// List implements HTTPHandler.List interface.
func (h *httpHandler) List(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	f, err := model.DecodeRepositoryFilter(req)
	if err != nil {
		return err
	}

	repos, err := h.service.List(ctx, f)
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, bunrouter.H{
		"repositories": &repos,
		"total":        len(repos),
	})
}

// Add implements HTTPHandler.Add interface.
func (h *httpHandler) Add(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	var r *AddRepositoryRequest
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		return err
	}

	repo, err := h.service.Add(ctx, r)
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, repo)
}

// Update implements HTTPHandler.Update interface.
func (h *httpHandler) Update(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	params := req.Params().Map()
	id, ok := params["id"]
	if !ok {
		log.Err(ErrInvalidParam).Msg("unable to update by repo ID")
		return ErrInvalidParam
	}

	var r *UpdateRepositoryRequest
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		return err
	}

	err := h.service.Update(ctx, id, r)
	if err != nil {
		return err
	}

	return nil
}

// Remove implements HTTPHandler.Remove interface.
func (h *httpHandler) Remove(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	params := req.Params().Map()
	id, ok := params["id"]
	if !ok {
		log.Err(ErrInvalidParam).Msg("unable to remove by repo ID")
		return ErrInvalidParam
	}

	return h.service.Remove(ctx, id)
}

// Scan implements HTTPHandler.Scan interface.
func (h *httpHandler) Scan(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	params := req.Params().Map()

	id, ok := params["id"]
	if !ok {
		log.Err(ErrInvalidParam).Msg("unable to scan by repo ID")
		return ErrInvalidParam
	}

	report, err := h.service.CreateReport(ctx, id)
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, &report)
}

// GetReport implements HTTPHandler.GetReport interface.
func (h *httpHandler) GetReport(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	params := req.Params().Map()

	id, ok := params["id"]
	if !ok {
		log.Err(ErrInvalidParam).
			Any("params", req.Params().Map()).
			Msg("unable to get report by repo ID")
		return ErrInvalidParam
	}

	response, err := h.service.GetReportByRepoId(ctx, id)
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, &response)
}
