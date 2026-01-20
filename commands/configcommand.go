package commands

import (
	"bytes"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/logging"
	"gopkg.in/yaml.v2"
)

// ConfigCommand displays the current configuration details
// @Command(config)
type ConfigCommand struct{}

// ShowConfig shows the working root path, search path and the various resource related directory paths.
// Config can be set using a file name '.ppconfig' in the pempal root directory.
// @Action
func (c *ConfigCommand) ShowConfig() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("Pempal Config:\n")
	buf.WriteString("root path: ")
	buf.WriteString(config.RootPath())
	buf.WriteString("\n")

	buf.WriteString("search path: ")
	buf.WriteString(config.SearchPath())
	buf.WriteString("\n")

	buf.WriteString("config file: ")
	cp := config.ConfigPath()
	if cp == "" {
		cp = "- no config file in use -"
	}
	buf.WriteString(cp)
	buf.WriteString("\n")

	data, err := yaml.Marshal(config.DefaultPPConfig)
	if err != nil {
		logging.Error("failed to marshal config %v", err)
	}
	buf.Write(data)
	buf.WriteString("\n\n")
	return buf.String()
}
