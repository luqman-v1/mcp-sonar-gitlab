# mcp-sonar-gitlab

[![Go Report Card](https://goreportcard.com/badge/github.com/luqman-v1/mcp-sonar-gitlab)](https://goreportcard.com/report/github.com/luqman-v1/mcp-sonar-gitlab)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Model Context Protocol (MCP) server written in Go for fetching SonarQube / SonarCloud issues and quality gates directly from GitLab Merge Request URLs.

Designed for engineering teams to use locally with AI coding assistants (Cursor, Windsurf, Claude Desktop, Antigravity CLI).

---

## ⚡ Key Features

- **Standard `stdio` MCP Server:** Runs directly inside your AI IDE without requiring a background daemon or port binding.
- **GitLab MR to Sonar Mapping:** Automatically parses GitLab MR URLs, resolves project component keys, and queries SonarQube / SonarCloud.
- **Dual Auth Support:** Transparently supports SonarCloud (`Bearer` token) and on-premise SonarQube (auto-fallback to HTTP Basic Auth).
- **Flexible Namespace Mapping:** Supports custom path prefix mapping via `NAMESPACE_PREFIX_MAP` or automatically reads `sonar.projectKey` from `sonar-project.properties`.
- **Clean Markdown Summaries:** Formats Quality Gate status, condition breakdowns, and open issues ordered by severity (`BLOCKER`, `CRITICAL`, `MAJOR`, `MINOR`, `INFO`).

---

## 📦 Installation

### 1. Install via `go install` (Recommended)

Requires Go 1.27+. Make sure `$GOPATH/bin` (or `~/go/bin`) is in your system's `PATH`:

```bash
go install github.com/luqman-v1/mcp-sonar-gitlab@latest
```

Verify the installation:

```bash
mcp-sonar-gitlab --version
# Output: mcp-sonar-gitlab v1.0.0
```

*(Optional) Alternatively, build from source:*

```bash
git clone git@github.com:luqman-v1/mcp-sonar-gitlab.git
cd mcp-sonar-gitlab
go build -o ~/go/bin/mcp-sonar-gitlab .
```

---

## 🔑 Environment Variables

Each team member configures their own credentials and preferences in their local MCP settings:

| Variable | Status | Default | Description |
|---|---|---|---|
| `SONAR_TOKEN` | **Required** | - | Your personal SonarQube / SonarCloud access token. |
| `SONAR_HOST_URL` | *Optional* | `https://sonarcloud.io` | Sonar instance URL. Required if using an internal on-premise SonarQube instance (e.g. `https://sonar.internal.company.com`). |
| `GITLAB_TOKEN` | *Optional* | - | GitLab Personal Access Token with `read_api` or `read_repository` scope. Used to read `sonar-project.properties` directly from the repository. If omitted, the resolver falls back to `NAMESPACE_PREFIX_MAP`. |
| `NAMESPACE_PREFIX_MAP` | *Optional* | - | Mapping rule from GitLab namespace path to Sonar component key (format: `old_prefix/=new_prefix/`, comma-separated for multiple mappings). Useful when repos lack `sonar-project.properties`. |
| `GITLAB_BASE_URL` | *Optional* | Inferred from MR URL | GitLab instance base URL. Automatically detected from the host of the provided MR URL. |
| `SONAR_ORGANIZATION` | *Optional* | - | SonarCloud organization key (only needed if your SonarCloud account requires an organization parameter). |

### Variable Behavior & Fallbacks:

1. **`SONAR_TOKEN` (Required)**:
   - Required for Sonar API authentication.
   - Works with both SonarCloud (Bearer token) and SonarQube on-premise (automatically falls back to Basic Auth `<token>:` if Bearer is rejected with `401 Unauthorized`).
2. **`SONAR_HOST_URL` (Optional)**:
   - Defaults to `https://sonarcloud.io`.
   - Set to your private/corporate SonarQube domain if not using public SonarCloud.
3. **`GITLAB_TOKEN` (Optional)**:
   - When provided, queries the GitLab API to check if `sonar-project.properties` exists in the repository root. If `sonar.projectKey=...` is found, that key is prioritized.
   - If omitted or properties file is not found, resolution proceeds directly to path mapping.
4. **`NAMESPACE_PREFIX_MAP` (Optional)**:
   - Translates repo path prefixes into Sonar component key format.
   - Format: `source_path/=target_prefix/` (e.g. `mygroup/project/backend/=org/be/`).
   - Slashes (`/`) are automatically normalized to colons (`:`).
5. **`GITLAB_BASE_URL` (Optional)**:
   - Not needed if the MR URL is a full URL (e.g. `https://gitlab.mycorp.com/...`), as the host is automatically parsed.

---

### 📋 Configuration Examples

#### Option 1: Recommended Team Setup
```json
"env": {
  "SONAR_TOKEN": "sqp_your_personal_sonar_token",
  "SONAR_HOST_URL": "https://sonarcloud.io",
  "GITLAB_TOKEN": "glpat_your_personal_gitlab_token",
  "NAMESPACE_PREFIX_MAP": "mygroup/project/backend/=org/be/"
}
```

#### Option 2: Minimal Setup (SonarCloud only)
```json
"env": {
  "SONAR_TOKEN": "sqp_your_personal_sonar_token"
}
```

---

## 🛠 IDE Configuration

### 1. Cursor

Add to `~/.cursor/mcp.json` (or via **Settings > Features > MCP Servers**):

```json
{
  "mcpServers": {
    "sonar-gitlab": {
      "command": "mcp-sonar-gitlab",
      "env": {
        "SONAR_TOKEN": "sqp_your_personal_sonar_token",
        "SONAR_HOST_URL": "https://sonarcloud.io",
        "GITLAB_TOKEN": "glpat_your_personal_gitlab_token",
        "NAMESPACE_PREFIX_MAP": "mygroup/backend/=org/be/"
      }
    }
  }
}
```

*Tip:* If `mcp-sonar-gitlab` is not in your global PATH, provide the absolute path: `"/Users/yourname/go/bin/mcp-sonar-gitlab"`.

---

### 2. Claude Desktop

Add to your Claude configuration file:
- **macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows:** `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "sonar-gitlab": {
      "command": "mcp-sonar-gitlab",
      "env": {
        "SONAR_TOKEN": "sqp_your_personal_sonar_token",
        "SONAR_HOST_URL": "https://sonarcloud.io",
        "GITLAB_TOKEN": "glpat_your_personal_gitlab_token",
        "NAMESPACE_PREFIX_MAP": "mygroup/backend/=org/be/"
      }
    }
  }
}
```

---

### 3. Windsurf

Add to your `~/.codeium/windsurf/mcp_config.json`:

```json
{
  "mcpServers": {
    "sonar-gitlab": {
      "command": "mcp-sonar-gitlab",
      "env": {
        "SONAR_TOKEN": "sqp_your_personal_sonar_token",
        "SONAR_HOST_URL": "https://sonarcloud.io",
        "GITLAB_TOKEN": "glpat_your_personal_gitlab_token",
        "NAMESPACE_PREFIX_MAP": "mygroup/backend/=org/be/"
      }
    }
  }
}
```

---

### 4. Antigravity CLI / Gemini

Add to your `~/.gemini/config/mcp_config.json`:

```json
{
  "mcpServers": {
    "sonar-gitlab": {
      "command": "mcp-sonar-gitlab",
      "env": {
        "SONAR_TOKEN": "sqp_your_personal_sonar_token",
        "SONAR_HOST_URL": "https://sonarcloud.io",
        "GITLAB_TOKEN": "glpat_your_personal_gitlab_token",
        "NAMESPACE_PREFIX_MAP": "mygroup/backend/=org/be/"
      }
    }
  }
}
```

---

## 💡 How to Use with AI

Once registered, your AI assistant will automatically discover the `fetch_sonar_issues_by_mr` tool.

Example prompts in your IDE:

> *"Check Sonar issues for this MR: https://gitlab.com/mygroup/backend/order-service/-/merge_requests/45"*

> *"Fetch all blocker and critical Sonar issues for MR https://gitlab.com/.../-/merge_requests/12 and suggest code fixes."*

> *"Has this MR passed the Quality Gate? List all failing conditions if any."*

---

## 📄 License

[MIT](LICENSE)
