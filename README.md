# mcp-sonar-gitlab

[![Go Report Card](https://goreportcard.com/badge/github.com/luqman-v1/mcp-sonar-gitlab)](https://goreportcard.com/report/github.com/luqman-v1/mcp-sonar-gitlab)
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

Requires Go 1.27+ installed. Make sure `$GOPATH/bin` (or `~/go/bin`) is in your system's `PATH`:

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

Setiap anggota tim mengatur token dan konfigurasinya sendiri di file setting MCP lokal mereka.

| Variable | Status | Default | Deskripsi |
|---|---|---|---|
| `SONAR_TOKEN` | **Required** | - | Token pribadi SonarQube / SonarCloud kamu. |
| `SONAR_HOST_URL` | *Optional* | `https://sonarcloud.io` | URL instance Sonar. Wajib diisi jika menggunakan SonarQube on-premise internal kantor (contoh: `https://sonar.internal.company.com`). |
| `GITLAB_TOKEN` | *Optional* | - | Personal Access Token GitLab dengan scope `read_api` atau `read_repository`. Digunakan untuk membaca file `sonar-project.properties` langsung dari repo. Jika dikosongkan, resolver akan menggunakan `NAMESPACE_PREFIX_MAP`. |
| `NAMESPACE_PREFIX_MAP` | *Optional* | - | Pemetaan prefix path namespace GitLab ke Sonar component key (format: `prefix_lama/=prefix_baru/`, bisa pisahkan dengan koma). Sangat berguna sebagai fallback jika repo tidak memiliki `sonar-project.properties`. |
| `GITLAB_BASE_URL` | *Optional* | Inferred dari MR URL | Base URL instance GitLab. Secara default otomatis diekstrak langsung dari domain MR URL yang kamu berikan. |
| `SONAR_ORGANIZATION` | *Optional* | - | Organization key pada SonarCloud (hanya diperlukan jika akun SonarCloud kamu memerlukan parameter organization). |

### Detail Cara Kerja Setiap Variabel:

1. **`SONAR_TOKEN` (Wajib)**:
   - Dibutuhkan untuk otentikasi API Sonar.
   - Mendukung SonarCloud (Bearer token) maupun SonarQube on-premise (otomatis fallback ke Basic Auth `token:` jika server menolak Bearer).
2. **`SONAR_HOST_URL` (Opsional)**:
   - Jika tidak diisi, otomatis menggunakan `https://sonarcloud.io`.
   - Isi dengan domain SonarQube kantor kamu jika bukan cloud publik.
3. **`GITLAB_TOKEN` (Opsional)**:
   - Jika diisi, tool akan memanggil GitLab API untuk mengecek apakah ada file `sonar-project.properties` di root repository. Jika properti `sonar.projectKey=...` ditemukan, key tersebut yang akan dipakai untuk query Sonar.
   - Jika tidak diisi atau file properties tidak ada, tool langsung lanjut ke langkah mapping path.
4. **`NAMESPACE_PREFIX_MAP` (Opsional)**:
   - Digunakan untuk mentranslasi struktur path repo menjadi component key Sonar.
   - Format: `source_path/=target_prefix/` (contoh: `mygroup/project/backend/=sb/be/`).
   - Setelah prefix diganti, seluruh karakter slash (`/`) otomatis dikonversi menjadi colon (`:`).
5. **`GITLAB_BASE_URL` (Opsional)**:
   - Tidak perlu diisi jika MR URL yang kamu masukkan adalah URL lengkap (misal `https://gitlab.mycorp.com/...`), karena host akan langsung terdeteksi.

---

### 📋 Contoh Konfigurasi

#### Opsi 1: Setup Lengkap (Direkomendasikan untuk Tim)
```json
"env": {
  "SONAR_TOKEN": "sqp_your_personal_token",
  "SONAR_HOST_URL": "https://sonarcloud.io",
  "GITLAB_TOKEN": "glpat_your_gitlab_token",
  "NAMESPACE_PREFIX_MAP": "mygroup/project/backend/=sb/be/"
}
```

#### Opsi 2: Setup Minimal (Hanya SonarCloud)
```json
"env": {
  "SONAR_TOKEN": "sqp_your_personal_token"
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

> *"Cek sonar issue untuk MR ini: https://gitlab.com/mygroup/backend/order-service/-/merge_requests/45"*

> *"Fetch all blocker and critical Sonar issues for MR https://gitlab.com/.../-/merge_requests/12 and suggest code fixes."*

> *"Periksa apakah Quality Gate lolos untuk MR ini sebelum kita merge."*

---

## 📄 License

[MIT](LICENSE)
