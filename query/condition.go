package query

import "fmt"

type Condition interface {
	PropertyName() string
	Operator() Operator
	Value() interface{}
	IsTrue(props ResourceProperties) bool
}

type Operator int

const (
	Not Operator = iota - 1
	Equals
	LessThan
	GreaterThan
	LessOrEqualThan
	GreaterOrEqualThan
	Like
)

// note they are in order of their length.
var operators = [...]string{
	"!",
	"=",
	"<",
	">",
	"<=",
	">=",
	"like",
}

func (o Operator) String() string {
	i := int(o + 1)
	if i < 0 || i >= len(operators) {
		return ""
	}
	return operators[i]
}

type condition struct {
	name string
	op   Operator
	val  interface{}
}

func (c condition) String() string {
	return fmt.Sprintf("%s %s %v", c.name, c.op, c.val)
}

func (c condition) PropertyName() string {
	return c.name
}

func (c condition) Operator() Operator {
	return c.op
}

func (c condition) Value() interface{} {
	return c.val
}

func (c condition) IsTrue(props ResourceProperties) bool {
	//TODO implement me
	panic("implement me")
}
