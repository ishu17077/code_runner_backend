package models

type TestCase struct {
	Problem_id     string `json:"problem_id"`
	Stdin          string `json:"stdin" binding:"required" validate:"required"`
	ExpectedOutput string `json:"expected_output" binding:"required" validate:"required"`
	Test_id        string `json:"test_id"`
}
