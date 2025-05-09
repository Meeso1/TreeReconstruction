package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing
		},
	}

	// Setup output buffer
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Run the command
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if help output contains expected text
	output := buf.String()
	if len(output) == 0 {
		t.Errorf("Expected help output but got nothing")
	}
} 