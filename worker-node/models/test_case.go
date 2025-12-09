package models

import "go.mongodb.org/mongo-driver/v2/bson"

type TestCase struct {
	ID             bson.ObjectID `bson:"_id"`
	Problem_id     string        `json:"problem_id" binding:"required" validate:"required"`
	Is_public      bool          `json:"is_public" binding:"required" validate:"required"`
	Stdin          string        `json:"stdin" binding:"required" validate:"required"`
	ExpectedOutput string        `json:"expected_output" binding:"required" validate:"required"`
	Test_id        string        `json:"test_id" binding:"required" validate:"required"`
}
