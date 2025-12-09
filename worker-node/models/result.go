package models

type Result struct {
	Status  string       `json:"status"`
	Results []ExecResult `json:"results"`
}
