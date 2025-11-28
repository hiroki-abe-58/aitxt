package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// These variables are set at build time using ldflags
	Version   = "0.1.0"
	GitCommit = "unknown"
	BuildDate = "unknown"
	GoVersion = runtime.Version()
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long: `Show detailed version and build information.

Examples:
  aitxt version
  aitxt version --short`,
	RunE: runVersion,
}

var versionShort bool

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&versionShort, "short", "s", false, "Show short version only")
}

func runVersion(cmd *cobra.Command, args []string) error {
	if versionShort {
		fmt.Printf("aitxt %s\n", Version)
		return nil
	}

	fmt.Println("🤖 aitxt - AI-powered text processing CLI")
	fmt.Println()
	fmt.Println("Version Information:")
	fmt.Printf("  Version:    %s\n", Version)
	fmt.Printf("  Git Commit: %s\n", GitCommit)
	fmt.Printf("  Build Date: %s\n", BuildDate)
	fmt.Println()
	fmt.Println("Runtime Information:")
	fmt.Printf("  Go Version: %s\n", GoVersion)
	fmt.Printf("  OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("  CPUs:       %d\n", runtime.NumCPU())
	fmt.Println()
	fmt.Println("Supported Providers:")
	fmt.Println("  - OpenAI (GPT-4o, GPT-4, GPT-3.5)")
	fmt.Println("  - Anthropic (Claude Sonnet, Claude Opus)")
	fmt.Println("  - Google (Gemini 1.5 Flash, Gemini Pro)")
	fmt.Println()
	fmt.Println("Repository:")
	fmt.Println("  https://github.com/hiroki-abe-58/aitxt")
	fmt.Println()
	fmt.Println("License: MIT")

	return nil
}
