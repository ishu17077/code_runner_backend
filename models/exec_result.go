package models

import (
	"time"
)

type ExecResult struct {
	Problem_id    string  `json:"problem_id" binding:"required"` //? if we wanna change just do bson:"problem_id"
	Test_id       string  `json:"test_id"`
	Status        *Status `json:"status"`
	ExecResult_id string  `json:"exec_result_id"`
}

type Status struct {
	Message        string    `json:"message"`
	Current_status string    `json:"current_status"`
	Stdout         string    `json:"stdout"`
	Stderr         string    `json:"stderr"`
	Completed_At   time.Time `json:"completed_at"`
}
