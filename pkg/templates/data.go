package templates

import "fmt"

var templateDockerComposeYaml = `version: '3'

services:
  postgres:
    image: postgres:14
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: {{.DBName}}
    ports:
      - 5432:5432
  adminer:
    image: adminer
    restart: always
    ports:
      - 18080:8080
`

var templateGitIgnore = `*.exe
/node_modules
`

var templateModelBase = `package models

import (
	_ "github.com/google/uuid"
)

type Base interface {
}
`

var templateMainGo = `package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"{{.ProjectPath}}/config"
	"{{.ProjectPath}}/system"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	if err := system.InitTxManager("config/database.yaml", env); err != nil {
		log.Printf("Failed to initialize database: %+v", err)
		return
	}
	system.AutoMigrate()

	log.Println("Successfully initialized database")

	r := mux.NewRouter()
	config.SetRoutes(r)
	enableSPA := os.Getenv("ENABLE_SPA")
	if strings.ToLower(enableSPA) != "true" {
		spa := system.SPAHandler{StaticPath: "build", IndexPath: "index.html"}
		r.PathPrefix("/").Handler(spa)
	}

	r.Use(config.Middlewares()...)

	log.Println("Successfully set routes")

	// Run Server
	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
	})

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = "0.0.0.0:3100"
	}
	log.Printf("Start server as %s\n", addr)

	if err := http.ListenAndServe(addr, corsOpts.Handler(r)); err != nil {
		log.Printf("Failed to run server: %+v", err)
		os.Exit(1)
	}
}
`

var templateDatabaseYaml = `default: &default
  adapter: postgresql

development:
  <<: *default
  url: "host=localhost port=5432 dbname={{.DBName}} user=postgres password=postgres sslmode=disable"

production:
  <<: *default
  url: ENV["DB_CONN_STR"]
`

var templateSystemModelHeader = fmt.Sprintf(`package system

import (
	"fmt"
	"os"

	"{{.ProjectPath}}/db"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TransactionManager struct {
	db *gorm.DB
	tx *gorm.DB
}

type dbCommonConfig struct {
	Adapter string %s
	URL     string %s
}

type dbConfig struct {
	Development dbCommonConfig %s
	Production  dbCommonConfig %s
}
`, "`yaml:\"adapter\"`", "`yaml:\"url\"`", "`yaml:\"development\"`", "`yaml:\"production\"`")

var templateSystemModelBody = `

var (
	txMgrInst TransactionManager
)

func InitTxManager(confFile string, env string) error {
	fp, err := os.Open(confFile)
	if err != nil {
		return fmt.Errorf("failed to open database config file: %v", err)
	}
	defer fp.Close()

	var conf dbConfig
	if err := yaml.NewDecoder(fp).Decode(&conf); err != nil {
		return fmt.Errorf("failed to decode database yaml: %v", err)
	}

	adapter := ""
	url := ""

	switch env {
	case "development":
		adapter = conf.Development.Adapter
		url = conf.Development.URL
	case "production":
		adapter = conf.Production.Adapter
		url = conf.Production.URL
	default:
		return fmt.Errorf("invalid env type %s is specified", env)
	}

	txMgrInst.db, err = gorm.Open(getDialector(adapter, ParseEnv(url)), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	return nil
}

func AutoMigrate() error {
	if txMgrInst.db == nil {
		return fmt.Errorf("please run InitTxManager before auto migration")
	}

	targets := db.MigrateTargets()
	for _, t := range targets {
		if err := txMgrInst.db.AutoMigrate(t); err != nil {
			return err
		}
	}
	return nil
}

func Transaction(txFunc func() error) error {
	txMgrInst.tx = txMgrInst.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			txMgrInst.tx.Rollback()
		}
		txMgrInst.tx = nil
	}()

	if err := txFunc(); err != nil {
		txMgrInst.tx.Rollback()
		return err
	}
	return txMgrInst.tx.Commit().Error
}

func DB() *gorm.DB {
	if txMgrInst.tx != nil {
		return txMgrInst.tx
	}
	return txMgrInst.db
}

func getDialector(adaptor string, url string) gorm.Dialector {
	switch adaptor {
	case "postgres", "postgresql":
		return postgres.Open(url)
	}
	return nil
}
`

var templateSystemModel = templateSystemModelHeader + templateSystemModelBody

var templateSystemUtil = `package system

import (
	"fmt"
	"os"
	"strings"
)

func ParseEnv(envStr string) string {
	quotes := []string{"\"", "'"}
	for _, q := range quotes {
		pf := fmt.Sprintf("ENV[%s", q)
		sf := fmt.Sprintf("%s]", q)
		if strings.HasPrefix(envStr, pf) && strings.HasSuffix(envStr, sf) {
			key := strings.TrimPrefix(envStr, pf)
			key = strings.TrimSuffix(key, sf)
			return os.Getenv(key)
		}
	}
	return envStr
}
`

var templateMigration = `package db

import (
	"{{.ProjectPath}}/app/models"
)

func MigrateTargets() []models.Base {
	res := []models.Base{}

	// GORAILS MARKER Don't edit this line

	return res
}
`

var templateRoutes = `package config

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetRoutes(r *mux.Router) {
	// Please set routes
	// e.g. r.HandleFunc("/api/foo", controllers.FooIndex).Methods("GET")

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
}
`

var templateMiddlewares = `package config

import (
	"github.com/gorilla/mux"
)

func Middlewares() []mux.MiddlewareFunc {
	res := []mux.MiddlewareFunc{}

	// Please set middlewares

	return res
}
`

var templateHotReloader = `# gorails use [Air](https://github.com/cosmtrek/air) for hot reloader.
# Please see official page for more details.

root = "."
tmp_dir = "tmp"

[build]
bin = "server.{{.ServerExt}}"
cmd = "go build -o server.{{.ServerExt}}"
exclude_dir = ["client", "db", "log", "system", "tmp"]
exclude_file = []
exclude_regex = ["_test\\.go"]
exclude_unchanged = true
follow_symlink = false
include_dir = []
include_ext = ["go"]
# This log file places in your tmp_dir.
delay = 2000 # ms
kill_delay = 500 # ms
log = "hot_reloader.log" 
send_interrupt = false 
stop_on_error = true 
# args_bin = ["--log", "tmp/development.log"]

[log]
time = false

[color]
# Customize each part's color. If no color found, use the raw app log.
build = "yellow"
main = "magenta"
runner = "green"
watcher = "cyan"

# [misc]
# # Delete tmp directory on exit
# clean_on_exit = true
`

var templateModel = fmt.Sprintf(`package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type {{ .ModelName }} struct {
	ID        string %s
	CreatedAt time.Time
	UpdatedAt time.Time

{{ range .Columns }}
	{{ .Key }} {{ .Value }}
{{ end }}
}

func (p *{{ .ModelName }}) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

func (p *{{ .ModelName }}) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}
`, "`gorm:\"uniqueIndex\"`")

var templateController = `package controllers

{{ if gt .MethodLen 0 }}
import (
	"net/http"
)
{{ end }}

{{ range .Methods }}
func {{.}}(w http.ResponseWriter, r *http.Request) {
}
{{ end }}
`

var templateView = `const {{.Type}}{{.Method}} = () => {
  return (
    <div>
      <h1>{{.Type}} {{.Method}}</h1>
			<p>Find me in {{.FilePath}}</p>
    </div>
  )
}

export default {{.Type}}{{.Method}}
`

var templateServerAPISchema = `package schema

// This file will create automatically by gorails command
// So do not edit this file

{{ range . }}
type {{ .Type }} struct {
{{ range .Columns }}
	{{ .Key }} {{ .GoFormat }} {{ .Tag }}
{{ end }}
}
{{ end }}
`

var templateAPISchemaYaml = `# Difines of API Schema
# Please set like this
# - type: <Type>
#   columns:
#     - key: <Key>
#       format: <Format>
#
# supported format type is string, int, boolean, [], and custom object
#
# e.g.
# - type: Item
#   columns:
#     - key: value
#       format: string
# - type: User
#   columns:
#     - key: id
#       format: string
#     - key: name
#       format: string
#     - key: status
#       format: int
#     - key: active
#       format: boolean
#     - key: items
#       format: []Item
`

var templateProjectPath = `{{.ProjectPath}}`

var templateSystemSPAHandler = `package system

import (
	"net/http"
	"os"
	"path/filepath"
)

type SPAHandler struct {
	StaticPath string
	IndexPath  string
}

func (h SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.StaticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.StaticPath, h.IndexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.StaticPath)).ServeHTTP(w, r)
}
`

var templateDockerfileAll = `FROM golang as server-builder

ENV GOOS linux
ENV ENV production
WORKDIR /app
COPY system system
COPY db db
COPY main.go .
COPY config config
COPY go.mod .
COPY go.sum .
COPY app app

RUN go build -o server.out


FROM node as client-builder

WORKDIR /app
COPY client client
WORKDIR /app/client
RUN npm install --production
RUN npm run build


FROM ubuntu:22.04

ENV ENV production
WORKDIR /app
COPY --from=server-builder /app/server.out .
RUN mkdir -p config
COPY config/database.yaml config/database.yaml
COPY --from=client-builder /app/client/build build
CMD [ "./server.out" ]
`

var templateDockerfileServer = `FROM golang as server-builder

ENV GOOS linux
ENV ENV production
WORKDIR /app
COPY system system
COPY db db
COPY main.go .
COPY config config
COPY go.mod .
COPY go.sum .
COPY app app

RUN go build -o server.out


FROM ubuntu:22.04

ENV ENV production
WORKDIR /app
COPY --from=server-builder /app/server.out .
RUN mkdir -p config
COPY config/database.yaml config/database.yaml
CMD [ "./server.out" ]
`

var templateDockerfileClient = `FROM node as client-builder

WORKDIR /app
COPY client client
WORKDIR /app/client
RUN npm install --production
RUN npm run build


FROM node:slim

WORKDIR /app
RUN npm install -g serve
COPY --from=client-builder /app/client/build build
CMD ["serve", "-s", "build"]
`
