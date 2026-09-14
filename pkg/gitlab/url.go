package gitlab

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var mrURLRegex = regexp.MustCompile(`^https?://([^/]+)/(.+?)/-/merge_requests/(\d+)`)

// MRInfo holds extracted components from a GitLab MR URL.
type MRInfo struct {
	Host        string
	BaseURL     string
	ProjectPath string
	MRIID       int
}

// ParseMRURL extracts host, project namespace path, and MR IID from a GitLab MR URL.
// Supports URLs such as:
// https://gitlab.com/owner/subgroup/repo/-/merge_requests/123
// https://gitlab.mycorp.com/owner/repo/-/merge_requests/45
func ParseMRURL(rawURL string) (*MRInfo, error) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return nil, fmt.Errorf("empty MR URL")
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return nil, fmt.Errorf("invalid URL format: %w", err)
	}

	matches := mrURLRegex.FindStringSubmatch(trimmed)
	if len(matches) < 4 {
		return nil, fmt.Errorf("URL does not match GitLab MR pattern (expected .../<path>/-/merge_requests/<iid>): %s", rawURL)
	}

	iid, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid MR IID %q: %w", matches[3], err)
	}

	baseURL := fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)

	return &MRInfo{
		Host:        parsed.Host,
		BaseURL:     baseURL,
		ProjectPath: strings.TrimSuffix(matches[2], ".git"),
		MRIID:       iid,
	}, nil
}
