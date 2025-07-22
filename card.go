// Package metabase provides a Go client for Metabase public card API
package metabase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client represents Metabase public API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient returns new client with optional custom http.Client
func NewClient(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		HTTPClient: httpClient,
	}
}

// CardQuery executes a GET request to Metabase card with filters
func (c *Client) CardQuery(
	ctx context.Context,
	uuid string,
	format Format,
	filters []Filter,
) ([]byte, error) {
	if !format.Valid() {
		return nil, fmt.Errorf("invalid format: %s", format)
	}

	endpoint := fmt.Sprintf("%s/api/public/card/%s/query/%s", c.BaseURL, uuid, format)
	filterBytes, err := json.Marshal(filters)
	if err != nil {
		return nil, fmt.Errorf("marshal filters: %w", err)
	}

	u := fmt.Sprintf("%s?%s", endpoint, url.Values{"parameters": {string(filterBytes)}}.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("metabase error: %s\n%s", resp.Status, string(b))
	}

	return io.ReadAll(resp.Body)
}
