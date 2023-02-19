package repository

import (
	"errors"
	"strconv"

	"github.com/uptrace/bun"
	"github.com/uptrace/bunrouter"
)

type RepositoryFilter struct {
	ID        string
	Name      string
	RemoteURL string

	Limit  int
	Offset int
}

func decodeRepositoryFilter(req bunrouter.Request) (*RepositoryFilter, error) {
	var err error

	limit := 100
	offset := 0

	query := req.URL.Query()

	f := &RepositoryFilter{
		ID:        query.Get("id"),
		Name:      query.Get("name"),
		RemoteURL: query.Get("remote_url"),
		Limit:     limit,
		Offset:    offset,
	}

	if query.Has("limit") {
		limit, err = strconv.Atoi(query.Get("limit"))
		if err != nil {
			return nil, errors.Join(errors.New("invalid query param value: limit"), err)
		}
		f.Limit = limit
	}

	if query.Has("offset") {
		offset, err = strconv.Atoi(query.Get("offset"))
		if err != nil {
			return nil, errors.Join(errors.New("invalid query param value: offset"), err)
		}
		f.Offset = offset
	}

	return f, nil
}

func (f *RepositoryFilter) query(q *bun.SelectQuery) *bun.SelectQuery {
	// TODO: Add relation query
	return q
}
