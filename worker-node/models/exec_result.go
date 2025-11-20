package models

import "time"

type ExecResult struct {
	ID         string  `bson:"_id"`
	Problem_id string  `json:"problem_id" binding:"required"` //? if we wanna change just do bson:"problem_id"
	Team_id    string  `json:"team_id" binding:"required"`
	Stderr     string  `json:"stderr"`
	Test_id    string  `json:"test_id"`
	Status     *Status `json:"status"`
}

type Status struct {
	Message        string    `json:"message"`
	Current_status string    `json:"current_status"`
	Stdout         string    `json:"stdout"`
	Stderr         string    `json:"stderr"`
	Completed_At   time.Time `json:"completed_at"`
}
