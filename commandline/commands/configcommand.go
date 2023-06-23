package commands

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/utils"
	"github.com/go-yaml/yaml"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ConfigCommand displays the current configuration settings.
// given with no parameters, it displays all the config settings.
// one or more parameters can optionally name specific properties and only those properties are displayed.
// placing a '=' directly after the property name (no space) will set the string following the '=' as that property value.
// If the current config is using an existing config file, that file is updated.
// If there is no current config file, one will be created.
type ConfigCommand struct {
	cfg config.Config
}

// Execute config command
func (cmd ConfigCommand) Execute(args []string, out io.Writer) error {
	cfg, err := config.CurrentConfig()
	if err != nil {
		return err
	}
	cmd.cfg = cfg

	if a, err := cleanArguments(args); err != nil {
		return err
	} else {
		args = a
	}

	assignments, names := parseArgsForAssignments(args)

	if err = cmd.setConfigValues(out, assignments); err != nil {
		return err
	}
	names, err = cleanAndCheckArgNames(names)
	if err != nil {
		return err
	}

	if !CommonFlags.Quiet {
		if err = cmd.writeHeader(out); err != nil {
			return err
		}
	}

	if err = cmd.writeProperties(out, names); err != nil {
		return err
	}

	return nil
}

func (cmd *ConfigCommand) setConfigValues(out io.Writer, assignments []string) error {
	if len(assignments) == 0 {
		return nil
	}
	m, err := parseAssignmentsToMap(assignments)
	if err != nil {
		return err
	}
	if err := cmd.mapIntoConfig(m); err != nil {
		return err
	}
	if !CommonFlags.Quiet {
		logger.Info("updated config with %d assignment", len(assignments))
	}

	// Save the new config
	if err = cmd.saveCurrentConfig(); err != nil {
		return err
	}
	if !CommonFlags.Quiet {
		logger.Info("saved config file to %s", cmd.cfg.ConfigLocation())
	}
	return nil
}

func (cmd ConfigCommand) writeHeader(out io.Writer) error {
	if _, err := fmt.Fprintln(out, "Pempal Configuration:"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(out, "config-path: "); err != nil {
		return err
	}
	if cmd.cfg.ConfigLocation() != "" {
		if _, err := fmt.Fprint(out, cmd.cfg.ConfigLocation()); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprint(out, "-- not set --"); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(out); err != nil {
		return err
	}
	return nil
}

func (cmd ConfigCommand) writeProperties(out io.Writer, names []string) error {
	colOut := createColumnOutput(out)
	mcfg, err := cmd.configAsMap()
	if err != nil {
		return fmt.Errorf("failed to read config as yaml  %v", err)
	}

	buf := bytes.NewBuffer(nil)
	for _, name := range names {
		buf.Reset()
		if !CommonFlags.Quiet {
			buf.WriteString(name)
			buf.WriteRune(':')
		}
		buf.WriteString(mcfg[name])
		if _, err := colOut.Write(buf.Bytes()); err != nil {
			return err
		}
		if err := colOut.WriteString("\n"); err != nil {
			return err
		}
	}
	return nil
}

func (cmd *ConfigCommand) mapIntoConfig(m map[string]string) error {
	by, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(by, cmd.cfg); err != nil {
		return err
	}
	return nil
}

func (cmd ConfigCommand) configAsMap() (map[string]string, error) {
	by, err := yaml.Marshal(cmd.cfg)
	if err != nil {
		return nil, err
	}
	m := map[string]string{}
	if err = yaml.Unmarshal(by, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (cmd ConfigCommand) saveCurrentConfig() error {
	p := cmd.cfg.ConfigLocation()
	if p == "" {
		p = filepath.Join(os.ExpandEnv("$PWD"), ".config")
	}
	return config.SaveConfig(cmd.cfg)
}

func cleanAndCheckArgNames(args []string) ([]string, error) {
	// If no name(s) given, return all names
	cfgNames := config.ConfigNames[:]
	if len(args) == 0 {
		return cfgNames, nil
	}

	var names []string
	for _, arg := range args {
		if name := containsString(arg, cfgNames); name != "" {
			names = append(names, name)
		} else {
			return nil, fmt.Errorf("property name '%s' is not known", arg)
		}
	}
	return names, nil
}

func parseAssignmentsToMap(assignments []string) (map[string]string, error) {
	// build a map of the assignements and apply it to the config
	m := map[string]string{}
	for _, arg := range assignments {
		ss := strings.SplitN(arg, "=", 2)
		key := strings.TrimSpace(ss[0])
		var val string
		if len(ss) > 1 {
			val = strings.TrimSpace(ss[1])
		}
		key = containsString(key, config.ConfigNames[:])
		if key == "" {
			return nil, fmt.Errorf("%s is not a known configuration name", ss[0])
		}
		m[key] = val
	}
	return m, nil
}

func parseArgsForAssignments(args []string) (assignments, names []string) {
	for i := 0; i < len(args); i++ {
		arg := strings.TrimSpace(args[i])
		if strings.Contains(arg, "=") {
			argz := strings.SplitN(arg, "=", 2)
			if len(argz) < 2 || argz[1] == "" {
				// No value found, combine with next argument
				if i+1 < len(args) {
					i++
					argz[1] = args[i]
				}
			}
			assignments = append(assignments, strings.Join(argz, "="))
			arg = argz[0]
		}
		names = append(names, arg)
	}
	return assignments, names
}

func containsString(s string, ss []string) string {
	for _, sz := range ss {
		if strings.EqualFold(sz, s) {
			return sz
		}
	}
	return ""
}

func createColumnOutput(out io.Writer) *utils.ColumnOutput {
	colOut := utils.NewColumnOutput(out)
	colOut.DataDelimiter = ":"
	colOut.ColumnWidths = []int{16}
	colOut.ColumnDelimiter = "  "
	return colOut
}
