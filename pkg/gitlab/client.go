package gitlab

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func NewClient(baseURL, token string) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}
	return &Client{
		baseURL: baseURL,
		token:   token,
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// GetRawFileContent fetches the raw content of a file from a GitLab project repository.
func (c *Client) GetRawFileContent(ctx context.Context, projectPath, filePath, ref string) (string, error) {
	if c.token == "" {
		return "", fmt.Errorf("gitlab token is empty")
	}

	encodedProjectPath := url.PathEscape(projectPath)
	encodedFilePath := url.PathEscape(filePath)

	endpoint := fmt.Sprintf("%s/api/v4/projects/%s/repository/files/%s/raw", c.baseURL, encodedProjectPath, encodedFilePath)
	if ref != "" {
		endpoint += "?ref=" + url.QueryEscape(ref)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("PRIVATE-TOKEN", c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute gitlab request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gitlab api returned status %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(bodyBytes), nil
}
