package cmd

import (
	"fmt"

	"pero/internal/config"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise pero configuration",
	Long: `Create the pero config directory and a default config.toml if they
don't already exist.

Default config file location: ~/.config/pero/config.toml`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := config.Init()
		if err != nil {
			return err
		}
		fmt.Printf("Config ready at: %s\n", path)
		return nil
	},
}
