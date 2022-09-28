package templates

import "fmt"

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
	return txMgrInst.db.AutoMigrate(targets...)
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
