package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"pempal/logger"
	"pempal/templates"
	"sort"
)

var ResourceTemplatesStore templates.TemplateStore

type TemplatesCommand struct {
	Name   string `flag:"name"`
	Add    string `flag:"add"`
	Remove bool   `flag:"remove"`
}

func (cmd TemplatesCommand) Execute(args []string, out io.Writer) error {
	ResourceTemplatesStore = ResourceTemplates.(templates.TemplateStore)
	if cmd.Remove {
		if err := cmd.removeTemplates(args); err != nil {
			return err
		}
	}
	if cmd.Add != "" {
		if err := cmd.addTemplate(args); err != nil {
			return err
		}
		logger.Log(logger.Info, "added new template: ")
	}
	names := ResourceTemplatesStore.Names(args...)
	sort.Strings(names)
	for _, n := range names {
		fmt.Fprintln(out, n)
	}
	return nil
}

func (cmd TemplatesCommand) addTemplate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("must provide a unique name for the new template")
	}
	if len(args) > 1 {
		return fmt.Errorf("multiple names can not be specified for a new template")
	}
	var data []byte
	var err error
	if cmd.Add == "-" {
		data, err = readStdIn()
	} else {
		data, err = os.ReadFile(cmd.Add)
	}
	if err != nil {
		return err
	}
	t, err := ResourceTemplates.ParseTemplate(data)
	if err != nil {
		return err
	}
	return ResourceTemplatesStore.SaveTemplate(args[0], t)
}

func (cmd TemplatesCommand) removeTemplates(names []string) error {
	for _, name := range names {
		if err := ResourceTemplatesStore.DeleteTemplate(name); err != nil {
			return err
		}
	}
	return nil
}

func readStdIn() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		buf.Write(scan.Bytes())
	}
	if scan.Err() != nil && scan.Err() != io.EOF {
		return nil, scan.Err()
	}
	return buf.Bytes(), nil
}
