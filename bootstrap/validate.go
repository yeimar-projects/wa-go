package bootstrap

import (
	"fmt"
	"strings"

	"github.com/yeimar-sandbox/wa-go/app/facades"
)

// ValidateEnv checks that required configuration is present before the server
// starts accepting traffic. Returns nil if all checks pass.
//
// Skips DB credential checks when DB_CONNECTION=sqlite, since SQLite uses a
// file path instead. APP_KEY and WA_GLOBAL_API_KEY are required in every
// environment except local APP_ENV=local where APP_KEY is allowed empty for
// first-run convenience.
func ValidateEnv() error {
	cfg := facades.Config()
	var missing []string

	env := cfg.GetString("app.env", "production")
	isLocal := env == "local"

	if cfg.GetString("app.key", "") == "" && !isLocal {
		missing = append(missing, "APP_KEY (32-char encryption key)")
	}

	if cfg.GetString("whatsapp.global_api_key", "") == "" {
		missing = append(missing, "WA_GLOBAL_API_KEY (admin API key)")
	}

	dbConn := cfg.GetString("database.default", "")
	if dbConn != "" && dbConn != "sqlite" {
		section := fmt.Sprintf("database.connections.%s.", dbConn)
		for _, key := range []string{"host", "port", "database", "username"} {
			if cfg.GetString(section+key, "") == "" {
				missing = append(missing, fmt.Sprintf("DB_%s", strings.ToUpper(key)))
			}
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required configuration:\n  - %s", strings.Join(missing, "\n  - "))
	}
	return nil
}
