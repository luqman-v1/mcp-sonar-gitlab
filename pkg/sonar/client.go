package sonar

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	baseURL      string
	token        string
	organization string
	http         *http.Client
}

func NewClient(baseURL, token, organization string) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	if baseURL == "" {
		baseURL = "https://sonarcloud.io"
	}
	return &Client{
		baseURL:      baseURL,
		token:        token,
		organization: organization,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) doRequest(ctx context.Context, reqURL *url.URL) (*http.Response, error) {
	if c.token == "" {
		return nil, fmt.Errorf("sonar token is required")
	}

	// 1. Try Bearer token (SonarCloud standard)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	// 2. Fallback to Basic Auth if unauthorized (SonarQube on-premise requires token as username and empty password)
	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()

		reqRetry, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to build retry request: %w", err)
		}
		reqRetry.SetBasicAuth(c.token, "")

		resp, err = c.http.Do(reqRetry)
		if err != nil {
			return nil, fmt.Errorf("failed to execute retry request: %w", err)
		}
	}

	return resp, nil
}

// GetMRIssues fetches open and confirmed issues for the given componentKey and pullRequestID.
func (c *Client) GetMRIssues(ctx context.Context, componentKey string, pullRequestID int) ([]SonarIssue, error) {
	endpoint := fmt.Sprintf("%s/api/issues/search", c.baseURL)
	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint URL: %w", err)
	}

	q := reqURL.Query()
	q.Set("componentKeys", componentKey)
	q.Set("pullRequest", strconv.Itoa(pullRequestID))
	q.Set("statuses", "OPEN,CONFIRMED")
	if c.organization != "" {
		q.Set("organization", c.organization)
	}
	reqURL.RawQuery = q.Encode()

	resp, err := c.doRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("sonar api returned status %d: %s", resp.StatusCode, string(body))
	}

	var searchRes SonarSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchRes); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	return searchRes.Issues, nil
}

// GetQualityGate fetches the quality gate status for the given componentKey and pullRequestID.
func (c *Client) GetQualityGate(ctx context.Context, componentKey string, pullRequestID int) (*SonarQualityGate, error) {
	endpoint := fmt.Sprintf("%s/api/qualitygates/project_status", c.baseURL)
	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint URL: %w", err)
	}

	q := reqURL.Query()
	q.Set("projectKey", componentKey)
	q.Set("pullRequest", strconv.Itoa(pullRequestID))
	if c.organization != "" {
		q.Set("organization", c.organization)
	}
	reqURL.RawQuery = q.Encode()

	resp, err := c.doRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("sonar quality gate api returned status %d: %s", resp.StatusCode, string(body))
	}

	var res sonarQualityGateResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode quality gate response: %w", err)
	}

	return &res.ProjectStatus, nil
}
