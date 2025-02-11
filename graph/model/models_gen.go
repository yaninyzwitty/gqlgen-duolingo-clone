// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Mutation struct {
}

type PageInfo struct {
	EndCursor   string `json:"endCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}

type Query struct {
}

type ChallengeType string

const (
	ChallengeTypeSelect ChallengeType = "SELECT"
	ChallengeTypeAssist ChallengeType = "ASSIST"
)

var AllChallengeType = []ChallengeType{
	ChallengeTypeSelect,
	ChallengeTypeAssist,
}

func (e ChallengeType) IsValid() bool {
	switch e {
	case ChallengeTypeSelect, ChallengeTypeAssist:
		return true
	}
	return false
}

func (e ChallengeType) String() string {
	return string(e)
}

func (e *ChallengeType) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ChallengeType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ChallengeType", str)
	}
	return nil
}

func (e ChallengeType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
