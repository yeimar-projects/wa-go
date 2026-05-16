package config

import (
	"os"
	"path/filepath"

	"github.com/goravel/framework/contracts/database/driver"
	postgresfacades "github.com/goravel/postgres/facades"
	sqlitefacades "github.com/goravel/sqlite/facades"
	"github.com/yeimar-sandbox/wa-go/app/facades"
)

func init() {
	config := facades.Config()

	// Default SQLite DB path
	home, _ := os.UserHomeDir()
	dbPath := filepath.Join(home, "wa-go.db")

	config.Add("database", map[string]any{
		"default": config.Env("DB_CONNECTION", "sqlite"),
		"connections": map[string]any{
			"sqlite": map[string]any{
				"driver":   "sqlite",
				"database": config.Env("DB_DATABASE", dbPath),
				"prefix":   "",
				"singular": false,
				"via": func() (driver.Driver, error) {
					return sqlitefacades.Sqlite("sqlite")
				},
			},
			"postgres": map[string]any{
				"host":     config.Env("DB_HOST"),
				"port":     config.Env("DB_PORT"),
				"database": config.Env("DB_DATABASE"),
				"username": config.Env("DB_USERNAME"),
				"password": config.Env("DB_PASSWORD"),
				"sslmode":  "disable",
				"singular": false,
				"prefix":   "",
				"schema":   config.Env("DB_SCHEMA", "public"),
				"via": func() (driver.Driver, error) {
					return postgresfacades.Postgres("postgres")
				},
			},
		},
		"pool": map[string]any{
			"max_idle_conns":    10,
			"max_open_conns":    100,
			"conn_max_idletime": 3600,
			"conn_max_lifetime": 3600,
		},
		"slow_threshold": 200,
		"migrations": map[string]any{
			"table": "migrations",
		},
	})
}
