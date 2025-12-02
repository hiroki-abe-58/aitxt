package main

import (
	"os"

	"github.com/hiroki-abe-58/aitxt/cmd"
	"github.com/hiroki-abe-58/aitxt/pkg/alias"
)

func main() {
	// Resolve aliases before executing command
	args := os.Args[1:]
	if len(args) > 0 {
		store, err := alias.NewStore()
		if err == nil {
			resolved := store.Resolve(args)
			if len(resolved) != len(args) || resolved[0] != args[0] {
				os.Args = append([]string{os.Args[0]}, resolved...)
			}
		}
	}

	cmd.Execute(cmd.Version)
}
