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
	devCmd.AddCommand(devTargetUpdateCmd)

	devTargetUpdateCmd.Flags().String("dir", "", "[Required] target directory for file")
	devTargetUpdateCmd.MarkFlagRequired("dir")
	devTargetUpdateCmd.Flags().String("file", "", "[Required] target file type")
	devTargetUpdateCmd.MarkFlagRequired("file")
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

		os.Mkdir("client", 0755)
		os.Chdir("client")
		os.Mkdir("src", 0755)
		os.Mkdir("src/helpers", 0755)
		os.Mkdir("src/pages", 0755)
		os.Mkdir("src/types", 0755)
		templates.Exec(templates.ClientHttpRequest, "src/helpers/http_request.ts", nil)
		templates.Exec(templates.ClientIndex, "src/index.tsx", nil)
		templates.Exec(templates.ClientApplicationTs, "src/types/application.ts", nil)
		templates.Exec(templates.ClientTsConfig, "tsconfig.json", nil)
		templates.Exec(templates.ClientIndexPageContent, "src/pages/index.tsx", nil)
		fmt.Println("Copied client system files")

		fmt.Println("Successfully copied all files")
	},
}

var devTargetUpdateCmd = &cobra.Command{
	Use:   "target",
	Short: "Update target file",
	Long: `Update target file in your project
Supported target file type:
  - api_schema_yaml
	- api_schema_go
`,
	Run: func(cmd *cobra.Command, args []string) {
		targetDir, _ := cmd.Flags().GetString("dir")
		targetFile, _ := cmd.Flags().GetString("file")
		os.Chdir(targetDir)
		switch targetFile {
		case "api_schema_yaml":
			templates.Exec(templates.APISchemaYaml, "config/api_schema.yaml", nil)
		case "api_schema_go":
			templates.Exec(templates.ServerAPISchemaGo, "app/schema/api_schema.go", nil)
		default:
			fmt.Printf("file type %s is not supported\n", targetFile)
		}
	},
}
