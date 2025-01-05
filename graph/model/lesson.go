package model

type Lesson struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Unit  *Unit  `json:"unit"`
	Order int32  `json:"order"`
}
