package metadata

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("Not found")
)

type Scrobbler interface {
	SearchArtist(ctx context.Context, name string) (*Artist, error)
}
