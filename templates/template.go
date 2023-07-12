package templates

type Template map[string]string

func MergeTemplates(temps ...Template) Template {
	temp := Template{}
	for _, t := range temps {
		MergeTemplate(temp, t)
	}
	return temp
}

func MergeTemplate(dst, src Template) {
	for k, v := range src {
		dst[k] = v
	}
}
