package errors

import (
	"errors"
)

var (
	ErrNotFound   = errors.New("URL not found in repository")
	ErrInternal   = errors.New("internal URL service error")
	ErrInvalid    = errors.New("invalid URL format")
	ErrDecompress = errors.New("request decompression error")
)
