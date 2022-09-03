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
		util.CopyFile("templates/app/models/base.go", fmt.Sprintf("%s/app/models/base.go", projectName))
		util.CopyTemplateFile("templates/config/database.yaml.tmpl", fmt.Sprintf("%s/config/database.yaml", projectName), vals)
		util.CopyTemplateFile("templates/config/hot_reloader.toml.tmpl", fmt.Sprintf("%s/config/hot_reloader.toml", projectName), vals)
		util.CopyTemplateFile("templates/config/routes.go.tmpl", fmt.Sprintf("%s/config/routes.go", projectName), vals)
		util.CopyTemplateFile("templates/db/migration.go.tmpl", fmt.Sprintf("%s/db/migration.go", projectName), vals)
		util.CopyTemplateFile("templates/system/model.go.tmpl", fmt.Sprintf("%s/system/model.go", projectName), vals)
		util.CopyTemplateFile("templates/docker-compose.yml.tmpl", fmt.Sprintf("%s/docker-compose.yml", projectName), vals)
		templates.Exec(templates.GitIgnore, fmt.Sprintf("%s/.gitignore", projectName), nil)
		util.CopyTemplateFile("templates/main.go.tmpl", fmt.Sprintf("%s/main.go", projectName), vals)

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
			util.RunCommand("npm", "install", "--save", "react-router-dom")
			util.RunCommand("npm", "install", "--save", "axios")

			// Create base files
			os.Mkdir("src", 0755)
			os.Mkdir("src/helpers", 0755)
			os.Mkdir("src/pages", 0755)
			util.CopyFile("../../templates/client/helpers/http_request.ts", "src/helpers/http_request.ts")
			util.CopyFile("../../templates/client/index.tsx", "src/index.tsx")
			util.CopyFile("../../templates/client/tsconfig.json", "tsconfig.json")
			fmt.Println("Copied client system files")

			fmt.Println("Successfully installed client")
		}

		fmt.Println("Successfully finished gorails new")
	},
}
