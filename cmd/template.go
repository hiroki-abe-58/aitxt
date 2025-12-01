package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hiroki-abe-58/aitxt/pkg/output"
	"github.com/hiroki-abe-58/aitxt/pkg/template"
	"github.com/spf13/cobra"
)

var templateOutput string

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage prompt templates",
	Long: `Manage reusable prompt templates.

Templates allow you to save and reuse custom prompts with variables.
Variables use Go template syntax: {{.VariableName}}

Subcommands:
  list    - List all templates
  show    - Show template details
  add     - Add a new template
  use     - Use a template
  delete  - Delete a template
  export  - Export templates to file
  import  - Import templates from file
  init    - Initialize with built-in templates

Examples:
  aitxt template list
  aitxt template show code-review-security
  aitxt template add --name my-template --prompt "Review: {{.Code}}"
  aitxt template use my-template --var Code="func main() {}"
  aitxt template export templates.json
  aitxt template import templates.json
  aitxt template init`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTemplateList(cmd, args)
	},
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all templates",
	RunE:  runTemplateList,
}

var templateShowCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show template details",
	Args:  cobra.ExactArgs(1),
	RunE:  runTemplateShow,
}

var (
	addName        string
	addDescription string
	addCommand     string
	addSystemMsg   string
	addPrompt      string
)

var templateAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new template",
	RunE:  runTemplateAdd,
}

var templateUseCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Use a template",
	Args:  cobra.ExactArgs(1),
	RunE:  runTemplateUse,
}

var useVars []string

var templateDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a template",
	Args:  cobra.ExactArgs(1),
	RunE:  runTemplateDelete,
}

var templateExportCmd = &cobra.Command{
	Use:   "export [filename]",
	Short: "Export templates to file",
	Args:  cobra.ExactArgs(1),
	RunE:  runTemplateExport,
}

var templateImportCmd = &cobra.Command{
	Use:   "import [filename]",
	Short: "Import templates from file",
	Args:  cobra.ExactArgs(1),
	RunE:  runTemplateImport,
}

var templateInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize with built-in templates",
	RunE:  runTemplateInit,
}

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateShowCmd)
	templateCmd.AddCommand(templateAddCmd)
	templateCmd.AddCommand(templateUseCmd)
	templateCmd.AddCommand(templateDeleteCmd)
	templateCmd.AddCommand(templateExportCmd)
	templateCmd.AddCommand(templateImportCmd)
	templateCmd.AddCommand(templateInitCmd)

	templateCmd.PersistentFlags().StringVarP(&templateOutput, "output", "o", "text", "Output format (text, json, yaml)")

	templateAddCmd.Flags().StringVarP(&addName, "name", "n", "", "Template name (required)")
	templateAddCmd.Flags().StringVarP(&addDescription, "description", "d", "", "Template description")
	templateAddCmd.Flags().StringVarP(&addCommand, "command", "c", "", "Associated command")
	templateAddCmd.Flags().StringVarP(&addSystemMsg, "system", "s", "", "System message")
	templateAddCmd.Flags().StringVarP(&addPrompt, "prompt", "p", "", "Prompt template")
	templateAddCmd.MarkFlagRequired("name")

	templateUseCmd.Flags().StringArrayVar(&useVars, "var", []string{}, "Variables (format: Key=Value)")
}

func runTemplateList(cmd *cobra.Command, args []string) error {
	store, err := template.NewStore()
	if err != nil {
		return err
	}

	templates, err := store.List()
	if err != nil {
		return err
	}

	if len(templates) == 0 {
		fmt.Println("No templates found. Run 'aitxt template init' to add built-in templates.")
		return nil
	}

	formatter := output.NewFormatter(templateOutput)
	if formatter.IsStructured() {
		return formatter.Print(templates)
	}

	fmt.Println("📝 Templates")
	fmt.Println(strings.Repeat("─", 60))
	for _, tmpl := range templates {
		fmt.Printf("\033[36m%s\033[0m", tmpl.Name)
		if tmpl.Command != "" {
			fmt.Printf(" [%s]", tmpl.Command)
		}
		fmt.Println()
		if tmpl.Description != "" {
			fmt.Printf("  %s\n", tmpl.Description)
		}
		if len(tmpl.Variables) > 0 {
			fmt.Printf("  Variables: %s\n", strings.Join(tmpl.Variables, ", "))
		}
		fmt.Println()
	}

	return nil
}

func runTemplateShow(cmd *cobra.Command, args []string) error {
	store, err := template.NewStore()
	if err != nil {
		return err
	}

	tmpl, err := store.Get(args[0])
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(templateOutput)
	if formatter.IsStructured() {
		return formatter.Print(tmpl)
	}

	fmt.Printf("\033[36m%s\033[0m\n", tmpl.Name)
	fmt.Println(strings.Repeat("─", 40))
	if tmpl.Description != "" {
		fmt.Printf("Description: %s\n", tmpl.Description)
	}
	if tmpl.Command != "" {
		fmt.Printf("Command:     %s\n", tmpl.Command)
	}
	fmt.Println()

	if tmpl.SystemMsg != "" {
		fmt.Println("\033[32mSystem Message:\033[0m")
		fmt.Println(tmpl.SystemMsg)
		fmt.Println()
	}

	fmt.Println("\033[32mPrompt:\033[0m")
	fmt.Println(tmpl.Prompt)
	fmt.Println()

	if len(tmpl.Variables) > 0 {
		fmt.Printf("\033[32mVariables:\033[0m %s\n", strings.Join(tmpl.Variables, ", "))
	}

	if len(tmpl.Defaults) > 0 {
		fmt.Println("\033[32mDefaults:\033[0m")
		for k, v := range tmpl.Defaults {
			fmt.Printf("  %s = %s\n", k, v)
		}
	}

	return nil
}

func runTemplateAdd(cmd *cobra.Command, args []string) error {
	store, err := template.NewStore()
	if err != nil {
		return err
	}

	prompt := addPrompt
	if prompt == "" {
		// Interactive prompt input
		fmt.Println("Enter prompt template (end with empty line):")
		reader := bufio.NewReader(os.Stdin)
		var lines []string
		for {
			line, _ := reader.ReadString('\n')
			line = strings.TrimRight(line, "\n\r")
			if line == "" {
				break
			}
			lines = append(lines, line)
		}
		prompt = strings.Join(lines, "\n")
	}

	tmpl := &template.Template{
		Name:        addName,
		Description: addDescription,
		Command:     addCommand,
		SystemMsg:   addSystemMsg,
		Prompt:      prompt,
	}

	if err := store.Add(tmpl); err != nil {
		return err
	}

	fmt.Printf("✅ Template '%s' added\n", addName)
	if len(tmpl.Variables) > 0 {
		fmt.Printf("   Variables: %s\n", strings.Join(tmpl.Variables, ", "))
	}

	return nil
}

func runTemplateUse(cmd *cobra.Command, args []string) error {
	store, err := template.NewStore()
	if err != nil {
		return err
	}

	// Parse variables
	vars := make(map[string]string)
	for _, v := range useVars {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}

	prompt, systemMsg, err := store.Render(args[0], vars)
	if err != nil {
		return err
	}

	fmt.Println("\033[32mSystem Message:\033[0m")
	fmt.Println(systemMsg)
	fmt.Println()
	fmt.Println("\033[32mPrompt:\033[0m")
	fmt.Println(prompt)
	fmt.Println()
	fmt.Println("\033[33mTip: Pipe this to 'aitxt ask' or copy to use\033[0m")

	return nil
}

func runTemplateDelete(cmd *cobra.Command, args []string) error {
	store, err := template.NewStore()
	if err != nil {
		return err
	}

	if err := store.Delete(args[0]); err != nil {
		return err
	}

	fmt.Printf("✅ Template '%s' deleted\n", args[0])
	return nil
}

func runTemplateExport(cmd *cobra.Command, args []string) error {
	store, err := template.NewStore()
	if err != nil {
		return err
	}

	if err := store.Export(args[0]); err != nil {
		return err
	}

	fmt.Printf("✅ Templates exported to %s\n", args[0])
	return nil
}

func runTemplateImport(cmd *cobra.Command, args []string) error {
	store, err := template.NewStore()
	if err != nil {
		return err
	}

	count, err := store.Import(args[0])
	if err != nil {
		return err
	}

	fmt.Printf("✅ Imported %d templates\n", count)
	return nil
}

func runTemplateInit(cmd *cobra.Command, args []string) error {
	store, err := template.NewStore()
	if err != nil {
		return err
	}

	builtins := template.GetBuiltinTemplates()
	count := 0
	for _, tmpl := range builtins {
		if err := store.Add(tmpl); err == nil {
			count++
		}
	}

	fmt.Printf("✅ Initialized %d built-in templates\n", count)
	fmt.Println("\nAvailable templates:")
	for _, tmpl := range builtins {
		fmt.Printf("  - %s: %s\n", tmpl.Name, tmpl.Description)
	}

	return nil
}
