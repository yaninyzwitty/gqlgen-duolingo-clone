package model

type Course struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	ImageSrc string  `json:"imageSrc"`
	Units    []*Unit `json:"units,omitempty"`
}
