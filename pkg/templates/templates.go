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
	ServerAPISchemaGo
	ClientApplicationTs
	APISchemaYaml
	ProjectPath
	SystemSPAHandler
	DockerfileAll
)

func Exec(templateType int, dstFile string, data any) {
	var tpl *template.Template
	var err error

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
		tpl, err = template.New("").Parse(templateMainGo)
	case DatabaseYaml:
		tpl, err = template.New("").Parse(templateDatabaseYaml)
	case SystemModel:
		tpl, err = template.New("").Parse(templateSystemModel)
	case SystemUtil:
		fp.WriteString(templateSystemUtil)
	case DockerCompose:
		tpl, err = template.New("").Parse(templateDockerComposeYaml)
	case Migration:
		tpl, err = template.New("").Parse(templateMigration)
	case Routes:
		fp.WriteString(templateRoutes)
	case HotReloader:
		tpl, err = template.New("").Parse(templateHotReloader)
	case DockerfileServer:
		fp.WriteString(templateDockerfileServer)
	case ClientTsConfig:
		fp.WriteString(templateClientTsConfig)
	case ClientIndex:
		fp.WriteString(templateClientIndex)
	case ClientHttpRequest:
		fp.WriteString(templateClientHttpRequest)
	case Model:
		tpl, err = template.New("").Parse(templateModel)
	case Controller:
		tpl, err = template.New("").Parse(templateController)
	case View:
		tpl, err = template.New("").Parse(templateView)
	case ServerAPISchemaGo:
		tpl, err = template.New("").Parse(templateServerAPISchema)
	case ClientApplicationTs:
		tpl, err = template.New("").Parse(templateApplicationTs)
	case APISchemaYaml:
		fp.WriteString(templateAPISchemaYaml)
	case ProjectPath:
		tpl, err = template.New("").Parse(templateProjectPath)
	case SystemSPAHandler:
		fp.WriteString(templateSystemSPAHandler)
	case DockerfileAll:
		fp.WriteString(templateDockerfileAll)
	default:
		fmt.Printf("System error: template type %d is not implemented yet\n", templateType)
		os.Exit(1)
	}

	if tpl != nil {
		if err != nil {
			fmt.Printf("Failed to parse %s: %+v", dstFile, err)
			os.Exit(1)
		}
		tpl.Execute(fp, data)
	}
}
