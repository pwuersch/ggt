package cmd

import (
	"fmt"

	"github.com/pwuersch/ggt/lib"
	"github.com/pwuersch/ggt/lib/api"
	"github.com/pwuersch/ggt/lib/globals"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cloneCmd)

	cloneCmd.Flags().StringVarP(&globals.CloneRemoteUrl, "url", "u", "", "git repository url (ssh or http)")
	cloneCmd.Flags().StringVarP(&globals.CloneDestDir, "dir", "d", "", "destination folder of the cloned repo")
}

var cloneCmd = &cobra.Command{
	Use:     "clone",
	Short:   "Clones a repository",
	Long:    "Clones a git repository with optional profiles",
	Example: "ggt cl -u https://github.com/pwuersch/ggt.git",
	Args:    cobra.NoArgs,
	Aliases: []string{"cl"},
	Run: func(cmd *cobra.Command, args []string) {
		profile, _, err := lib.PromptSelectProfile("Select a profile to use")
		if err != nil {
			lib.ExitWithError(err)
		}

		url := globals.CloneRemoteUrl
		if url == "" {
			if profile.Api.Adapter == "" {
				url, err = lib.PromptInput("Enter the git url (ssh or http[s])", lib.ValidateGitUrl)
				if err != nil {
					lib.ExitWithError(err)
				}
			} else {
				adapter, err := api.GetAdapterByName(profile.Api.Adapter)
				if err != nil {
					lib.ExitWithError(err)
				}

				adapterInfo := adapter.Info()

				search := ""
				if adapterInfo.ProvidesSearch {
					search, err = lib.PromptInput("Search for (optional)", nil)
					if err != nil {
						lib.ExitWithError(err)
					}
				}

				scope := profile.Api.DefaultScope
				if adapterInfo.ProvidesScope {
					scope, err = lib.PromptInput(fmt.Sprintf("Select the scope to use (default: \"%s\")", scope), nil)
					if err != nil {
						lib.ExitWithError(err)
					}
				}

				projects, err := adapter.GetProjects(profile, search, scope)
				if err != nil {
					lib.ExitWithError(err)
				}

				projectNames := globals.Map(projects, func(project api.Project) string {
					return project.DisplayName
				})

				index, _, err := lib.SelectItem("Select the project you want to clone", projectNames)
				if err != nil {
					lib.ExitWithError(err)
				}
				url = projects[index].SshUrl
			}
		} else {
			err = lib.ValidateGitUrl(url)
			if err != nil {
				lib.ExitWithError(err)
			}
		}

		dest := globals.CloneDestDir
		if dest == "" {
			defaultDest := lib.ParseGitDestDir(url)
			dest, err = lib.PromptInput(fmt.Sprintf("Enter the destination for your repo (default: \"%s\")", defaultDest), nil)
			if err != nil {
				lib.ExitWithError(err)
			}

			if dest == "" {
				dest = defaultDest
			}
		}

		err = lib.Clone(url, dest, profile)
		if err != nil {
			lib.ExitWithError(err)
		}
	},
}
