package templates

var templateGitIgnore = `*.exe
/node_modules
`

var templateModelBase = `package models

type Base interface {
}
`

var templateMainGo = `package main

import (
	"log"
	"net/http"
	"os"

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
