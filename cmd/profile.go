package cmd

import (
	"fmt"

	"github.com/pwuersch/ggt/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(profileCmd)

	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileAddCmd)
	profileCmd.AddCommand(profileRemoveCmd)
	profileCmd.AddCommand(profilePathCmd)
}

var profileCmd = &cobra.Command{
	Use:     "profile",
	Short:   "Manages profiles",
	Aliases: []string{"profiles", "settings", "p", "s"},
}

var profileListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Lists all saved profiles",
	Example: "ggt profiles ls",
	Args:    cobra.NoArgs,
	Aliases: []string{"ls", "l"},
	Run: func(cmd *cobra.Command, args []string) {
		profiles, err := lib.ReadProfiles()
		if err != nil {
			lib.ExitWithError(err)
		}

		for index, profile := range profiles {
			fmt.Printf("%d: %+v\n", index+1, profile)
		}
	},
}

var profileAddCmd = &cobra.Command{
	Use:     "add",
	Short:   "Adds a profile",
	Example: "ggt profiles add",
	Args:    cobra.NoArgs,
	Aliases: []string{"a"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Enter the required fields for your new profile")
		newProfile, err := lib.PromptNewProfile()
		if err != nil {
			lib.ExitWithError(err)
		}

		err = lib.AddProfile(newProfile)
		if err != nil {
			lib.ExitWithError(err)
		}
	},
}

var profileRemoveCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Removes a profile",
	Example: "ggt profiles rm",
	Args:    cobra.NoArgs,
	Aliases: []string{"rm"},
	Run: func(cmd *cobra.Command, args []string) {
		profiles, err := lib.ReadProfiles()
		if err != nil {
			lib.ExitWithError(err)
		}

		_, index, err := lib.PromptSelectProfile("Select the profile you want to delete")
		if err != nil {
			lib.ExitWithError(err)
		}

		err = lib.RemoveProfile(profiles, index)
		if err != nil {
			lib.ExitWithError(err)
		}
	},
}

var profilePathCmd = &cobra.Command{
	Use:     "path",
	Short:   "Prints the path to the profiles file",
	Example: "ggt profiles path",
	Args:    cobra.NoArgs,
	Aliases: []string{"p"},
	Run: func(cmd *cobra.Command, args []string) {
		profilePath, err := lib.GetProfilePath()
		if err != nil {
			lib.ExitWithError(err)
		}

		fmt.Print(profilePath)
	},
}
