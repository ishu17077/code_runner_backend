package models

type Submission struct {
	ID        string     `json:"id"`
	ProblemID string     `json:"problem_id"`
	Language  string     `json:"language" binding:"required" validate:"required"`
	Code      string     `json:"code" binding:"required" validate:"required"`
	Tests     []TestCase `json:"tests" binding:"required" validate:"required"`
}
