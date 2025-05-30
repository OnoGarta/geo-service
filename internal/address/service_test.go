package address

import (
	"context"
	"testing"
)

// Мок для Upstream (реализация)
type mockUpstream struct{}

func (m *mockUpstream) Search(ctx context.Context, q string) ([]*Address, error) {
	return []*Address{{City: "Москва"}}, nil
}
func (m *mockUpstream) Geocode(ctx context.Context, lat, lng string) ([]*Address, error) {
	return []*Address{{City: "Сочи"}}, nil
}

func TestService_Search(t *testing.T) {
	svc := NewService(&mockUpstream{})
	res, err := svc.Search(context.Background(), "Москва")
	if err != nil || len(res) == 0 || res[0].City != "Москва" {
		t.Fatalf("unexpected result: %+v, err: %v", res, err)
	}
}

func TestService_Geocode(t *testing.T) {
	svc := NewService(&mockUpstream{})
	res, err := svc.Geocode(context.Background(), "43.6", "39.7")
	if err != nil || len(res) == 0 || res[0].City != "Сочи" {
		t.Fatalf("unexpected result: %+v, err: %v", res, err)
	}
}
