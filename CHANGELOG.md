# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-09-14

### Added
- Initial release of `mcp-sonar-gitlab`.
- Standard I/O (`stdio`) Model Context Protocol (MCP) server built with `mark3labs/mcp-go`.
- MCP tool `fetch_sonar_issues_by_mr` to fetch SonarQube / SonarCloud issues and Quality Gate status for any GitLab Merge Request URL.
- Transparent dual authentication: attempts SonarCloud `Bearer` token first, with automatic fallback to HTTP Basic Auth (`<token>:`) for on-premise SonarQube instances.
- GitLab URL parser supporting both `gitlab.com` and self-hosted GitLab instances.
- Component key resolver with automatic detection of `sonar-project.properties` via GitLab API and custom namespace prefix mapping via `NAMESPACE_PREFIX_MAP`.
- Structured Markdown reporting with Quality Gate status, condition breakdown, and issues ordered by severity (`BLOCKER`, `CRITICAL`, `MAJOR`, `MINOR`, `INFO`).
- Built for Go 1.27+.
