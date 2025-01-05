package model

type UserProgressResponse struct {
	Error        *string       `json:"error,omitempty"`
	UserProgress *UserProgress `json:"userProgress,omitempty"`
}
