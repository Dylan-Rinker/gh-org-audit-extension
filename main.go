package main

import (
	"os"

	"github.com/dylan-rinker/gh-org-audit-extension/cmd"
)

func main() {

	// Instantiate and execute root command
	cmd := cmd.NewCmd()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
