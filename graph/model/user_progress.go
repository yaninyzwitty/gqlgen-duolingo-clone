package model

type UserProgress struct {
	UserID       string  `json:"userId"`
	UserName     string  `json:"userName"`
	ActiveCourse *Course `json:"activeCourse,omitempty"`
	Hearts       int32   `json:"hearts"`
	Points       int32   `json:"points"`
}
