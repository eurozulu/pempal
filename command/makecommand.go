package command

import (
	"context"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"pempal/builders"
	"pempal/finder"
	"pempal/pemtypes"
	"pempal/templates"
	"strings"
)

type makeCommand struct {
	VerboseOutput bool `flag:"verbose,v"`
}

func (mc makeCommand) Run(ctx context.Context, args Arguments, out io.Writer) error {
	params := args.Parameters()
	if len(params) == 0 {
		return fmt.Errorf("must specify at least the pem type or template to make.  <pem type> [template[ template...]]")
	}
	pt := pemtypes.ParsePEMType(params[0])

	builder := mc.prepareBuilder(ctx, pt)
	if builder == nil {
		return fmt.Errorf("%s is not a know type which can be built", params[0])
	}
	t, err := mc.buildTemplate(ctx, pt, params[1:])
	if err != nil {
		return err
	}
	if err := builder.AddTemplate(t); err != nil {
		return err
	}
	errs := builder.Validate()
	for len(errs) > 0 {
		if err = mc.populateBuilder(builder, errs); err != nil {
			return err
		}
		errs = builder.Validate()
	}
}

func (mc makeCommand) prepareBuilder(ctx context.Context, pemType pemtypes.PEMType) builders.Builder {
	builder := builders.NewBuilder(pemType)
	if builder == nil {
		return nil
	}

	blocks, err := mc.collectStdInResource(ctx)
	if err != nil && mc.VerboseOutput {
		log.Println(err)
	}
	if len(blocks) > 0 {
		builder.AddResource(blocks...)
	}
}

func (mc makeCommand) collectStdInResource(ctx context.Context) ([]*pem.Block, error) {
	pf := finder.NewPemFinderResources(false, mc.VerboseOutput)
	var blocks []*pem.Block
	for _, l := range pf.FindAll(ctx, "-") {
		blks, err := finder.ReadLocationPems(l)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, blks...)
	}
	return blocks, nil
}

func (mc makeCommand) buildTemplate(ctx context.Context, pemType pemtypes.PEMType, names []string) (templates.Template, error) {
	if len(names) == 0 {
		return nil, nil
	}
	tf := finder.NewTemplateFinder(true, mc.VerboseOutput)
	locs := tf.FindAll(ctx, names...)
	t := templates.NewTemplate(pemType)
	if mc.VerboseOutput {
		log.Printf("found %d template locations\n", len(locs))
	}
	for _, l := range locs {
		data, err := l.MarshalText()
		if err != nil {
			return nil, fmt.Errorf("failed to read template in %s  %v", l.Path(), err)
		}
		if err = yaml.Unmarshal(data, &t); err != nil {
			return nil, fmt.Errorf("failed to write template in %s  %v", l.Path(), err)
		}
	}
	return t, nil
}

func (mc makeCommand) populateBuilder(builder builders.Builder, errs []error) error {
	m := templates.BlankTemplate{}
	for _, err := range errs {
		if strings.HasPrefix(err.Error(), "missing ") {
			name := err.Error()[len("missing "):]
			v, err := mc.requestProperty(name)
			if err != nil {
				return err
			}
			m[name] = v
		}
	}
	return builder.AddTemplate(m)
}

func (mc makeCommand) requestProperty(name string) (string, error) {
}
