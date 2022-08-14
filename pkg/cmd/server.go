package cmd

import (
	"github.com/sh-miyoshi/gorails/pkg/cmd/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run development server",
	Long:  `Run development server`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO linuxの場合
		util.RunCommand("go", "build", "-o", "a.exe")
		util.RunCommand("./a.exe")

		// TODO hot reloading
	},
}
