package cmd

import (
	"fmt"
	"strings"

	"github.com/hiroki-abe-58/aitxt/pkg/alias"
	"github.com/hiroki-abe-58/aitxt/pkg/output"
	"github.com/spf13/cobra"
)

var aliasOutput string

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Manage command aliases",
	Long: `Manage command aliases for quick shortcuts.

Aliases allow you to create shortcuts for frequently used commands.

Subcommands:
  list   - List all aliases
  add    - Add a new alias
  delete - Delete an alias
  init   - Initialize with suggested aliases

Examples:
  aitxt alias list
  aitxt alias add s summarize
  aitxt alias add tj "translate --to ja" --desc "Translate to Japanese"
  aitxt alias delete s
  aitxt alias init

Usage after creating alias:
  aitxt s document.txt       # Same as: aitxt summarize document.txt
  aitxt tj "Hello world"     # Same as: aitxt translate --to ja "Hello world"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAliasList(cmd, args)
	},
}

var aliasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all aliases",
	RunE:  runAliasList,
}

var (
	aliasDesc string
)

var aliasAddCmd = &cobra.Command{
	Use:   "add [name] [command]",
	Short: "Add a new alias",
	Args:  cobra.MinimumNArgs(2),
	RunE:  runAliasAdd,
}

var aliasDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete an alias",
	Args:  cobra.ExactArgs(1),
	RunE:  runAliasDelete,
}

var aliasInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize with suggested aliases",
	RunE:  runAliasInit,
}

func init() {
	rootCmd.AddCommand(aliasCmd)

	aliasCmd.AddCommand(aliasListCmd)
	aliasCmd.AddCommand(aliasAddCmd)
	aliasCmd.AddCommand(aliasDeleteCmd)
	aliasCmd.AddCommand(aliasInitCmd)

	aliasCmd.PersistentFlags().StringVarP(&aliasOutput, "output", "o", "text", "Output format (text, json, yaml)")
	aliasAddCmd.Flags().StringVarP(&aliasDesc, "desc", "d", "", "Alias description")
}

func runAliasList(cmd *cobra.Command, args []string) error {
	store, err := alias.NewStore()
	if err != nil {
		return err
	}

	aliases := store.List()

	if len(aliases) == 0 {
		fmt.Println("No aliases defined. Run 'aitxt alias init' to add suggested aliases.")
		return nil
	}

	formatter := output.NewFormatter(aliasOutput)
	if formatter.IsStructured() {
		return formatter.Print(aliases)
	}

	fmt.Println("⚡ Aliases")
	fmt.Println(strings.Repeat("─", 50))
	for _, a := range aliases {
		fmt.Printf("\033[36m%-8s\033[0m → %s", a.Name, a.Command)
		if a.Desc != "" {
			fmt.Printf("  \033[90m# %s\033[0m", a.Desc)
		}
		fmt.Println()
	}
	fmt.Println()
	fmt.Println("\033[33mUsage: aitxt <alias> [args...]\033[0m")

	return nil
}

func runAliasAdd(cmd *cobra.Command, args []string) error {
	store, err := alias.NewStore()
	if err != nil {
		return err
	}

	name := args[0]
	command := strings.Join(args[1:], " ")

	if err := store.Add(name, command, aliasDesc); err != nil {
		return err
	}

	fmt.Printf("✅ Alias added: %s → %s\n", name, command)
	return nil
}

func runAliasDelete(cmd *cobra.Command, args []string) error {
	store, err := alias.NewStore()
	if err != nil {
		return err
	}

	if err := store.Delete(args[0]); err != nil {
		return err
	}

	fmt.Printf("✅ Alias '%s' deleted\n", args[0])
	return nil
}

func runAliasInit(cmd *cobra.Command, args []string) error {
	store, err := alias.NewStore()
	if err != nil {
		return err
	}

	builtins := alias.GetBuiltinAliases()
	count := 0
	for _, a := range builtins {
		if err := store.Add(a.Name, a.Command, a.Desc); err == nil {
			count++
		}
	}

	fmt.Printf("✅ Initialized %d aliases\n\n", count)
	fmt.Println("Available aliases:")
	for _, a := range builtins {
		fmt.Printf("  \033[36m%-4s\033[0m → %-25s \033[90m# %s\033[0m\n", a.Name, a.Command, a.Desc)
	}

	return nil
}
