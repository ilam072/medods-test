package postgres

import "errors"

const (
	ErrUniqueViolationCode = "23505"
)

var (
	ErrUniqueContraintFailed = errors.New("unique constraint failed")
)
