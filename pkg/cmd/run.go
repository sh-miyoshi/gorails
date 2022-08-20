package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.AddCommand(serverCmd)
	runCmd.AddCommand(clientCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run development server or client",
	Long:  `Run development server or client`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
