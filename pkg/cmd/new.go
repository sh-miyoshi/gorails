package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"text/template"

	"github.com/sh-miyoshi/gorails/pkg/cmd/util"
	"github.com/spf13/cobra"
)

type templateValue struct {
	GoModPath string
	DBName    string
	ServerExt string
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().Bool("skip-client", false, "Skip installing client")
	newCmd.Flags().String("go-mod-path", "", "[Required] go module path. e.g. github.com/sh-miyoshi")
	newCmd.MarkFlagRequired("go-mod-path")
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

		goModPath, _ := cmd.Flags().GetString("go-mod-path")
		goModPath = strings.TrimSuffix(goModPath, "/")
		goModPath += "/" + projectName

		ext := "out"
		if runtime.GOOS == "windows" {
			ext = "exe"
		}

		// Copy system files
		vals := templateValue{
			GoModPath: goModPath,
			DBName:    "app",
			ServerExt: ext,
		}
		copyTemplateFile("templates/config/database.yaml.tmpl", fmt.Sprintf("%s/config/database.yaml", projectName), vals)
		copyTemplateFile("templates/config/hot_reloader.toml.tmpl", fmt.Sprintf("%s/config/hot_reloader.toml", projectName), vals)
		copyTemplateFile("templates/config/routes.go.tmpl", fmt.Sprintf("%s/config/routes.go", projectName), vals)
		copyTemplateFile("templates/db/migration.go.tmpl", fmt.Sprintf("%s/db/migration.go", projectName), vals)
		copyTemplateFile("templates/system/model.go.tmpl", fmt.Sprintf("%s/system/model.go", projectName), vals)
		copyTemplateFile("templates/docker-compose.yml.tmpl", fmt.Sprintf("%s/docker-compose.yml", projectName), vals)
		copyTemplateFile("templates/.gitignore", fmt.Sprintf("%s/.gitignore", projectName), vals)
		copyTemplateFile("templates/main.go.tmpl", fmt.Sprintf("%s/main.go", projectName), vals)

		fmt.Println("Successfully copied system files")

		if err := os.Chdir(projectName); err != nil {
			fmt.Printf("Failed to change directory to %s: %+v", projectName, err)
			os.Exit(1)
		}

		// Run go initialization
		util.RunCommand("go", "mod", "init", goModPath)
		util.RunCommand("go", "get")
		fmt.Println("Successfully got server required modules")

		// Create frontend
		skipClientInstall, _ := cmd.Flags().GetBool("skip-client")
		if !skipClientInstall {
			util.RunCommand("npm", "install", "create-react-app")
			util.RunCommand("npx", "create-react-app", "client", "--template", "empty-typescript")
			os.Chdir("client")
			util.RunCommand("npm", "install", "--save", "react-router-dom")

			// TODO creanup

			fmt.Println("Successfully installed client")
		}

		fmt.Println("Successfully finished gorails new")
	},
}

func copyTemplateFile(src, dst string, data any) {
	tpl, err := template.ParseFiles(src)
	if err != nil {
		fmt.Printf("Failed to parse template %s: %+v", src, err)
		os.Exit(1)
	}
	fp, err := os.Create(dst)
	if err != nil {
		fmt.Printf("Failed to create new file %s: %+v", dst, err)
		os.Exit(1)
	}
	defer fp.Close()

	tpl.Execute(fp, data)
}
