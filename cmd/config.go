package cmd

import (
	"github.com/pwuersch/ggt/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "Configures a repository",
	Long:    "Clones a local git repository with the settings of a provided profile",
	Example: "ggt config",
	Args:    cobra.NoArgs,
	Aliases: []string{"configure", "cfg"},
	Run: func(cmd *cobra.Command, args []string) {
		profile, _, err := lib.PromptSelectProfile("Select a profile to configure the current repo with")
		if err != nil {
			lib.ExitWithError(err)
		}

		err = lib.Config("", profile)
		if err != nil {
			lib.ExitWithError(err)
		}
	},
}
