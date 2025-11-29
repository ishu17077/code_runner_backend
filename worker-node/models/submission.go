package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Submission struct {
	ID            bson.ObjectID `bson:"_id"`
	Team_id       string        `json:"team_id"`
	ProblemID     string        `json:"problem_id" binding:"required"`
	Language      string        `json:"language" binding:"required"`
	Status        string        `json:"status"`
	Code          string        `json:"code" binding:"required"`
	Submission_id string        `json:"submission_id"`
	Tests         []TestCase    `json:"tests" binding:"required"`
}
