package model

type Challenge struct {
	ID       string        `json:"id"`
	Lesson   *Lesson       `json:"lesson"`
	Type     ChallengeType `json:"type"`
	Question string        `json:"question"`
	Order    int32         `json:"order"`
}
