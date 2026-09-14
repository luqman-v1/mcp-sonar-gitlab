package gitlab_test

import (
	"testing"

	"github.com/luqman-v1/mcp-sonar-gitlab/pkg/gitlab"
)

func TestParseMRURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		wantHost    string
		wantBase    string
		wantPath    string
		wantIID     int
		wantErr     bool
	}{
		{
			name:     "standard gitlab.com MR",
			url:      "https://gitlab.com/myorg/team/backend/trade-service/-/merge_requests/123",
			wantHost: "gitlab.com",
			wantBase: "https://gitlab.com",
			wantPath: "myorg/team/backend/trade-service",
			wantIID:  123,
			wantErr:  false,
		},
		{
			name:     "custom domain MR with diffs tab suffix",
			url:      "https://gitlab.internal-corp.net/infra/devops/-/merge_requests/42/diffs",
			wantHost: "gitlab.internal-corp.net",
			wantBase: "https://gitlab.internal-corp.net",
			wantPath: "infra/devops",
			wantIID:  42,
			wantErr:  false,
		},
		{
			name:    "invalid URL",
			url:     "https://gitlab.com/not-an-mr",
			wantErr: true,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := gitlab.ParseMRURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseMRURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if info.Host != tt.wantHost {
				t.Errorf("Host = %v, want %v", info.Host, tt.wantHost)
			}
			if info.BaseURL != tt.wantBase {
				t.Errorf("BaseURL = %v, want %v", info.BaseURL, tt.wantBase)
			}
			if info.ProjectPath != tt.wantPath {
				t.Errorf("ProjectPath = %v, want %v", info.ProjectPath, tt.wantPath)
			}
			if info.MRIID != tt.wantIID {
				t.Errorf("MRIID = %v, want %v", info.MRIID, tt.wantIID)
			}
		})
	}
}
