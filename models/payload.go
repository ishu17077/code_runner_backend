package models

type Payload struct {
	Submission Submission `json:"submission" validate:"required"`
	TestCases  []TestCase `json:"test_cases" validate:"required"`
}
