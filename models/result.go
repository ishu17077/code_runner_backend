package models

type Result struct {
	Problem_id string       `json:"problem_id"`
	Status     string       `json:"status"`
	Results    []ExecResult `json:"results"`
	Error      string       `json:"error,omitempty"`
}
