package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version 版本号
	Version string
	Date    string
	Commit  string
	Branch  string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if Version == "" {
			Version = "Unknow"
		}
		if Commit == "" {
			Commit = "Unknow"
		}
		fmt.Printf("Version: %s\n", Version)
		if Date != "" {
			fmt.Printf("Date: %s\n", Date)
		}
		if Commit != "" {
			fmt.Printf("Commit: %s\n", Commit)
		}
		if Branch != "" {
			fmt.Printf("Branch: %s\n", Branch)
		}
		os.Exit(1)
	},
}
