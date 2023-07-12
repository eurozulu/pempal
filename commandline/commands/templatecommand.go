package commands

//import (
//	"fmt"
//	"github.com/eurozulu/argdecoder"
//	"github.com/eurozulu/pempal/builder"
//	"github.com/eurozulu/pempal/config"
//	"github.com/eurozulu/pempal/logger"
//	"github.com/eurozulu/pempal/templates"
//	"github.com/eurozulu/pempal/utils"
//	"io"
//	"sort"
//	"strings"
//)
//
//// TemplateCommand, when used with no parameters, lists the names of all the available templates.
//type TemplateCommand struct {
//	// ShowExtends, when present and true, will show the templates which are being extended by a named template, preceeding the named template.
//	ShowExtends bool `flag:"show-extends,e,extends"`
//
//	// Apply will format the given name(s) into a single template suitable to pass to make.
//	// Each named template extends any named template preceeding it in the argument list.
//	// If a named template alread extends another template, that extension chain is applied prior to any named template preceeding it is applied.
//	Apply bool `flag:"apply"`
//
//	// AllNames when true, displays the built-in template names as well as any user template names.
//	// This flag only applies when no named templates or property flags are given, i.e. listing names, otherwise it is ignored.
//	AllNames bool `flag:"all,a"`
//
//	// Delete, when true deletes any given template names.
//	// If no names given, it is ignored.
//	// Will confirm deletion from the stdIn unless the -Quiet flag is set
//	Delete bool `flag:"delete"`
//
//	templateManager templates.TemplateManager
//}
//
//func (cmd TemplateCommand) Execute(args []string, out io.Writer) error {
//	// set up the template manager
//	if tm, err := config.TemplateManager(); err == nil {
//		cmd.templateManager = tm
//	} else {
//		return fmt.Errorf("template manager unavailable, %v", err)
//	}
//
//	if logger.DefaultLogger.Level() >= logger.LevelDebug {
//		cfg, err := config.CurrentConfig()
//		if err != nil {
//			return err
//		}
//		logger.Info("Template root location: %s", cfg.Templates())
//	}
//
//	cleanArgs, flags := argdecoder.ParseArgs(args)
//	assignments, names := parseArgsForAssignments(cleanArgs)
//
//	// perform deletes before assignments to allow assignment to replace existing
//	if cmd.Delete {
//		if len(names) > 0 {
//			if err := cmd.deleteTemplates(names); err != nil {
//				return err
//			}
//			names = nil
//		} else {
//			logger.Warning("ignoring delete flag as no names given")
//		}
//	}
//	if len(assignments) > 0 {
//		if err := cmd.addNewAssignments(assignments); err != nil {
//			return err
//		}
//	}
//
//	var flagTemp templates.Template
//	if len(flags) > 0 {
//		t, err := flagTemplate(flags)
//		if err != nil {
//			return err
//		}
//		flagTemp = t
//	}
//
//	// no templates, show the available template names
//	if len(names) == 0 && flagTemp == nil {
//		return cmd.writeTemplateNames(out)
//	}
//
//	if cmd.Apply {
//		return cmd.applyTemplates(out, names, flagTemp)
//	}
//	return cmd.writeTemplates(out, names, flagTemp)
//}
//
//func (cmd TemplateCommand) writeTemplateNames(out io.Writer) error {
//	var names []string
//	if cmd.AllNames {
//		names = cmd.templateManager.AllNames()
//		if !CommonFlags.Quiet {
//			logger.Info("All template names:")
//		}
//	} else {
//		names = cmd.templateManager.Names()
//		if !CommonFlags.Quiet {
//			logger.Info("User template names:")
//		}
//	}
//	if !CommonFlags.Quiet && len(names) == 0 {
//		logger.Info("No templates found")
//		return nil
//	}
//	sort.Strings(names)
//	for _, n := range names {
//		if _, err := fmt.Fprintln(out, n); err != nil {
//			return err
//		}
//	}
//	if !CommonFlags.Quiet {
//		logger.Info("%d templates found", len(names))
//	}
//	return nil
//}
//
//func (cmd TemplateCommand) applyTemplates(out io.Writer, names []string, flagTemplate templates.Template) error {
//	temps, err := cmd.templateManager.ExtendedTemplatesByName(names...)
//	if err != nil {
//		return err
//	}
//	if flagTemplate != nil {
//		temps = append(temps, flagTemplate)
//	}
//	tb := builder.TemplateBuilder(temps)
//	t, err := tb.MergeTemplates()
//	if err != nil {
//		return err
//	}
//	return cmd.writeTemplate(out, "", t, 0)
//}
//
//func (cmd TemplateCommand) writeTemplates(out io.Writer, names []string, flagTemplate templates.Template) error {
//	temps, err := cmd.templateManager.TemplateByName(names...)
//	if err != nil {
//		return err
//	}
//	if flagTemplate != nil {
//		temps = append(temps, flagTemplate)
//	}
//	for i, t := range temps {
//		cmd.writeTemplate(out, templateName(names, i), t, 0)
//		if _, err := out.Write([]byte{'\n'}); err != nil {
//			return err
//		}
//	}
//	if !CommonFlags.Quiet {
//		fmt.Fprintf(out, "\n%d template%s\n", len(temps), plualString(len(temps)))
//	}
//	return nil
//}
//
//func (cmd TemplateCommand) writeTemplate(out io.Writer, name string, t templates.Template, indent int) error {
//	if !CommonFlags.Quiet && name != "" {
//		writeIndentedString(out, indent, fmt.Sprintf("Template Name: %s", name))
//		if _, err := out.Write([]byte{'\n'}); err != nil {
//			return err
//		}
//	}
//	var cOut utils.ColourOut
//	if indent == 0 {
//		cOut = utils.ColourOut{Out: out, Col: utils.ColourBlue}
//	} else {
//		cOut = utils.ColourOut{Out: out, Col: utils.ColourCyan}
//	}
//	if err := writeIndentedString(cOut, indent, t.String()); err != nil {
//		return err
//	}
//	if _, err := cOut.Write([]byte{'\n'}); err != nil {
//		return err
//	}
//
//	if cmd.ShowExtends && t.Tags().ContainsTag(templates.TAG_EXTENDS) {
//		name := t.Tags().TagByName(templates.TAG_EXTENDS).Value()
//		et, err := cmd.templateManager.TemplateByName(name)
//		if err != nil {
//			return err
//		}
//		if err := cmd.writeTemplate(out, name, et[0], indent+1); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (cmd TemplateCommand) deleteTemplates(names []string) error {
//	if !CommonFlags.Quiet {
//		y, err := PromptYorN(fmt.Sprintf("delete templates: %v", names), false)
//		if err != nil {
//			return err
//		}
//		if !y {
//			return fmt.Errorf("abandoned")
//		}
//	}
//	for _, n := range names {
//		if err := cmd.templateManager.DeleteTemplate(n); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (cmd TemplateCommand) addNewAssignments(names []string) error {
//	for _, name := range names {
//		ss := strings.SplitN(name, "=", 2)
//		name = ss[0]
//		var temp templates.Template
//		if len(ss) > 1 {
//			t, err := templates.ParseInlineTemplate(ss[1])
//			if err != nil {
//				return err
//			}
//			temp = t
//		}
//		if CommonFlags.ForceOut && cmd.templateManager.Exists(name) {
//			if err := cmd.templateManager.DeleteTemplate(name); err != nil {
//				return err
//			}
//		}
//		if err := cmd.templateManager.SaveTemplate(name, temp); err != nil {
//			return fmt.Errorf("Failed to save template '%s'  %v", name, err)
//		}
//		if !CommonFlags.Quiet {
//			logger.Info("created template %s\n", name)
//		}
//	}
//	return nil
//}
//
//func flagTemplate(flags utils.FlatMap) (templates.Template, error) {
//	if len(flags) == 0 {
//		return nil, nil
//	}
//	return builder.TemplateFromValue(flags)
//}
//
//func templateName(names []string, index int) string {
//	if index < len(names) {
//		return names[index]
//	}
//	return "--anonymous--"
//}
//
//func plualString(size int) string {
//	if size == 1 {
//		return ""
//	}
//	return "s"
//}
//
//func writeIndentedString(out io.Writer, indent int, s string) error {
//	widths := make([]int, indent+1)
//	for i := range widths {
//		widths[i] = 8
//	}
//	colOut := utils.NewColumnOutput(out, widths...)
//	_, err := colOut.WriteSlice(append(make([]string, indent), s))
//	return err
//}
