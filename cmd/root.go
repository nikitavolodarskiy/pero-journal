package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"pero/internal/config"
	"pero/internal/service"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pero",
	Short: "A personal journalling CLI",
	Long: `Pero is a minimal journalling app for the terminal.

Running 'pero' with no arguments opens today's journal entry in your editor,
creating it (and any missing directories) if it doesn't exist yet.

Journal entries are stored under:
  <journal_dir>/<year>/<month>/<yyyy-mm-dd>.md

Configuration lives at: ~/.config/pero/config.toml`,

	SilenceErrors: true,
	SilenceUsage:  true,

	RunE: func(cmd *cobra.Command, args []string) error {
		svc, err := buildService()
		if err != nil {
			return err
		}

		entry, err := svc.Today()
		if err != nil {
			return err
		}

		return openInEditor(svc.Editor(), entry.Path)
	},
}

// Execute is the entry point called from main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "pero:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(statsCmd)
}

// buildService loads config and constructs the service.
// All cmd handlers call this — config is never loaded more than once per command.
func buildService() (service.Service, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return service.New(cfg), nil
}

// openInEditor launches editor with path, connecting stdin/stdout/stderr so
// interactive editors (e.g. nvim) work correctly.
// This stays in the cmd layer: launching a process is a UI concern, not a service concern.
func openInEditor(editor, path string) error {
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor %q exited with error: %w", editor, err)
	}
	return nil
}
