package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

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
		// check required program
		c := exec.Command("go", "version")
		out, err := c.Output()
		if err != nil {
			fmt.Printf("Failed to get go version: %+v\n", err)
			if len(out) > 0 {
				fmt.Println(string(out))
			}
			os.Exit(1)
		}

		c = exec.Command("npm", "--version")
		out, err = c.Output()
		if err != nil {
			fmt.Printf("Failed to get npm version: %+v\n", err)
			if len(out) > 0 {
				fmt.Println(string(out))
			}
			os.Exit(1)
		}
	})
}

// Execute method run root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}
