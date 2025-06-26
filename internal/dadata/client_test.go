package dadata

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

// mockTransport реализует http.RoundTripper
type mockTransport struct {
	response *http.Response
	err      error
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func TestClient_Search(t *testing.T) {
	// Фейковый ответ Dadata API
	mockBody := `{
		"suggestions": [{
			"data": {
				"city": "Москва",
				"street": "Тверская",
				"house": "1",
				"geo_lat": "55.7558",
				"geo_lon": "37.6173"
			}
		}]
	}`

	mt := &mockTransport{
		response: &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(mockBody))),
			Header:     make(http.Header),
		},
	}

	client := New("apiKey", "secretKey", 2*time.Second)
	client.http.Transport = mt // подменяем http

	addrs, err := client.Search(context.Background(), "Москва")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(addrs) != 1 || addrs[0].City != "Москва" {
		t.Fatalf("wrong address: %+v", addrs)
	}
}

func TestClient_Geocode(t *testing.T) {
	mockBody := `{
		"suggestions": [{
			"data": {
				"city": "Сочи",
				"street": "Островского",
				"house": "2",
				"geo_lat": "43.6",
				"geo_lon": "39.7"
			}
		}]
	}`

	mt := &mockTransport{
		response: &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(mockBody))),
			Header:     make(http.Header),
		},
	}

	client := New("apiKey", "secretKey", 2*time.Second)
	client.http.Transport = mt

	addrs, err := client.Geocode(context.Background(), "43.6", "39.7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(addrs) != 1 || addrs[0].City != "Сочи" {
		t.Fatalf("wrong address: %+v", addrs)
	}
}

func TestClient_call_Upstream500(t *testing.T) {
	mt := &mockTransport{
		response: &http.Response{
			StatusCode: 500,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
			Header:     make(http.Header),
		},
	}

	client := New("apiKey", "secretKey", 2*time.Second)
	client.http.Transport = mt

	_, err := client.Search(context.Background(), "fail")
	if err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}
