package cmd

import (
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open <date>",
	Short: "Open a specific journal entry",
	Long: `Open a journal entry for the given date in your editor.

The date must be in yyyy-mm-dd format, e.g.:

  pero open 2026-05-21`,
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, err := buildService()
		if err != nil {
			return err
		}

		entry, err := svc.Get(args[0])
		if err != nil {
			return err
		}

		return openInEditor(svc.Editor(), entry.Path)
	},
}
