package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sh-miyoshi/gorails/pkg/cmd/util"
	"github.com/sh-miyoshi/gorails/pkg/templates"
	"github.com/spf13/cobra"
)

type templateValue struct {
	ProjectPath string
	DBName      string
	ServerExt   string
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().Bool("skip-client", false, "Skip installing client")
	newCmd.Flags().String("project-path", "", "[Required] project path. e.g. github.com/sh-miyoshi")
	newCmd.MarkFlagRequired("project-path")
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "New gorails project",
	Long:  `New gorails project`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("new command requires project name")
			fmt.Println("e.g. gorails new sample-project")
			os.Exit(1)
		}

		projectName := args[0]
		fmt.Printf("Project name: %s\n", projectName)

		if err := os.Mkdir(projectName, 0755); err != nil {
			fmt.Println("Failed to create project directory")
			fmt.Printf("Directory %s will not empty\n", projectName)
			os.Exit(1)
		}

		os.Mkdir(fmt.Sprintf("%s/app", projectName), 0755)
		os.Mkdir(fmt.Sprintf("%s/app/controllers", projectName), 0755)
		os.Mkdir(fmt.Sprintf("%s/app/models", projectName), 0755)
		os.Mkdir(fmt.Sprintf("%s/app/schema", projectName), 0755)
		os.Mkdir(fmt.Sprintf("%s/config", projectName), 0755)
		os.Mkdir(fmt.Sprintf("%s/db", projectName), 0755)
		os.Mkdir(fmt.Sprintf("%s/log", projectName), 0755)
		os.Mkdir(fmt.Sprintf("%s/system", projectName), 0755)

		fmt.Println("Successfully created base directories")

		projectPath, _ := cmd.Flags().GetString("project-path")
		projectPath = strings.TrimSuffix(projectPath, "/")
		projectPath += "/" + projectName

		ext := "out"
		if runtime.GOOS == "windows" {
			ext = "exe"
		}

		// Copy system files
		vals := templateValue{
			ProjectPath: projectPath,
			DBName:      "app",
			ServerExt:   ext,
		}
		templates.Exec(templates.ModelBase, fmt.Sprintf("%s/app/models/base.go", projectName), nil)
		templates.Exec(templates.DatabaseYaml, fmt.Sprintf("%s/config/database.yaml", projectName), vals)
		templates.Exec(templates.HotReloader, fmt.Sprintf("%s/config/hot_reloader.toml", projectName), vals)
		templates.Exec(templates.Routes, fmt.Sprintf("%s/config/routes.go", projectName), nil)
		templates.Exec(templates.Migration, fmt.Sprintf("%s/db/migration.go", projectName), vals)
		templates.Exec(templates.DockerCompose, fmt.Sprintf("%s/docker-compose.yml", projectName), vals)
		templates.Exec(templates.GitIgnore, fmt.Sprintf("%s/.gitignore", projectName), nil)
		templates.Exec(templates.MainGo, fmt.Sprintf("%s/main.go", projectName), vals)
		templates.Exec(templates.SystemModel, fmt.Sprintf("%s/system/model.go", projectName), vals)
		templates.Exec(templates.SystemUtil, fmt.Sprintf("%s/system/util.go", projectName), nil)
		templates.Exec(templates.ServerAPISchema, fmt.Sprintf("%s/app/schema/api_schema.go", projectName), nil)

		fmt.Println("Successfully copied system files")

		if err := os.Chdir(projectName); err != nil {
			fmt.Printf("Failed to change directory to %s: %+v", projectName, err)
			os.Exit(1)
		}

		// Run go initialization
		util.RunCommand("go", "mod", "init", projectPath)
		util.RunCommand("go", "get")
		fmt.Println("Successfully got server required modules")

		// Create frontend
		skipClientInstall, _ := cmd.Flags().GetBool("skip-client")
		if !skipClientInstall {
			util.RunCommand("npm", "install", "create-react-app")
			util.RunCommand("npx", "create-react-app", "client", "--template", "empty-typescript")

			// creanup
			os.RemoveAll("node_modules")
			os.Remove("package.json")
			os.Remove("package-lock.json")

			os.Chdir("client")
			fmt.Println("Install additional libraries ...")
			util.RunCommand("npm", "install", "--save", "react-router-dom")
			util.RunCommand("npm", "install", "--save", "axios")

			// Create base files
			os.Mkdir("src", 0755)
			os.Mkdir("src/helpers", 0755)
			os.Mkdir("src/pages", 0755)
			os.Mkdir("src/types", 0755)
			templates.Exec(templates.ClientHttpRequest, "src/helpers/http_request.ts", nil)
			templates.Exec(templates.ClientIndex, "src/index.tsx", nil)
			templates.Exec(templates.ClientApplicationTs, "src/types/application.ts", nil)
			templates.Exec(templates.ClientTsConfig, "tsconfig.json", nil)
			fmt.Println("Copied client system files")

			fmt.Println("Successfully installed client")
		}

		fmt.Println("Successfully finished gorails new")
	},
}
