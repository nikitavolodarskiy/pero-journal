package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:           "stats",
	Short:         "Show journal statistics",
	Long:          `Display aggregated statistics about your journal entries.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, err := buildService()
		if err != nil {
			return err
		}

		st, err := svc.Stats()
		if err != nil {
			return err
		}

		if st.TotalEntries == 0 {
			fmt.Println("No journal entries yet. Run 'pero' to create your first one.")
			return nil
		}

		fmt.Printf("Entries        %d\n", st.TotalEntries)
		fmt.Printf("Words          %d\n", st.TotalWords)
		fmt.Printf("Avg words      %d\n", st.AverageWords)
		fmt.Printf("Streak         %d days\n", st.CurrentStreak)
		fmt.Printf("Longest streak %d days\n", st.LongestStreak)

		if st.FirstEntry != nil {
			fmt.Printf("First entry    %s\n", st.FirstEntry.Format("2006-01-02"))
		}
		if st.LastEntry != nil {
			fmt.Printf("Last entry     %s\n", st.LastEntry.Format("2006-01-02"))
		}

		return nil
	},
}
