package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Launch the interactive TUI browser",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !isTerminal(os.Stdin) || !isTerminal(os.Stdout) {
			return fmt.Errorf("the ui command requires an interactive terminal")
		}
		return runTUI()
	},
}

// maybeRunTUI launches the TUI when stdout and stdin are both terminals.
// It reports whether the TUI ran.
func maybeRunTUI() bool {
	if !isTerminal(os.Stdin) || !isTerminal(os.Stdout) {
		return false
	}
	if err := runTUI(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
	}
	return true
}
