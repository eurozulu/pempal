package commands

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/utils"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

type ConfigCommand struct {
	Output io.Writer
}

func (cmd ConfigCommand) Exec(args ...string) error {
	if cmd.Output == nil {
		cmd.Output = os.Stdout
	}
	buf := bytes.NewBuffer(nil)
	cmd.writeConfigFilePaths(buf)

	cfg := config.Config
	if err := yaml.NewEncoder(buf).Encode(cfg); err != nil {
		return err
	}
	if _, err := buf.WriteTo(cmd.Output); err != nil {
		return err
	}
	return nil
}

func (cmd ConfigCommand) writeConfigFilePaths(buf *bytes.Buffer) error {
	globalPath := config.GlobalFilePath()
	gfExists := utils.FileExists(globalPath)
	if gfExists {
		fmt.Fprintf(cmd.Output, "Global config: %s\n", globalPath)
	} else {
		fmt.Fprintf(cmd.Output, "No Global config found in %s\n", globalPath)
	}
	localPath := config.Config.LocalFilePath()
	lfExists := utils.FileExists(localPath)
	if lfExists {
		fmt.Fprintf(cmd.Output, "Local config: %s\n", localPath)
	} else {
		fmt.Fprintf(cmd.Output, "No local config found in %s\n", localPath)
	}
	return nil
}
