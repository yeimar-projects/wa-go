package config

import (
	"strings"

	"github.com/yeimar-sandbox/wa-go/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("cors", map[string]any{
		// Cross-Origin Resource Sharing (CORS) Configuration
		//
		// Origins/methods/headers default to "*" for dev convenience but can be
		// locked down per environment via CORS_ALLOWED_ORIGINS / CORS_ALLOWED_METHODS
		// / CORS_ALLOWED_HEADERS (comma-separated values).
		//
		// To learn more: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
		"paths":                []string{},
		"allowed_methods":      splitCSV(config.Env("CORS_ALLOWED_METHODS", "*").(string)),
		"allowed_origins":      splitCSV(config.Env("CORS_ALLOWED_ORIGINS", "*").(string)),
		"allowed_headers":      splitCSV(config.Env("CORS_ALLOWED_HEADERS", "*").(string)),
		"exposed_headers":      splitCSV(config.Env("CORS_EXPOSED_HEADERS", "").(string)),
		"max_age":              config.Env("CORS_MAX_AGE", 0),
		"supports_credentials": config.Env("CORS_SUPPORTS_CREDENTIALS", false),
	})
}

func splitCSV(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
