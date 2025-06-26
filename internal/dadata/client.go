package dadata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/your-org/geo-service-swagger/internal/address"
)

type Client struct {
	http      *http.Client
	apiKey    string
	secretKey string
	base      *url.URL
}

func New(apiKey, secret string, timeout time.Duration) *Client {
	return &Client{
		http:      &http.Client{Timeout: timeout},
		apiKey:    apiKey,
		secretKey: secret,
		base:      mustURL("https://suggestions.dadata.ru/suggestions/api/4_1/rs/"),
	}
}

func mustURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}

func (c *Client) Search(ctx context.Context, q string) ([]*address.Address, error) {
	reqBody, _ := json.Marshal(map[string]string{"query": q})
	endpoint := c.base.ResolveReference(&url.URL{Path: "suggest/address"})
	return c.call(ctx, endpoint.String(), reqBody)
}

func (c *Client) Geocode(ctx context.Context, lat, lng string) ([]*address.Address, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"lat": lat,
		"lon": lng, //  API ожидает lon
	})
	endpoint := c.base.ResolveReference(&url.URL{Path: "geolocate/address"})
	return c.call(ctx, endpoint.String(), reqBody)
}

func (c *Client) call(ctx context.Context, url string, body []byte) ([]*address.Address, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.apiKey))

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 500 {
		return nil, fmt.Errorf("dadata upstream 5xx: %d", res.StatusCode)
	}

	var raw struct {
		Suggestions []struct {
			Data struct {
				City   string `json:"city"`
				Street string `json:"street"`
				House  string `json:"house"`
				GeoLat string `json:"geo_lat"`
				GeoLon string `json:"geo_lon"`
			} `json:"data"`
		} `json:"suggestions"`
	}

	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var list []*address.Address
	for _, s := range raw.Suggestions {
		list = append(list, &address.Address{
			City:   s.Data.City,
			Street: s.Data.Street,
			House:  s.Data.House,
			Lat:    s.Data.GeoLat,
			Lon:    s.Data.GeoLon,
		})
	}
	return list, nil
}
