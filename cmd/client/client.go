package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ivansevryukov1995/github-user-activity/internal/model"
)

type GitHubRequest struct {
	User        string
	CountEvents int
}

type GitHubClient struct {
	URL string
}

func (srv *GitHubClient) GetEvents(ctx context.Context, req GitHubRequest) ([]model.Event, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	events := make([]model.Event, req.CountEvents)

	urlEvents := fmt.Sprintf("%s/users/%s/events", srv.URL, req.User)

	eventReq, err := http.NewRequestWithContext(ctx, http.MethodGet, urlEvents, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(eventReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %v", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("failed decode: %w", err)
	}

	return events, nil
}
