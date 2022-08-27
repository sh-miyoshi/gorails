package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build docker images",
	Long:  `Build docker images`,
	Run: func(cmd *cobra.Command, args []string) {
		/*
			案
				Dockerfileはnew時にコピーしておく
				引数でallinoneかprd用を選択
				docker build

				newコマンドのgo-mod-pathをproject-pathとかにする
				それをファイルに書き込んでおく
				image名は<project-path>/<project-name>:latestとか
		*/
	},
}
