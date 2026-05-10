package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version    = "dev"
	buildTime  = ""
)

var rootCmd = &cobra.Command{
	Use:   "windo",
	Short: "windo is a CLI tool for various project templates and operations",
	Long:  `windo is a CLI tool for various project templates and operations, such as generating kratos project templates.`,
}

func Execute() error {
	return rootCmd.Execute()
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of windo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("windo version %s\n", version)
		if buildTime != "" {
			fmt.Printf("build time: %s\n", buildTime)
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(versionCmd)
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func printErr(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

// SilenceUsage is to silence usage when an error occurs.
// Cobra will print usage by default on error.
func SilenceUsage(cmd *cobra.Command) {
	cmd.SilenceUsage = true
}
