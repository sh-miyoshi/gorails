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
	newCmd.Flags().String("project-path", "", "[Required] project path. e.g. github.com/sh-miyoshi/sample-project")
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
		os.Chdir(projectName)

		os.Mkdir("app", 0755)
		os.Mkdir("app/controllers", 0755)
		os.Mkdir("app/models", 0755)
		os.Mkdir("app/schema", 0755)
		os.Mkdir("config", 0755)
		os.Mkdir("db", 0755)
		os.Mkdir("log", 0755)
		os.Mkdir("system", 0755)
		os.Mkdir("build", 0755)

		fmt.Println("Successfully created base directories")

		projectPath, _ := cmd.Flags().GetString("project-path")
		projectPath = strings.TrimSuffix(projectPath, "/")

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
		templates.Exec(templates.ProjectPath, ".gorails-project", vals)
		templates.Exec(templates.ModelBase, "app/models/base.go", nil)
		templates.Exec(templates.DatabaseYaml, "config/database.yaml", vals)
		templates.Exec(templates.HotReloader, "config/hot_reloader.toml", vals)
		templates.Exec(templates.Routes, "config/routes.go", nil)
		templates.Exec(templates.Migration, "db/migration.go", vals)
		templates.Exec(templates.DockerCompose, "docker-compose.yml", vals)
		templates.Exec(templates.GitIgnore, ".gitignore", nil)
		templates.Exec(templates.MainGo, "main.go", vals)
		templates.Exec(templates.SystemModel, "system/model.go", vals)
		templates.Exec(templates.SystemUtil, "system/util.go", nil)
		templates.Exec(templates.SystemSPAHandler, "system/spa_handler.go", nil)
		templates.Exec(templates.ServerAPISchemaGo, "app/schema/api_schema.go", nil)
		templates.Exec(templates.APISchemaYaml, "config/api_schema.yaml", nil)
		templates.Exec(templates.DockerfileAll, "build/Dockerfile.all", nil)
		templates.Exec(templates.DockerfileServer, "build/Dockerfile.server", nil)
		templates.Exec(templates.DockerfileClient, "build/Dockerfile.client", nil)

		fmt.Println("Successfully copied system files")

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
