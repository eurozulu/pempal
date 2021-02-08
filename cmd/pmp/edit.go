package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pempal/pempal"
	"github.com/pempal/templates"
	"gopkg.in/yaml.v3"
	"strings"
)

func ConfirmTemplate(prompt string, t templates.Template) error {
	lines, err := templateLines(t)
	if err != nil {
		return err
	}

	fmt.Println(prompt)
	for i, l := range lines {
		fmt.Printf("%d) %s\n", i+1, l)
	}
	fmt.Println("\n0) Abort")

	for {
		fmt.Println("Confirm these details by hitting enter or", strings.ToLower(prompt))
		fmt.Println("Select 1-%d to edit that property", len(lines))

		for {
			index := pempal.PromptInputNumber(len(lines))
			if index == 0 {
				return fmt.Errorf("aborted")
			}
			if index < 0 {
				return nil
			}
			if index > len(lines) {
				fmt.Println("Select 1-%d to edit that property", len(lines))
				continue
			}
			if err := EditTemplate(t, lines[index]); err != nil {
				return err
			}
		}
	}
}

func templateLines(t templates.Template) ([]string, error) {
	by, err := yaml.Marshal(t)
	if err != nil {
		return nil, err
	}

	var lines []string
	scn := bufio.NewScanner(bytes.NewReader(by))
	for scn.Scan() {
		if !strings.Contains(scn.Text(), ":") {
			continue
		}
		lines = append(lines, scn.Text())
	}
	return lines, nil
}

func ChooseTemplate(prompt string, tps []*pempal.QueryResult, options []string) int {
	ListPems(tps, true, false, true)
	for i, s := range options {
		fmt.Printf("%d)\t%s", len(tps)+i, s)
	}
	fmt.Println("0)\texit")
	fmt.Println()
	total := len(tps) + len(options)

	for {
		fmt.Println("\nSelect 1-%d %s", total, prompt)
		i := pempal.PromptInputNumber(total)
		if i < 0 {
			return -1
		}
		if i < 0 || i >= total {
			continue
		}
		return i
	}
}

func EditTemplate(t templates.Template, line string) error {
	panic("Not yet implemented")
}
