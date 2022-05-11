package cmd

import (
	"os"

	"github.com/pwuersch/ggt/lib/globals"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&globals.RootDebug, "debug", false, "Print debugging information")
}

var rootCmd = &cobra.Command{
	Use:   "ggt",
	Short: "Git helper tool",
	Long:  "Useful for managing git profiles and applying them in different scenarios",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
