package models

type Attachment struct {
	FileName    string `json:"file_name" validate:"required,min=3,max=50" binding:"required,min=3,max=50"`
	FileContent string `json:"file_content" validate:"required,min=1,max=7000000" binding:"required,min=3,max=7000000"`
}
