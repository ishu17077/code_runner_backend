package models

type TestCase struct {
	// ID             string `bson:"_id"`
	Problem_id     string `json:"problem_id" binding:"required"`
	Is_public      bool   `json:"is_public" binding:"required"`
	Stdin          string `json:"stdin" binding:"required"`
	ExpectedOutput string `json:"expectedOutput" binding:"required"`
	Test_id        string `json:"test_id" binding:"required"`
}
