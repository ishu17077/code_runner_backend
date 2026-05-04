package models

import (
	"time"
)

type ExecResult struct {
	Test_id       string  `json:"test_id"`
	Status        *Status `json:"status"`
	ExecResult_id string  `json:"exec_result_id"`
}

type Status struct {
	Message         string    `json:"message"`
	Current_status  string    `json:"current_status"`
	Stdout          string    `json:"stdout"`
	Stderr          string    `json:"stderr"`
	Exec_time_ms    uint16    `json:"exec_time_ms"`
	Stdin           string    `json:"stdin"`
	Expected_output string    `json:"expected_output"`
	Completed_At    time.Time `json:"completed_at"`
}
