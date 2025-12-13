package models

type Payload struct {
	Class_name string     `json:"class_name" validate:"required"`
	Exec_time  int        `json:"exec_time" validate:"required"`
	Tests      []TestCase `json:"tests" validate:"required"`
}
