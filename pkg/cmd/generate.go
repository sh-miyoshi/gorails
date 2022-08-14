package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate resources",
	Long:  `Generate resources`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("generate command requires resource type and name")
			fmt.Println("e.g. gorails generate controller user")
			os.Exit(1)
		}

		resType := args[0]
		resName := args[1]
		fmt.Printf("generating resource type: %s\n", resType)
		fmt.Printf("generating resource name: %s\n", resName)

		switch resType {
		case "controller", "c":
			fmt.Println("WIP")
		case "model", "m":
			fmt.Println("WIP")
		case "view", "v":
			fmt.Println("WIP")
		default:
			fmt.Printf("Invalid resource type %s is specified", resType)
			os.Exit(1)
		}
	},
}
