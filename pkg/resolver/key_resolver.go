package resolver

import (
	"context"
	"os"
	"strings"

	"github.com/luqman-v1/mcp-sonar-gitlab/pkg/gitlab"
)

// ResolveComponentKey resolves the Sonar component key from GitLab repository properties or namespace path.
func ResolveComponentKey(ctx context.Context, glClient *gitlab.Client, projectPath string) string {
	// 1. Try to read sonar-project.properties from repository root if GitLab client is available
	if glClient != nil && projectPath != "" {
		content, err := glClient.GetRawFileContent(ctx, projectPath, "sonar-project.properties", "")
		if err == nil && content != "" {
			lines := strings.Split(content, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "sonar.projectKey=") {
					key := strings.TrimSpace(strings.TrimPrefix(line, "sonar.projectKey="))
					if key != "" {
						return key
					}
				}
			}
		}
	}

	// 2. Fallback to path-based mapping
	return MapPathToComponentKey(projectPath)
}

// MapPathToComponentKey converts a GitLab project path into Sonar component key convention.
// Supports custom namespace prefix mapping via NAMESPACE_PREFIX_MAP environment variable
// (e.g. NAMESPACE_PREFIX_MAP="mygroup/sub/backend/=sb/be/").
func MapPathToComponentKey(projectPath string) string {
	path := strings.Trim(projectPath, "/")
	if path == "" {
		return ""
	}

	// Apply custom prefix mapping if configured: "fromA=toA,fromB=toB"
	if mapping := os.Getenv("NAMESPACE_PREFIX_MAP"); mapping != "" {
		pairs := strings.Split(mapping, ",")
		for _, pair := range pairs {
			parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
			if len(parts) == 2 && parts[0] != "" {
				if strings.HasPrefix(path, parts[0]) {
					path = strings.Replace(path, parts[0], parts[1], 1)
					break
				}
			}
		}
	}

	return strings.ReplaceAll(path, "/", ":")
}
