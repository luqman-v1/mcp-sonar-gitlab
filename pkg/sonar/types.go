package sonar

type SonarIssue struct {
	Key       string          `json:"key"`
	Rule      string          `json:"rule"`
	Severity  string          `json:"severity"`
	Component string          `json:"component"`
	Project   string          `json:"project"`
	Line      int             `json:"line"`
	Hash      string          `json:"hash"`
	TextRange *SonarTextRange `json:"textRange,omitempty"`
	Message   string          `json:"message"`
	Type      string          `json:"type"`
	Status    string          `json:"status"`
	Effort    string          `json:"effort,omitempty"`
}

type SonarTextRange struct {
	StartLine   int `json:"startLine"`
	EndLine     int `json:"endLine"`
	StartOffset int `json:"startOffset"`
	EndOffset   int `json:"endOffset"`
}

type SonarQualityGate struct {
	Status     string                  `json:"status"` // OK, ERROR, WARN
	Conditions []SonarQualityCondition `json:"conditions"`
}

type SonarQualityCondition struct {
	Status         string `json:"status"` // OK, ERROR, WARN
	MetricKey      string `json:"metricKey"`
	Comparator     string `json:"comparator"`
	ErrorThreshold string `json:"errorThreshold"`
	ActualValue    string `json:"actualValue"`
}

type SonarSearchResponse struct {
	Total  int          `json:"total"`
	Issues []SonarIssue `json:"issues"`
}

type sonarQualityGateResponse struct {
	ProjectStatus SonarQualityGate `json:"projectStatus"`
}
