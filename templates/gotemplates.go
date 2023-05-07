package templates

import (
	"bytes"
	gotemplate "text/template"
	"time"
)

func containsGoTemplates(text []byte) bool {
	i := bytes.Index(text, []byte("{{"))
	if i < 0 {
		return false
	}
	return bytes.Index(text[i+2:], []byte("}}")) >= 0
}

func executeGoTemplate(text []byte, data map[string]interface{}) ([]byte, error) {
	gt, err := gotemplate.New("template-manager").
		Funcs(buildFuncMap()).
		Parse(string(text))

	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	if err = gt.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func buildFuncMap() gotemplate.FuncMap {
	return gotemplate.FuncMap{
		"now": func() string { return time.Now().Format(time.RFC3339) },
		"nowPlusDays": func(days int) string {
			day := time.Hour * 24
			return time.Now().Add(day * time.Duration(days)).Format(time.RFC3339)
		},
	}
}
