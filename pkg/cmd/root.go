package cmd

import (
	"fmt"
	"os"

	"github.com/sh-miyoshi/gorails/pkg/config"
	"github.com/spf13/cobra"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "gorails",
	Short: "gorails is a building tool for web apps",
	Long:  "gorails is a building tool for web apps",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cobra.OnInitialize(func() {
		if err := config.Init(configFile); err != nil {
			fmt.Printf("Failed to init config: %+v", err)
			os.Exit(1)
		}

		// TODO check required program check(go, npm)
	})

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "config.yaml", "The file path of gorails config")
}

// Execute method run root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}
