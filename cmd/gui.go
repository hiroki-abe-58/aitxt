package cmd

import (
	"github.com/hiroki-abe-58/aitxt/pkg/gui"
	"github.com/spf13/cobra"
)

var (
	guiPort int
)

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Start web-based GUI",
	Long: `Start a web-based graphical user interface for aitxt.

The GUI provides an easy-to-use interface for all aitxt features:
  - Ask: Ask AI any question
  - Translate: Translate text between languages
  - Summarize: Summarize long texts
  - Proofread: Check grammar and spelling
  - Style: Transform writing style
  - Explain: Explain error messages
  - Review: Review code for issues
  - Document: Generate documentation
  - Chat: Interactive conversation mode

The GUI runs on a local web server and opens in your default browser.

Examples:
  aitxt gui
  aitxt gui --port 3000`,
	RunE: runGUI,
}

func init() {
	rootCmd.AddCommand(guiCmd)

	guiCmd.Flags().IntVarP(&guiPort, "port", "p", 8080, "Port to run the GUI server on")
}

func runGUI(cmd *cobra.Command, args []string) error {
	return gui.Start(guiPort)
}
