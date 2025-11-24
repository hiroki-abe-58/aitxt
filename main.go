package main

import (
	"fmt"
	"os"

	"github.com/hiroki-abe-58/aitxt/cmd"
)

var (
	version = "0.1.0"
)

func main() {
	if err := cmd.Execute(version); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
