package model

type CourseConnection struct {
	Edges    []*CourseEdge `json:"edges,omitempty"`
	PageInfo *PageInfo     `json:"pageInfo"`
}
