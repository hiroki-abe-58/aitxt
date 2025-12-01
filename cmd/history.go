package cmd

import (
	"fmt"
	"strings"

	"github.com/hiroki-abe-58/aitxt/pkg/history"
	"github.com/hiroki-abe-58/aitxt/pkg/output"
	"github.com/spf13/cobra"
)

var (
	historyLimit  int
	historyOutput string
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Manage command history",
	Long: `View and manage aitxt command history.

Subcommands:
  list    - List recent history entries
  search  - Search history by keyword
  show    - Show details of a specific entry
  clear   - Clear all history
  export  - Export history to file
  stats   - Show history statistics

Examples:
  aitxt history list
  aitxt history list --limit 20
  aitxt history search "translate"
  aitxt history show <id>
  aitxt history clear
  aitxt history export history.md --format md
  aitxt history stats`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runHistoryList(cmd, args)
	},
}

var historyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent history entries",
	RunE:  runHistoryList,
}

var historySearchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "Search history by keyword",
	Args:  cobra.ExactArgs(1),
	RunE:  runHistorySearch,
}

var historyShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show details of a specific entry",
	Args:  cobra.ExactArgs(1),
	RunE:  runHistoryShow,
}

var historyClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all history",
	RunE:  runHistoryClear,
}

var historyExportCmd = &cobra.Command{
	Use:   "export [filename]",
	Short: "Export history to file",
	Args:  cobra.ExactArgs(1),
	RunE:  runHistoryExport,
}

var historyStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show history statistics",
	RunE:  runHistoryStats,
}

var exportFormat string

func init() {
	rootCmd.AddCommand(historyCmd)

	historyCmd.AddCommand(historyListCmd)
	historyCmd.AddCommand(historySearchCmd)
	historyCmd.AddCommand(historyShowCmd)
	historyCmd.AddCommand(historyClearCmd)
	historyCmd.AddCommand(historyExportCmd)
	historyCmd.AddCommand(historyStatsCmd)

	historyCmd.PersistentFlags().IntVarP(&historyLimit, "limit", "n", 10, "Number of entries to show")
	historyCmd.PersistentFlags().StringVarP(&historyOutput, "output", "o", "text", "Output format (text, json, yaml)")

	historyExportCmd.Flags().StringVarP(&exportFormat, "format", "f", "md", "Export format (json, md, text)")
}

func runHistoryList(cmd *cobra.Command, args []string) error {
	store, err := history.NewStore()
	if err != nil {
		return err
	}

	entries, err := store.List(historyLimit)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		fmt.Println("No history entries found.")
		return nil
	}

	formatter := output.NewFormatter(historyOutput)
	if formatter.IsStructured() {
		return formatter.Print(entries)
	}

	fmt.Println("📜 Recent History")
	fmt.Println(strings.Repeat("─", 60))
	for _, entry := range entries {
		printHistoryEntry(entry, false)
	}

	return nil
}

func runHistorySearch(cmd *cobra.Command, args []string) error {
	store, err := history.NewStore()
	if err != nil {
		return err
	}

	entries, err := store.Search(args[0], historyLimit)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		fmt.Printf("No entries found matching '%s'\n", args[0])
		return nil
	}

	formatter := output.NewFormatter(historyOutput)
	if formatter.IsStructured() {
		return formatter.Print(entries)
	}

	fmt.Printf("🔍 Search results for '%s'\n", args[0])
	fmt.Println(strings.Repeat("─", 60))
	for _, entry := range entries {
		printHistoryEntry(entry, false)
	}

	return nil
}

func runHistoryShow(cmd *cobra.Command, args []string) error {
	store, err := history.NewStore()
	if err != nil {
		return err
	}

	entry, err := store.Get(args[0])
	if err != nil {
		return fmt.Errorf("entry not found: %s", args[0])
	}

	formatter := output.NewFormatter(historyOutput)
	if formatter.IsStructured() {
		return formatter.Print(entry)
	}

	printHistoryEntry(entry, true)
	return nil
}

func runHistoryClear(cmd *cobra.Command, args []string) error {
	store, err := history.NewStore()
	if err != nil {
		return err
	}

	if err := store.Clear(); err != nil {
		return err
	}

	fmt.Println("✅ History cleared")
	return nil
}

func runHistoryExport(cmd *cobra.Command, args []string) error {
	store, err := history.NewStore()
	if err != nil {
		return err
	}

	if err := store.Export(args[0], exportFormat); err != nil {
		return err
	}

	fmt.Printf("✅ History exported to %s\n", args[0])
	return nil
}

func runHistoryStats(cmd *cobra.Command, args []string) error {
	store, err := history.NewStore()
	if err != nil {
		return err
	}

	stats, err := store.GetStats()
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(historyOutput)
	if formatter.IsStructured() {
		return formatter.Print(stats)
	}

	fmt.Println("📊 History Statistics")
	fmt.Println(strings.Repeat("─", 40))
	fmt.Printf("Total entries: %d\n", stats["total"])
	fmt.Printf("Total tokens:  %d\n", stats["total_tokens"])
	fmt.Println()

	if byCommand, ok := stats["by_command"].(map[string]int); ok {
		fmt.Println("By Command:")
		for cmd, count := range byCommand {
			fmt.Printf("  %-12s %d\n", cmd, count)
		}
	}
	fmt.Println()

	if byProvider, ok := stats["by_provider"].(map[string]int); ok {
		fmt.Println("By Provider:")
		for provider, count := range byProvider {
			fmt.Printf("  %-12s %d\n", provider, count)
		}
	}

	return nil
}

func printHistoryEntry(entry *history.Entry, detailed bool) {
	fmt.Printf("\033[36m%s\033[0m [%s] \033[33m%s\033[0m (%s)\n",
		entry.ID[:8],
		entry.Timestamp.Format("01/02 15:04"),
		entry.Command,
		entry.Provider)

	if detailed {
		fmt.Println()
		fmt.Println("\033[32mInput:\033[0m")
		fmt.Println(entry.Input)
		fmt.Println()
		fmt.Println("\033[32mOutput:\033[0m")
		fmt.Println(entry.Output)
		if entry.Tokens > 0 {
			fmt.Printf("\n[Tokens: %d]\n", entry.Tokens)
		}
	} else {
		input := strings.ReplaceAll(entry.Input, "\n", " ")
		if len(input) > 60 {
			input = input[:60] + "..."
		}
		fmt.Printf("  %s\n", input)
	}
	fmt.Println()
}
