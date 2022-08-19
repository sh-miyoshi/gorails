package cmd

import (
	"os"

	"github.com/sh-miyoshi/gorails/pkg/cmd/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(clientCmd)
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Run development client",
	Long:  `Run development client`,
	Run: func(cmd *cobra.Command, args []string) {
		os.Chdir("client")
		util.RunCommand("npm", "start")
	},
}
