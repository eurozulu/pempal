package resourceio

import (
	"bytes"
	"fmt"
	"github.com/go-yaml/yaml"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var queryFuncMap = template.FuncMap{
	"contains": fnContains,
	"compare":  fnCompare,
}

var operatorMaping = map[string]string{
	"=":        "eq .%v %v",
	"!=":       "ne .%v %v",
	"NOT":      "ne .%v %v",
	"<":        "lt .%v %v",
	"<=":       "le .%v %v",
	">":        "gt .%v %v",
	">=":       "ge .%v %v",
	"contains": "contains .%v %v",
	"compare":  "compare .%v %v",
}

type Query interface {
	Fields() []string
	Match(m map[string]interface{}) (map[string]bool, error)
}

type query struct {
	expressions map[string]*template.Template
}

func (q query) Fields() []string {
	flds := make([]string, len(q.expressions))
	var index int
	for k := range q.expressions {
		flds[index] = k
	}
	return flds
}

func (q query) Match(m map[string]interface{}) (map[string]bool, error) {
	buf := bytes.NewBuffer(nil)
	for _, t := range q.expressions {
		if err := t.Execute(buf, cleanKeyNames(m)); err != nil {
			return nil, err
		}
	}
	bm := map[string]bool{}
	if err := yaml.NewDecoder(buf).Decode(&bm); err != nil {
		return nil, err
	}
	return bm, nil
}

func fnContains(op1, op2 string) bool {
	return strings.Contains(op1, op2)
}

func fnCompare(op1, op2 interface{}) bool {
	return strings.EqualFold(valueToString(op1), valueToString(op2))
}

func valueToString(v interface{}) string {
	switch vt := v.(type) {
	case string:
		return vt
	case int, int64, int32, int16, int8:
		return strconv.FormatInt(v.(int64), 10)
	case uint, uint64, uint32, uint16, uint8:
		return strconv.FormatUint(v.(uint64), 10)
	case float64, float32:
		return strconv.FormatFloat(v.(float64), 0, 1, 10)
	case bool:
		return strconv.FormatBool(vt)
	case time.Time:
		return vt.Format(http.TimeFormat)
	default:
		return ""
	}
}
func findOp(s string) string {
	for k := range operatorMaping {
		if strings.Contains(s, k) {
			return k
		}
	}
	return ""
}

func splitByOp(s, op string) (op1, op2 string) {
	i := strings.Index(s, op)
	if i < 0 {
		i = len(s)
	}
	op1 = strings.Replace(strings.TrimSpace(s[:i]), "-", "_", -1)
	op2 = strings.Replace(strings.TrimSpace(s[i+len(op):]), "'", "`", -1)
	return op1, op2
}

func cleanKeyNames(m map[string]interface{}) map[string]interface{} {
	mm := map[string]interface{}{}
	for k, v := range m {
		nk := strings.Replace(k, "-", "_", -1)
		if mv, ok := v.(map[string]interface{}); ok {
			v = cleanKeyNames(mv)
		}
		mm[nk] = v
	}
	return mm
}

func ParseQuery(q string) (Query, error) {
	qy := &query{expressions: map[string]*template.Template{}}
	qs := strings.Split(q, ",")
	for _, v := range qs {
		op := findOp(v)
		if op == "" {
			return nil, fmt.Errorf("no operator found after %s", v)
		}
		op1, op2 := splitByOp(v, op)
		fun := fmt.Sprintf("%s: {{ %s }}\n", op1, fmt.Sprintf(operatorMaping[op], op1, op2))
		t, err := template.New(op1).Funcs(queryFuncMap).Parse(fun)
		if err != nil {
			return nil, err
		}
		qy.expressions[op1] = t
	}
	return qy, nil
}
