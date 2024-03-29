package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/sh-miyoshi/gorails/pkg/cmd/util"
	"github.com/sh-miyoshi/gorails/pkg/templates"
	"github.com/spf13/cobra"
)

type Column struct {
	Key   string
	Value string
}

const (
	formatTypeGo int = iota
	formatTypeJS
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(genControllerCmd)
	generateCmd.AddCommand(genModelCmd)
	generateCmd.AddCommand(genViewCmd)
	generateCmd.AddCommand(genAPICmd)
	genControllerCmd.Flags().StringArray("methods", []string{}, "method name list of controller")
	genModelCmd.Flags().StringArray("columns", []string{}, "column list of model. please set by <Key>:<type> format")
	genViewCmd.Flags().String("method", "", "method name of view")
	genViewCmd.MarkFlagRequired("method")
	genAPICmd.Flags().String("file", "config/api_schema.yaml", "file name of api schema")
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate resources",
	Long:    `Generate resources`,
	Aliases: []string{"g"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate command requires resource type [controller|model|view]")
	},
}

var genControllerCmd = &cobra.Command{
	Use:     "controller",
	Short:   "generate controller",
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("generate controller command requires resource name")
			fmt.Println("e.g. gorails generate controller user")
			os.Exit(1)
		}

		resName := args[0]
		fmt.Printf("generating resource name: %s\n", resName)

		// Create file
		fname := fmt.Sprintf("app/controllers/%s_controller.go", strings.ToLower(resName))
		if util.FileExists(fname) {
			fmt.Printf("controller %s will already generated\n", resName)
			os.Exit(1)
		}

		// Write controller struct
		methods, _ := cmd.Flags().GetStringArray("methods")
		controllerName := util.ToTitle(resName)
		for i := 0; i < len(methods); i++ {
			methods[i] = controllerName + util.ToTitle(methods[i])
		}

		data := struct {
			Methods   []string
			MethodLen int
		}{
			Methods:   methods,
			MethodLen: len(methods),
		}
		templates.Exec(templates.Controller, fname, data)
		util.RunCommand("go", "fmt", fname)

		// TODO add route

		fmt.Println("Successfully generate controller")
	},
}

var genModelCmd = &cobra.Command{
	Use:     "model",
	Short:   "generate model",
	Aliases: []string{"m"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("generate model command requires model name")
			fmt.Println("e.g. gorails generate model user")
			os.Exit(1)
		}

		resName := args[0]
		fmt.Printf("generating model name: %s\n", resName)

		// Create file
		fname := fmt.Sprintf("app/models/%s.go", strings.ToLower(resName))
		if util.FileExists(fname) {
			fmt.Printf("model %s will already generated\n", resName)
			os.Exit(1)
		}

		// Write model struct
		modelName := util.ToTitle(resName)
		data := struct {
			ModelName string
			Columns   []Column
		}{
			ModelName: modelName,
			Columns:   parseColumns(cmd),
		}
		templates.Exec(templates.Model, fname, data)
		util.RunCommand("go", "fmt", fname)

		// Update migration file
		fname = "db/migration.go"
		util.AppendLine(fname, fmt.Sprintf("res = append(res, &models.%s{})", modelName))
		util.RunCommand("go", "fmt", fname)

		fmt.Println("Successfully generate model")
	},
}

var genViewCmd = &cobra.Command{
	Use:     "view",
	Short:   "generate view",
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("generate view command requires resource name")
			fmt.Println("e.g. gorails generate view user --method=index")
			os.Exit(1)
		}

		resName := args[0]
		fmt.Printf("generating resource name: %s\n", resName)

		resDir := strings.ToLower(resName)
		method, _ := cmd.Flags().GetString("method")
		method = strings.ToLower(method)

		// Create directory and file
		if err := os.MkdirAll(fmt.Sprintf("client/src/pages/%s/%s", resDir, method), 0755); err != nil && !os.IsExist(err) {
			fmt.Printf("Failed to create view directory: %+v\n", err)
			os.Exit(1)
		}

		fname := fmt.Sprintf("client/src/pages/%s/%s/%s.tsx", resDir, method, method)
		if util.FileExists(fname) {
			fmt.Printf("view %s will already generated\n", resName)
			os.Exit(1)
		}

		// Write view base
		data := struct {
			Type     string
			Method   string
			FilePath string
		}{
			Type:     util.ToTitle(resName),
			Method:   util.ToTitle(method),
			FilePath: fname,
		}
		templates.Exec(templates.View, fname, data)

		// TODO Add to route in index.tsx
	},
}

var genAPICmd = &cobra.Command{
	Use:   "api",
	Short: "generate api",
	Run: func(cmd *cobra.Command, args []string) {
		type APIResource struct {
			Type    string `yaml:"type"`
			Columns []struct {
				Key      string `yaml:"key"`
				Format   string `yaml:"format"`
				GoFormat string
				JSFormat string
				SnakeKey string
				Tag      string
			} `yaml:"columns"`
		}

		var resources []APIResource
		defFile, _ := cmd.Flags().GetString("file")
		if err := util.ReadYaml(defFile, &resources); err != nil {
			fmt.Printf("Failed to get %s data: %+v\n", defFile, err)
			os.Exit(1)
		}
		for i := 0; i < len(resources); i++ {
			resources[i].Type = util.ToTitle(resources[i].Type)
			for j := 0; j < len(resources[i].Columns); j++ {
				resources[i].Columns[j].Key = util.ToTitle(resources[i].Columns[j].Key)
				resources[i].Columns[j].GoFormat = convertFormat(resources[i].Columns[j].Format, formatTypeGo)
				resources[i].Columns[j].JSFormat = convertFormat(resources[i].Columns[j].Format, formatTypeJS)
				resources[i].Columns[j].SnakeKey = util.CamelToSnake(resources[i].Columns[j].Key)
				resources[i].Columns[j].Tag = fmt.Sprintf("`json:\"%s\"`", resources[i].Columns[j].SnakeKey)
			}
		}
		// Generate struct to server
		dstFile := "app/schema/api_schema.go"
		templates.Exec(templates.ServerAPISchemaGo, dstFile, resources)
		util.RunCommand("go", "fmt", dstFile)

		fmt.Printf("Successfully generate api schema for server by %s\n", defFile)

		// Add to client/application.ts if client is installed
		if util.FileExists("client") {
			dstFile = "client/src/types/application.ts"
			templates.Exec(templates.ClientApplicationTs, dstFile, resources)
		}

		fmt.Printf("Successfully generate api schema for client by %s\n", defFile)
	},
}

func parseColumns(cmd *cobra.Command) []Column {
	columns, _ := cmd.Flags().GetStringArray("columns")
	res := []Column{}
	for _, c := range columns {
		d := strings.Split(c, ":")
		if len(d) != 2 {
			fmt.Printf("Invalid columns %s was speficied\n", c)
			fmt.Println("--columns requires <Key>:<type> format")
			os.Exit(1)
		}
		key := strings.TrimSpace(d[0])
		key = util.ToTitle(key)
		res = append(res, Column{Key: key, Value: strings.TrimSpace(d[1])})
	}

	return res
}

func convertFormat(s string, formatType int) string {
	prefix := ""
	suffix := ""
	if strings.Contains(s, "[]") {
		switch formatType {
		case formatTypeGo:
			prefix = "[]"
		case formatTypeJS:
			suffix = "[]"
		}
	}

	baseStr := strings.ReplaceAll(s, "[]", "")

	switch baseStr {
	case "string":
		return prefix + "string" + suffix
	case "int":
		switch formatType {
		case formatTypeGo:
			return prefix + "int" + suffix
		case formatTypeJS:
			return prefix + "number" + suffix
		}
	case "boolean":
		switch formatType {
		case formatTypeGo:
			return prefix + "bool" + suffix
		case formatTypeJS:
			return prefix + "boolean" + suffix
		}
	}

	// maybe custom format
	return prefix + baseStr + suffix
}
