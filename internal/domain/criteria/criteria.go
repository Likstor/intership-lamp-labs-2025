package criteria

import "fmt"

type Operator string

var (
	GT Operator = ">"
	LT Operator = "<"
	EQ Operator = "="
)

type Condition interface {
	ToSQL() string
}

type Criteria struct {
	Condition Condition
}

type SimpleCondition struct {
	Field    string
	Operator Operator
	Value    any
}

func (sc *SimpleCondition) ToSQL() string {
	return fmt.Sprintf("WHERE %s %s %v", sc.Field, sc.Operator, sc.Value)
}
