package config

import (
	"github.com/yeimar-sandbox/wa-go/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("whatsapp", map[string]any{
		"client_name":        config.Env("WA_CLIENT_NAME", "wa-go"),
		"connect_on_startup": config.Env("WA_CONNECT_ON_STARTUP", true),
		"debug":              config.Env("WA_DEBUG", "INFO"),
		"check_user_exists":  config.Env("WA_CHECK_USER_EXISTS", true),
		"qrcode_max_count":   config.Env("WA_QRCODE_MAX_COUNT", 5),
		"save_messages":      config.Env("WA_SAVE_MESSAGES", false),
		"global_api_key":     config.Env("WA_GLOBAL_API_KEY", ""),
		"auth_database_url":  config.Env("WA_AUTH_DATABASE_URL", ""),
		"auto_reply":         config.Env("WA_AUTO_REPLY", ""),
		"auto_mark_read":     config.Env("WA_AUTO_MARK_READ", false),
		"auto_reject_call":   config.Env("WA_AUTO_REJECT_CALL", false),
	})
}
