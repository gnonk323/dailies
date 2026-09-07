package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"dailies/pkg/types"
)

type APIClient struct {
	BaseURL string
	Client  *http.Client
}

func NewFromEnv() (*APIClient, error) {
	base := os.Getenv("DAILIES_SERVER")
	if base == "" {
		base = "http://127.0.0.1:8080"
	}
	return New(base), nil
}

func New(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Client:  &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *APIClient) GetConfig() (types.DailiesConfig, error) {
	var cfg types.DailiesConfig
	if err := c.do("GET", "/api/config", nil, &cfg); err != nil {
		return types.DailiesConfig{}, err
	}
	return cfg, nil
}

func (c *APIClient) SaveConfig(cfg types.DailiesConfig) error {
	return c.do("POST", "/api/config", cfg, nil)
}

func (c *APIClient) GetEntry(date string) (*types.DailyEntry, error) {
	var entry types.DailyEntry
	err := c.do("GET", "/api/entries/"+date, nil, &entry)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &entry, nil
}

func (c *APIClient) ListEntries() ([]types.DailyEntry, error) {
	var entries []types.DailyEntry
	if err := c.do("GET", "/api/entries", nil, &entries); err != nil {
		return nil, err
	}
	if entries == nil {
		return []types.DailyEntry{}, nil
	}
	return entries, nil
}

func (c *APIClient) SaveEntry(entry *types.DailyEntry) error {
	return c.do("POST", "/api/entries", entry, nil)
}

func (c *APIClient) FetchIntegration(date, name string) (*types.DailyEntry, error) {
	var entry types.DailyEntry
	if err := c.do("POST", "/api/entries/"+date+"/fetch/"+name, nil, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

func (c *APIClient) FetchAllIntegrations(date string) error {
	return c.do("POST", "/api/entries/"+date+"/fetch", nil, nil)
}

type httpError struct {
	Status int
	Body   string
}

func (e *httpError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("server returned %d: %s", e.Status, e.Body)
	}
	return fmt.Sprintf("server returned %d", e.Status)
}

func isNotFound(err error) bool {
	var httpErr *httpError
	if errors.As(err, &httpErr) {
		return httpErr.Status == http.StatusNotFound
	}
	return false
}

func (c *APIClient) do(method, path string, body any, dest any) error {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reader)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("could not reach dailies server at %s: %w", c.BaseURL, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return &httpError{Status: resp.StatusCode, Body: strings.TrimSpace(string(raw))}
	}
	if dest == nil || len(raw) == 0 {
		return nil
	}
	return json.Unmarshal(raw, dest)
}
