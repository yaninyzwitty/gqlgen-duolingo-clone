package model

type Unit struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Course      *Course `json:"course"`
	Order       int32   `json:"order"`
}
