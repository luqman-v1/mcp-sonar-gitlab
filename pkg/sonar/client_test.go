package sonar_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luqman-v1/mcp-sonar-gitlab/pkg/sonar"
)

func TestSonarClient_GetMRIssues_BearerAuth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer valid-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if r.URL.Path != "/api/issues/search" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"total": 1,
			"issues": [
				{
					"key": "ISSUE-1",
					"rule": "go:S1186",
					"severity": "CRITICAL",
					"component": "acme:be:trade-service:pkg/handler.go",
					"project": "acme:be:trade-service",
					"line": 42,
					"message": "Add a nested comment explaining why this method is empty",
					"type": "CODE_SMELL",
					"status": "OPEN"
				}
			]
		}`))
	}))
	defer ts.Close()

	client := sonar.NewClient(ts.URL, "valid-token", "")
	issues, err := client.GetMRIssues(context.Background(), "acme:be:trade-service", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "CRITICAL" {
		t.Errorf("expected CRITICAL severity, got %s", issues[0].Severity)
	}
	if issues[0].Line != 42 {
		t.Errorf("expected line 42, got %d", issues[0].Line)
	}
}

func TestSonarClient_GetMRIssues_BasicAuthFallback(t *testing.T) {
	attempt := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt++
		auth := r.Header.Get("Authorization")

		// First attempt with Bearer fails (simulating SonarQube on-premise)
		if auth == "Bearer onprem-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Second attempt should use Basic Auth (username = token, password = "")
		user, pass, ok := r.BasicAuth()
		if ok && user == "onprem-token" && pass == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"total": 0, "issues": []}`))
			return
		}

		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	client := sonar.NewClient(ts.URL, "onprem-token", "")
	issues, err := client.GetMRIssues(context.Background(), "my-component", 5)
	if err != nil {
		t.Fatalf("expected fallback to succeed, got error: %v", err)
	}

	if len(issues) != 0 {
		t.Errorf("expected 0 issues, got %d", len(issues))
	}
	if attempt != 2 {
		t.Errorf("expected 2 attempts (Bearer then Basic), got %d", attempt)
	}
}

func TestSonarClient_GetQualityGate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"projectStatus": {
				"status": "ERROR",
				"conditions": [
					{
						"status": "ERROR",
						"metricKey": "new_coverage",
						"comparator": "LT",
						"errorThreshold": "80",
						"actualValue": "72.5"
					}
				]
			}
		}`))
	}))
	defer ts.Close()

	client := sonar.NewClient(ts.URL, "tok", "")
	qg, err := client.GetQualityGate(context.Background(), "comp", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if qg.Status != "ERROR" {
		t.Errorf("expected status ERROR, got %s", qg.Status)
	}
	if len(qg.Conditions) != 1 {
		t.Fatalf("expected 1 condition, got %d", len(qg.Conditions))
	}
	if qg.Conditions[0].MetricKey != "new_coverage" {
		t.Errorf("expected metric new_coverage, got %s", qg.Conditions[0].MetricKey)
	}
}
