package cmd

import (
	"fmt"
	"os"

	"github.com/cosmtrek/air/runner"
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
		cfg, err := runner.InitConfig("config/hot_reloader.toml")
		if err != nil {
			fmt.Printf("Failed to initialize hot reloader: %+v", err)
			os.Exit(1)
		}
		r, err := runner.NewEngineWithConfig(cfg, true)
		if err != nil {
			fmt.Printf("Failed to create hot reloader engine: %+v", err)
			os.Exit(1)
		}
		defer func() {
			r.Stop()
		}()

		r.Run()
	},
}
