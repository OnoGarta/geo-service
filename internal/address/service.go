package address

import (
	"context"
)

type Upstream interface {
	Search(ctx context.Context, query string) ([]*Address, error)
	Geocode(ctx context.Context, lat, lng string) ([]*Address, error)
}

type Service interface {
	Search(ctx context.Context, query string) ([]*Address, error)
	Geocode(ctx context.Context, lat, lng string) ([]*Address, error)
}

type service struct{ upstream Upstream }

func NewService(up Upstream) Service { return &service{upstream: up} }

func (s *service) Search(ctx context.Context, query string) ([]*Address, error) {
	return s.upstream.Search(ctx, query)
}

func (s *service) Geocode(ctx context.Context, lat, lng string) ([]*Address, error) {
	return s.upstream.Geocode(ctx, lat, lng)
}
