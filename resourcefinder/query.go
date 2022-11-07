package resourcefinder

import "fmt"

type Query map[string]QueryValue

type QueryValue interface {
	fmt.Stringer
}

type stringQueryValue struct {
	Value string
}

func (q stringQueryValue) String() string {
	return q.Value
}

func NewQueryValue(s string) QueryValue {

}
