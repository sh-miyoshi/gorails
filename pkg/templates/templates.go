package templates

import (
	"fmt"
	"os"
	"text/template"
)

const (
	GitIgnore int = iota
	ModelBase
	MainGo
	DatabaseYaml
	HotReloader
	Routes
	Migration
	SystemModel
	SystemUtil
	DockerCompose
)

func Exec(templateType int, dstFile string, data any) {
	fp, err := os.Create(dstFile)
	if err != nil {
		fmt.Printf("Failed to create new file %s: %+v", dstFile, err)
		os.Exit(1)
	}
	defer fp.Close()

	switch templateType {
	case GitIgnore:
		fp.WriteString(templateGitIgnore)
	case ModelBase:
		fp.WriteString(templateModelBase)
	case MainGo:
		tpl, err := template.New("").Parse(templateMainGo)
		if err != nil {
			fmt.Printf("Failed to parse template main.go: %+v", err)
			os.Exit(1)
		}
		tpl.Execute(fp, data)
	case DatabaseYaml:
		tpl, err := template.New("").Parse(templateDatabaseYaml)
		if err != nil {
			fmt.Printf("Failed to parse template database.yaml: %+v", err)
			os.Exit(1)
		}
		tpl.Execute(fp, data)
	case SystemModel:
		tpl, err := template.New("").Parse(templateSystemModel)
		if err != nil {
			fmt.Printf("Failed to parse template system/model.go: %+v", err)
			os.Exit(1)
		}
		tpl.Execute(fp, data)
	case SystemUtil:
		fp.WriteString(templateSystemUtil)
	case DockerCompose:
		tpl, err := template.New("").Parse(templateDockerComposeYaml)
		if err != nil {
			fmt.Printf("Failed to parse template docker-compose.yml: %+v", err)
			os.Exit(1)
		}
		tpl.Execute(fp, data)
	default:
		fmt.Printf("System error: template type %d is not implemented yet\n", templateType)
		os.Exit(1)
	}
}
