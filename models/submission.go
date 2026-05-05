package models

import (
	"errors"
	"regexp"

	"github.com/ishu17077/code_runner_backend/models/enums/language"
)

type Submission struct {
	ID          string       `json:"id"`
	Problem_id  string       `json:"problem_id"`
	Language    string       `json:"language" binding:"required" validate:"required,min=3,max=10"`
	Code        string       `json:"code" binding:"required" validate:"required,min=3,max=7000000"`
	Attachments []Attachment `json:"attachments"`
	Pipe        bool         `json:"pipe"`
	Tests       []TestCase   `json:"tests"`
}

var validBase64Regex = regexp.MustCompile(`^[a-zA-Z0-9+/]*={0,2}$`)

const maxCodeLength = 50000

func ValidateSubmission(submission *Submission) error {
	if submission.Language == "" {
		return errors.New("language is required")
	}

	lang := language.LanguageParser(submission.Language)
	if lang == language.Undefined {
		return errors.New("unsupported language: " + submission.Language)
	}

	if submission.Code == "" {
		return errors.New("code payload is required")
	}

	if len(submission.Code) > maxCodeLength {
		return errors.New("code payload exceeds maximum allowed size")
	}

	if !validBase64Regex.MatchString(submission.Code) {
		return errors.New("code payload is not a valid base64 string")
	}

	return nil
}
