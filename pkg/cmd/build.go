package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/sh-miyoshi/gorails/pkg/cmd/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/stew/slice"
)

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().String("target", "all", "build target [all | separate]")
	buildCmd.Flags().String("tag", "latest", "docker image tag")
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build docker images",
	Long:  `Build docker images`,
	Run: func(cmd *cobra.Command, args []string) {
		targets := []string{"all", "separate"}
		target, _ := cmd.Flags().GetString("target")
		if !slice.Contains(targets, target) {
			fmt.Printf("Invalid target type %s was specified\n", target)
			os.Exit(1)
		}
		clientExists := util.FileExists("client")
		if target == "all" && !clientExists {
			fmt.Printf("[Warning] cannot find client directory, so will create server image only")
			target = "separate"
		}
		project, err := os.ReadFile(".gorails-project")
		if err != nil {
			fmt.Printf("Failed to read .gorails-project file: %+v\n", err)
			os.Exit(1)
		}
		tag, _ := cmd.Flags().GetString("tag")
		image := strings.Trim(string(project), "\n") + ":" + tag

		switch target {
		case "all":
			util.RunCommand("docker", "build", "-f", "build/Dockerfile.all", "-t", image, ".")
		case "separate":
			panic("TODO: WIP")
			// TODO build server, build client if exists
		}
	},
}
