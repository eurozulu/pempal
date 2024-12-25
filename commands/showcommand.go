package commands

import (
	"context"
	"fmt"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"sort"
)

type ShowCommand struct {
	Output io.Writer
	Types  []model.ResourceType
}

func (cmd ShowCommand) Exec(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide at least one file or directory path ")
	}
	if cmd.Output == nil {
		cmd.Output = os.Stdout
	}

	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	scan := resources.NewPemScan(cmd.Types...)
	for pr := range scan.ScanPath(ctx, args...) {
		if err := cmd.showPems(pr); err != nil {
			return err
		}
	}
	return nil
}

func (cmd ShowCommand) showPems(resource resources.PemResource) error {
	pems := resource.Content
	sort.Slice(pems, func(i, j int) bool {
		return pems[i].Type < pems[j].Type
	})

	for _, pem := range pems {
		t, err := templates.TemplateOfPem(pem)
		if err != nil {
			return err
		}
		by, err := yaml.Marshal(t)
		if err != nil {
			return err
		}
		_, err = cmd.Output.Write(by)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.Output)

		if !CommonFlags.Quiet {
			rt := model.ParseResourceTypeFromPEMType(pem.Type)
			fmt.Fprintln(cmd.Output, rt.String())
		}
		fmt.Fprintln(cmd.Output, "\n")

	}
	return nil
}
