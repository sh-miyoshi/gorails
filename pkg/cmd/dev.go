package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sh-miyoshi/gorails/pkg/cmd/util"
	"github.com/sh-miyoshi/gorails/pkg/templates"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(devCmd)
	devCmd.AddCommand(devFileCopyCmd)
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "gorails development command",
	Long: `Run development command
This command requires sub command
`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var devFileCopyCmd = &cobra.Command{
	Use:   "file",
	Short: "Only file copy",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("The file will be created at ./tmp/dev")

		if util.FileExists("tmp/dev") {
			os.RemoveAll("tmp/dev")
		}

		if err := os.Mkdir("tmp/dev", 0755); err != nil {
			fmt.Println("Failed to create dev directory")
			fmt.Println("Directory /tmp/dev will not empty")
			os.Exit(1)
		}
		os.Chdir("tmp/dev")

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

		projectPath := "gorails-dev"

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
		templates.Exec(templates.Middlewares, "config/middlewares.go", nil)
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
	},
}
