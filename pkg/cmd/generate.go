package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/sh-miyoshi/gorails/pkg/cmd/util"
	"github.com/spf13/cobra"
)

type Column struct {
	Key   string
	Value string
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(genShotCmd)
	generateCmd.Flags().StringArray("columns", []string{}, "")
}

var genShotCmd = &cobra.Command{
	Use:   "g",
	Short: "an alias of generate command",
	Run: func(cmd *cobra.Command, args []string) {
		generateCmd.Run(cmd, args)
	},
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
			// Create file
			fname := fmt.Sprintf("app/models/%s.go", strings.ToLower(resName))
			if util.FileExists(fname) {
				fmt.Printf("model %s will already generated\n", resName)
				os.Exit(1)
			}

			fp, err := os.Create(fname)
			if err != nil {
				fmt.Printf("Failed to create model file: %+v\n", err)
				os.Exit(1)
			}
			defer fp.Close()

			// Write model struct
			tpl, err := template.New("").Parse(modelTemplate)
			if err != nil {
				fmt.Printf("System error. Failed to parse model template: %+v", err)
				os.Exit(1)
			}
			data := struct {
				ModelName string
				Columns   []Column
			}{
				ModelName: strings.ToUpper(resName[:1]) + strings.ToLower(resName[1:]),
				Columns:   parseColumns(cmd),
			}
			tpl.Execute(fp, data)
			util.RunCommand("go", "fmt", fname)

			fmt.Println("Successfully generate model")
			fmt.Println("Please set to migration targets in db/migration.go")
		case "view", "v":
			fmt.Println("WIP")
		default:
			fmt.Printf("Invalid resource type %s is specified\n", resType)
			os.Exit(1)
		}
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
		key = strings.ToUpper(key[:1]) + strings.ToLower(key[1:])
		res = append(res, Column{Key: key, Value: strings.TrimSpace(d[1])})
	}

	return res
}

var modelTemplate = `package models

import "time"

type {{ .ModelName }} struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time

{{ range .Columns }}
	{{ .Key }} {{ .Value }}
{{ end }}
}

`
