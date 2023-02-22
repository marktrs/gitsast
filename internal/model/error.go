package model

import "errors"

var (
	ErrReportInProgress = errors.New(
		`the report for this repository already initialized, only completed/failed report can retry`)
)
