# mcp-sonar-gitlab-

[![Go Report Card](https://goreportcard.com/badge/github.com/luqman-v1/mcp-sonar-gitlab-)](https://goreportcard.com/report/github.com/luqman-v1/mcp-sonar-gitlab-)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Model Context Protocol (MCP) server written in Go for fetching SonarQube / SonarCloud issues and quality gates directly from GitLab Merge Request URLs.

Designed for engineering teams to use locally with AI coding assistants (Cursor, Windsurf, Claude Desktop, Antigravity CLI).

---

## ⚡ Key Features

- **Standard `stdio` MCP Server:** Runs directly inside your AI IDE without requiring a long-running background web server or port binding.
- **GitLab MR to Sonar Mapping:** Automatically parses GitLab MR URLs, resolves project component keys, and queries SonarQube / SonarCloud.
- **Dual Auth Support:** Transparently supports SonarCloud (`Bearer` token) and on-premise SonarQube (auto-fallback to HTTP Basic Auth).
- **Flexible Namespace Mapping:** Supports custom path prefix mapping via `NAMESPACE_PREFIX_MAP` or automatically reads `sonar.projectKey` from `sonar-project.properties`.
- **Clean Markdown Summaries:** Formats Quality Gate status, condition breakdowns, and open issues ordered by severity (`BLOCKER`, `CRITICAL`, `MAJOR`, `MINOR`, `INFO`).

---

## 📦 Installation

### 1. Install via `go install` (Recommended)

Make sure Go is installed and `$GOPATH/bin` (or `~/go/bin`) is in your system's `PATH`:

```bash
go install github.com/luqman-v1/mcp-sonar-gitlab-@latest
```

Verify the installation:

```bash
mcp-sonar-gitlab- --version
# Output: mcp-sonar-gitlab v1.0.0
```

*(Optional) Alternatively, build from source:*

```bash
git clone git@github.com:luqman-v1/mcp-sonar-gitlab-.git
cd mcp-sonar-gitlab-
go build -o ~/go/bin/mcp-sonar-gitlab- .
```

---

## 🔑 Environment Variables

Each team member configures their own tokens in their local MCP settings:

| Variable | Required | Default | Description |
|---|---|---|---|
| `SONAR_TOKEN` | **Yes** | - | Your personal SonarQube / SonarCloud token |
| `SONAR_HOST_URL` | No | `https://sonarcloud.io` | Sonar URL (e.g. `https://sonar.internal.company.com`) |
| `SONAR_ORGANIZATION` | No | - | SonarCloud organization key (if applicable) |
| `GITLAB_TOKEN` | No | - | GitLab personal access token (used to read `sonar-project.properties`) |
| `GITLAB_BASE_URL` | No | Inferred from MR | GitLab instance base URL (e.g. `https://gitlab.com` or self-hosted) |
| `NAMESPACE_PREFIX_MAP` | No | - | Custom namespace prefix mapping (e.g. `myorg/backend/=org/be/`) |

---

## 🛠 IDE Configuration

### 1. Cursor

Add to `~/.cursor/mcp.json` (or via **Settings > Features > MCP Servers**):

```json
{
  "mcpServers": {
    "sonar-gitlab": {
      "command": "mcp-sonar-gitlab-",
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

*Tip:* If `mcp-sonar-gitlab-` is not in your global PATH, provide the absolute path: `"/Users/yourname/go/bin/mcp-sonar-gitlab-"`.

---

### 2. Claude Desktop

Add to your Claude configuration file:
- **macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows:** `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "sonar-gitlab": {
      "command": "mcp-sonar-gitlab-",
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
      "command": "mcp-sonar-gitlab-",
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
      "command": "mcp-sonar-gitlab-",
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

> *"Cek sonar issue untuk MR ini: https://gitlab.com/mygroup/backend/order-service/-/merge_requests/45"*

> *"Fetch all blocker and critical Sonar issues for MR https://gitlab.com/.../-/merge_requests/12 and suggest code fixes."*

> *"Periksa apakah Quality Gate lolos untuk MR ini sebelum kita merge."*

---

## 📄 License

[MIT](LICENSE)
