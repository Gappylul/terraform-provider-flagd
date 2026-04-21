package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Flag struct {
	Name        string    `json:"name"`
	Enabled     bool      `json:"enabled"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type flagdClient struct {
	baseURL    string
	adminKey   string
	httpClient *http.Client
}

func newFlagdClient(baseURL, adminKey string) *flagdClient {
	return &flagdClient{
		baseURL:  baseURL,
		adminKey: adminKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *flagdClient) GetFlag(ctx context.Context, name string) (*Flag, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/flags/%s", c.baseURL, name), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("flagd request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("flagd returned status %d", resp.StatusCode)
	}

	var f Flag
	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &f, nil
}

func (c *flagdClient) UpsertFlag(ctx context.Context, name, description string, enabled bool) (*Flag, error) {
	body, _ := json.Marshal(map[string]any{
		"enabled":     enabled,
		"description": description,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPut,
		fmt.Sprintf("%s/flags/%s", c.baseURL, name),
		bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.adminKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("flagd request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("flagd returned status %d", resp.StatusCode)
	}

	var f Flag
	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &f, nil
}

func (c *flagdClient) DeleteFlag(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete,
		fmt.Sprintf("%s/flags/%s", c.baseURL, name), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.adminKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("flagd request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("flagd returned status %d", resp.StatusCode)
	}
	return nil
}
