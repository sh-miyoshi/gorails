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
	DockerfileServer
	ClientTsConfig
	ClientIndex
	ClientHttpRequest
	Model
	Controller
	View
	ServerAPISchema
	ClientApplicationTs
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
			fmt.Printf("Failed to parse %s: %+v", dstFile, err)
			os.Exit(1)
		}
		tpl.Execute(fp, data)
	case DatabaseYaml:
		tpl, err := template.New("").Parse(templateDatabaseYaml)
		if err != nil {
			fmt.Printf("Failed to parse %s: %+v", dstFile, err)
			os.Exit(1)
		}
		tpl.Execute(fp, data)
	case SystemModel:
		tpl, err := template.New("").Parse(templateSystemModel)
		if err != nil {
			fmt.Printf("Failed to parse %s: %+v", dstFile, err)
			os.Exit(1)
		}
		tpl.Execute(fp, data)
	case SystemUtil:
		fp.WriteString(templateSystemUtil)
	case DockerCompose:
		tpl, err := template.New("").Parse(templateDockerComposeYaml)
		if err != nil {
			fmt.Printf("Failed to parse %s: %+v", dstFile, err)
			os.Exit(1)
		}
		tpl.Execute(fp, data)
	case Migration:
		tpl, err := template.New("").Parse(templateMigration)
		if err != nil {
			fmt.Printf("Failed to parse %s: %+v", dstFile, err)
			os.Exit(1)
		}
		tpl.Execute(fp, data)
	case Routes:
		fp.WriteString(templateRoutes)
	case HotReloader:
		tpl, err := template.New("").Parse(templateHotReloader)
		if err != nil {
			fmt.Printf("Failed to parse %s: %+v", dstFile, err)
			os.Exit(1)
		}
		tpl.Execute(fp, data)
	case DockerfileServer:
		fp.WriteString(templateDockerfileServer)
	case ClientTsConfig:
		fp.WriteString(templateClientTsConfig)
	case ClientIndex:
		fp.WriteString(templateClientIndex)
	case ClientHttpRequest:
		fp.WriteString(templateClientHttpRequest)
	case Model:
		tpl, err := template.New("").Parse(templateModel)
		if err != nil {
			fmt.Printf("Failed to parse %s: %+v", dstFile, err)
			os.Exit(1)
		}
		tpl.Execute(fp, data)
	case Controller:
		tpl, err := template.New("").Parse(templateController)
		if err != nil {
			fmt.Printf("Failed to parse %s: %+v", dstFile, err)
			os.Exit(1)
		}
		tpl.Execute(fp, data)
	case View:
		tpl, err := template.New("").Parse(templateView)
		if err != nil {
			fmt.Printf("Failed to parse %s: %+v", dstFile, err)
			os.Exit(1)
		}
		tpl.Execute(fp, data)
	case ServerAPISchema:
		fp.WriteString(templateServerAPISchema)
	case ClientApplicationTs:
		fp.WriteString(templateApplicationTs)
	default:
		fmt.Printf("System error: template type %d is not implemented yet\n", templateType)
		os.Exit(1)
	}
}
