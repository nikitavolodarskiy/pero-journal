package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:           "list",
	Short:         "List all journal entries",
	Long:          `Print every journal entry, newest first, annotating today's entry if it exists.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, err := buildService()
		if err != nil {
			return err
		}

		entries, err := svc.List()
		if err != nil {
			return err
		}

		if len(entries) == 0 {
			fmt.Println("No journal entries found.")
			return nil
		}

		today := time.Now().Format("2006-01-02")
		for _, e := range entries {
			ds := e.DateString()
			if ds == today {
				fmt.Printf("%s  (today)\n", ds)
			} else {
				fmt.Println(ds)
			}
		}

		return nil
	},
}
