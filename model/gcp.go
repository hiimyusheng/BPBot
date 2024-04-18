package model

type Gcp struct {
	Version  string `json:"version"`
	Incident struct {
		ProjectId  string `json:"scoping_project_id"`
		PolicyName string `json:"policy_name"`
		State      string `json:"state"`
		Summary    string `json:"summary"`
		Started    int64  `json:"started_at"`
	} `json:"incident"`
}
type GcpEvent struct {
	Event string
}
