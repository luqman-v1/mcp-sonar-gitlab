package resolver_test

import (
	"context"
	"testing"

	"github.com/luqman-v1/mcp-sonar-gitlab-/pkg/resolver"
)

func TestMapPathToComponentKey(t *testing.T) {
	t.Setenv("NAMESPACE_PREFIX_MAP", "myorg/enterprise/backend/=app/be/,finance/services/=fin/svc/")

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "custom mapped namespace",
			path:     "myorg/enterprise/backend/order-service",
			expected: "app:be:order-service",
		},
		{
			name:     "second mapped namespace",
			path:     "finance/services/ledger",
			expected: "fin:svc:ledger",
		},
		{
			name:     "generic repo without prefix mapping",
			path:     "openorg/mobile/android/core",
			expected: "openorg:mobile:android:core",
		},
		{
			name:     "path with leading and trailing slashes",
			path:     "/myorg/enterprise/backend/wallet/api/",
			expected: "app:be:wallet:api",
		},
		{
			name:     "empty path",
			path:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolver.MapPathToComponentKey(tt.path)
			if got != tt.expected {
				t.Errorf("MapPathToComponentKey(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}

func TestResolveComponentKeyWithoutGitLabClient(t *testing.T) {
	t.Setenv("NAMESPACE_PREFIX_MAP", "myorg/enterprise/backend/=app/be/")
	key := resolver.ResolveComponentKey(context.Background(), nil, "myorg/enterprise/backend/portfolio-service")
	expected := "app:be:portfolio-service"
	if key != expected {
		t.Errorf("ResolveComponentKey() = %q, want %q", key, expected)
	}
}
