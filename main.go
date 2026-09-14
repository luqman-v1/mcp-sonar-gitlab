package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/luqman-v1/mcp-sonar-gitlab/pkg/gitlab"
	"github.com/luqman-v1/mcp-sonar-gitlab/pkg/resolver"
	"github.com/luqman-v1/mcp-sonar-gitlab/pkg/sonar"
)

var (
	version = "1.0.0"
)

func main() {
	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("mcp-sonar-gitlab v%s\n", version)
		os.Exit(0)
	}

	mcpServer := server.NewMCPServer(
		"mcp-sonar-gitlab",
		version,
		server.WithToolCapabilities(true),
	)

	tool := mcp.NewTool("fetch_sonar_issues_by_mr",
		mcp.WithDescription("Fetch SonarQube / SonarCloud issues and quality gate status for a given GitLab Merge Request URL."),
		mcp.WithString("mr_url",
			mcp.Required(),
			mcp.Description("The GitLab Merge Request URL (e.g. https://gitlab.com/group/repo/-/merge_requests/123)"),
		),
	)

	mcpServer.AddTool(tool, handleFetchSonarIssues)

	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Fprintf(os.Stderr, "MCP server stdio error: %v\n", err)
		os.Exit(1)
	}
}

func handleFetchSonarIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mrURL := mcp.ParseString(req, "mr_url", "")
	if mrURL == "" {
		return mcp.NewToolResultError("Argument 'mr_url' is required"), nil
	}

	sonarToken := os.Getenv("SONAR_TOKEN")
	if sonarToken == "" {
		return mcp.NewToolResultError("SONAR_TOKEN environment variable is not set. Please provide it in your MCP configuration under 'env'."), nil
	}

	sonarHost := os.Getenv("SONAR_HOST_URL")
	if sonarHost == "" {
		sonarHost = "https://sonarcloud.io"
	}
	sonarOrg := os.Getenv("SONAR_ORGANIZATION")

	mrInfo, err := gitlab.ParseMRURL(mrURL)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid MR URL: %v", err)), nil
	}

	gitlabBaseURL := os.Getenv("GITLAB_BASE_URL")
	if gitlabBaseURL == "" {
		gitlabBaseURL = mrInfo.BaseURL
	}
	gitlabToken := os.Getenv("GITLAB_TOKEN")

	var glClient *gitlab.Client
	if gitlabToken != "" {
		glClient = gitlab.NewClient(gitlabBaseURL, gitlabToken)
	}

	componentKey := resolver.ResolveComponentKey(ctx, glClient, mrInfo.ProjectPath)
	if componentKey == "" {
		return mcp.NewToolResultError(fmt.Sprintf("Could not resolve Sonar component key for path: %s", mrInfo.ProjectPath)), nil
	}

	sonarClient := sonar.NewClient(sonarHost, sonarToken, sonarOrg)

	issues, err := sonarClient.GetMRIssues(ctx, componentKey, mrInfo.MRIID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch Sonar issues for %s MR !%d: %v", componentKey, mrInfo.MRIID, err)), nil
	}

	qg, qgErr := sonarClient.GetQualityGate(ctx, componentKey, mrInfo.MRIID)
	if qgErr != nil {
		// Log warning but don't fail, Quality Gate might not be activated
		fmt.Fprintf(os.Stderr, "[mcp-sonar-gitlab] warning: failed to fetch quality gate for %s: %v\n", componentKey, qgErr)
	}

	resultText := formatOutput(componentKey, mrInfo.MRIID, issues, qg)
	return mcp.NewToolResultText(resultText), nil
}

func formatOutput(componentKey string, mrIID int, issues []sonar.SonarIssue, qg *sonar.SonarQualityGate) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("### Sonar Analysis: `%s` (MR !%d)\n\n", componentKey, mrIID))

	// 1. Quality Gate Summary
	if qg != nil {
		statusBadge := qg.Status
		switch strings.ToUpper(qg.Status) {
		case "OK":
			statusBadge = "PASSED (OK)"
		case "ERROR":
			statusBadge = "FAILED (ERROR)"
		case "WARN":
			statusBadge = "WARNING (WARN)"
		}
		sb.WriteString(fmt.Sprintf("**Quality Gate Status:** %s\n\n", statusBadge))

		if len(qg.Conditions) > 0 {
			sb.WriteString("**Conditions:**\n")
			for _, cond := range qg.Conditions {
				sb.WriteString(fmt.Sprintf("- `%s`: actual `%s` (comparator: `%s`, threshold: `%s`) -> **%s**\n",
					cond.MetricKey, cond.ActualValue, cond.Comparator, cond.ErrorThreshold, cond.Status))
			}
			sb.WriteString("\n")
		}
	}

	// 2. Issues Summary
	sb.WriteString(fmt.Sprintf("**Total Open Issues:** %d\n\n", len(issues)))

	if len(issues) == 0 {
		sb.WriteString("No open Sonar issues found for this MR. Everything looks clean!\n")
		return sb.String()
	}

	// Severity sort ordering
	severityRank := map[string]int{
		"BLOCKER":  1,
		"CRITICAL": 2,
		"MAJOR":    3,
		"MINOR":    4,
		"INFO":     5,
	}

	sort.SliceStable(issues, func(i, j int) bool {
		rI := severityRank[strings.ToUpper(issues[i].Severity)]
		if rI == 0 {
			rI = 99
		}
		rJ := severityRank[strings.ToUpper(issues[j].Severity)]
		if rJ == 0 {
			rJ = 99
		}
		return rI < rJ
	})

	for i, issue := range issues {
		line := issue.Line
		if line == 0 && issue.TextRange != nil {
			line = issue.TextRange.StartLine
		}

		location := issue.Component
		if line > 0 {
			location = fmt.Sprintf("%s:%d", issue.Component, line)
		}

		sb.WriteString(fmt.Sprintf("%d. **[%s]** %s\n", i+1, strings.ToUpper(issue.Severity), issue.Message))
		sb.WriteString(fmt.Sprintf("   - **Location:** `%s`\n", location))
		sb.WriteString(fmt.Sprintf("   - **Rule:** `%s`", issue.Rule))
		if issue.Type != "" {
			sb.WriteString(fmt.Sprintf(" | **Type:** `%s`", issue.Type))
		}
		if issue.Effort != "" {
			sb.WriteString(fmt.Sprintf(" | **Effort:** `%s`", issue.Effort))
		}
		sb.WriteString("\n\n")
	}

	return sb.String()
}
