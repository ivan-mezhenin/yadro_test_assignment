package domain

import (
	"context"
	"errors"
)

var (
	ErrAlreadyExists = errors.New("dns server already exists")
	ErrNotFound      = errors.New("dns server not found")
	ErrInvalidAddr   = errors.New("invalid address")
)

type DnsRepository interface {
	Add(ctx context.Context, address string) error
	Remove(ctx context.Context, address string) error
	List(ctx context.Context) ([]string, error)
}
