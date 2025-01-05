package model

type CourseEdge struct {
	Cursor string  `json:"cursor"`
	Node   *Course `json:"node"`
}
