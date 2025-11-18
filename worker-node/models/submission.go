package models

import "time"

type Submission struct {
	ID            string  `bson:"_id"`
	User_id       string  `json:"user_id"`
	ProblemID     string  `json:"problem_id" binding:"required"`
	Language      string  `json:"language" binding:"required"`
	Status        *Status `json:"status"`
	Code          string  `json:"code" binding:"required"`
	Submission_id string  `json:"submission_id"`
}

type Status struct {
	Message        string    `json:"message"`
	Current_status string    `json:"current_status"`
	Stdout         string    `json:"stdout"`
	Stderr         string    `json:"stderr"`
	Completed_At   time.Time `json:"completed_at"`
}
