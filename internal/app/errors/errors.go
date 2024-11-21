package errors

import (
	"errors"
)

var (
	ErrNotFound   = errors.New("URL not found in repository")
	ErrInternal   = errors.New("internal URL service error")
	ErrInvalid    = errors.New("invalid URL format")
	ErrDecompress = errors.New("request decompression error")

	ErrDBOpen     = errors.New("utils db: failed to open db")
	ErrDBPing     = errors.New("utils db: failed to ping db")
	ErrDBDSNParse = errors.New("utils db: failed to parse dsn")
	ErrDBClose    = errors.New("utils db: close error")
)
