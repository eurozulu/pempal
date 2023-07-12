package commands

import (
	"fmt"
	"github.com/eurozulu/pempal/builders"
	"github.com/eurozulu/pempal/commandline/valueeditors"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/utils"
	"io"
)

type newKeyCommand struct {
	flagDTO resources.PrivateKeyDTO
}

var keyPropertyEditors = []valueeditors.ValueEditor{
	valueeditors.StringEditor{
		PropertyName:  "key-algorithm",
		Choice:        utils.PublicKeyAlgorithms,
		DefaultChoice: 0,
	},
	valueeditors.NumberEditor{
		PropertyName:  "key-length",
		Choice:        []int{512, 1024, 2048, 4096},
		DefaultChoice: 2,
		AllowedInput:  true,
	},
	valueeditors.StringEditor{
		PropertyName:  "key-curve",
		Choice:        utils.ECDSACurveNames[1:],
		DefaultChoice: 3,
	},
	valueeditors.BoolEditor{
		PropertyName:  "is-encrypted",
		DefaultChoice: true,
	},
}

func (cmd newKeyCommand) Execute(args []string, out io.Writer) error {
	builder, err := builders.NewBuilder(resources.PrivateKey)

	// check if any named templates or file resources given in args
	temps, err := argumentsToTemplates(args)
	if err != nil {
		return err
	}
	builder.AddTemplate(temps...)

	// get flag properties as template
	t, err := resources.DTOToTemplate(&cmd.flagDTO)
	if err != nil {
		return err
	}
	builder.AddTemplate(t)

	if err = confirmBuild("Create new key", keyPropertyEditors, builder); err != nil {
		return err
	}

	kr, err := builder.Build()
	if err != nil {
		return err
	}

	// flip new keyresource into a PEM string
	dto, err := resources.NewResourceDTO(kr)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, dto.String())
	return err
}

func (cmd newKeyCommand) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&cmd.flagDTO)
}

func (cmd newKeyCommand) MarshalYAML() (interface{}, error) {
	return &cmd.flagDTO, nil
}
