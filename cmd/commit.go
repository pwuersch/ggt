package cmd

import (
	"fmt"
	"strings"

	"github.com/pwuersch/ggt/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(commitCmd)
}

var commitCmd = &cobra.Command{
	Use:     "commit",
	Short:   "Commits all changes",
	Long:    "Commits all local changes in the current repo and pushes them to the current upstream branch",
	Example: "ggt c add navigation bar",
	Args:    cobra.ArbitraryArgs,
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		message := strings.Join(args, " ")

		var err error
		if message == "" {
			message, err = lib.PromptInput("Enter a commit message", lib.ValidateEmpty)
			if err != nil {
				lib.ExitWithError(err)
			}
		}

		fmt.Println("Adding all files in the current directory")
		err = lib.Add()
		if err != nil {
			lib.ExitWithError(err)
		}

		fmt.Println("Committing all changes")
		err = lib.Commit(message)
		if err != nil {
			lib.ExitWithError(err)
		}

		fmt.Println("Pushing created commit")
		err = lib.Push()
		if err != nil {
			lib.ExitWithError(err)
		}
	},
}
