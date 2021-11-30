package templates

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path"
	"pempal/pemreader"
	"sort"
	"strings"
	"sync"
)

const FileTag = "#"

const ENV_TemplatePath = "PP_TEMPLATEPATH"

var TemplatePath = strings.TrimSpace(os.Getenv(ENV_TemplatePath))

// TemplateNames lists the known names of all the templates, including buuld in ones
func TemplateNames(includeBuiltIn bool) []string {
	tp := []string{os.ExpandEnv("$PWD")}
	if TemplatePath != "" {
		tp = append(tp, strings.Split(TemplatePath, ":")...)
	}

	names := findAll(tp)
	sort.Strings(names)

	// add built in names
	if includeBuiltIn {
		names = append(names, sortedMapKeys(builtInTemplates)...)
	}
	for i, n := range names {
		names[i] = cleanName(n)
	}
	return names
}

func FindTemplate(name string) (Template, error) {
	key := strings.TrimLeft(name, FileTag)
	if t, ok := builtInTemplates[key]; ok {
		return t, nil
	}

	// not built in, scan for file with that name
	tp := []string{os.ExpandEnv("$PWD")}
	if TemplatePath != "" {
		tp = append(tp, strings.Split(TemplatePath, ":")...)
	}
	p := findFirst(tp, name)
	if p == "" {
		return nil, fmt.Errorf("%s not found", name)
	}
	return parseTemplate(p)
}

func findAll(rootpaths []string) []string {
	fs := &pemreader.FileScanner{Filter: &pemreader.ExtensionFilter{"yaml": true}}
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()

	found := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(rootpaths))
	go func(wg *sync.WaitGroup) {
		wg.Wait()
		defer close(found)
	}(&wg)

	for _, p := range rootpaths {
		pCh := fs.Find(ctx, p)
		go func(ch <-chan string, wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case fp, ok := <-pCh:
					if !ok {
						return
					}
					if !strings.HasPrefix(path.Base(fp), FileTag) {
						continue
					}
					select {
					case <-ctx.Done():
						return
					case found <- fp:
					}
				}
			}
		}(pCh, &wg)
	}

	var names []string
	for {
		select {
		case p, ok := <-found:
			if !ok {
				return names
			}
			names = append(names, p)
		}
	}
}

func cleanName(p string) string {
	n := path.Base(p)
	e := path.Ext(n)
	if len(e) > 0 {
		n = n[:len(n)-len(e)]
	}
	if !strings.HasPrefix(n, FileTag) {
		n = strings.Join([]string{FileTag, n}, "")
	}
	return n
}

func findFirst(rootpaths []string, name string) string {
	fs := &pemreader.FileScanner{Filter: &pemreader.NameFilter{Name: name}}
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()

	found := make(chan string)
	go func() {
		var wg sync.WaitGroup
		wg.Add(len(rootpaths))
		go func() {
			wg.Wait()
			defer close(found)
		}()

		for _, p := range rootpaths {
			pCh := fs.Find(ctx, p)
			go func(ch <-chan string, wg *sync.WaitGroup) {
				defer wg.Done()
				select {
				case <-ctx.Done():
					return
				case fp, ok := <-pCh:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case found <- fp:
					}
					return
				}
			}(pCh, &wg)
		}
	}()
	// return first one found, triggers ctx cancel, killing other searches.
	return <-found
}

func sortedMapKeys(m map[string]Template) []string {
	keys := make([]string, len(m))
	var i int
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func parseTemplate(name string) (Template, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()
	m := map[string]interface{}{}
	if err = yaml.NewDecoder(f).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}
