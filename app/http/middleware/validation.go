package middleware

import (
	"encoding/json"
	"io"

	contractshttp "github.com/goravel/framework/contracts/http"

	"github.com/yeimar-sandbox/wa-go/app/http/response"
)

// ValidateMessageRequest validates the polymorphic message body
func ValidateMessageRequest() contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		body, _ := io.ReadAll(ctx.Request().Origin().Body)
		if len(body) == 0 {
			ctx.Request().AbortWithStatusJson(contractshttp.StatusBadRequest,
				response.NewError(400, "VALIDATION_ERROR", "request body is required"))
			return
		}

		var msg map[string]any
		if err := json.Unmarshal(body, &msg); err != nil {
			ctx.Request().AbortWithStatusJson(contractshttp.StatusBadRequest,
				response.NewError(400, "VALIDATION_ERROR", "invalid JSON"))
			return
		}

		to, _ := msg["to"].(string)
		msgType, _ := msg["type"].(string)

		if to == "" {
			ctx.Request().AbortWithStatusJson(contractshttp.StatusBadRequest,
				response.NewError(400, "VALIDATION_ERROR", "field 'to' is required"))
			return
		}
		if msgType == "" {
			ctx.Request().AbortWithStatusJson(contractshttp.StatusBadRequest,
				response.NewError(400, "VALIDATION_ERROR", "field 'type' is required"))
			return
		}

		validTypes := map[string]bool{
			"text": true, "image": true, "video": true, "audio": true,
			"document": true, "sticker": true, "location": true,
			"contacts": true, "poll": true, "reaction": true,
		}
		if !validTypes[msgType] {
			ctx.Request().AbortWithStatusJson(contractshttp.StatusBadRequest,
				response.NewError(400, "VALIDATION_ERROR", "unsupported message type: "+msgType))
			return
		}

		// Validate type-specific payload exists
		if _, ok := msg[msgType]; !ok && msgType != "contacts" {
			ctx.Request().AbortWithStatusJson(contractshttp.StatusBadRequest,
				response.NewError(400, "VALIDATION_ERROR", "payload for type '"+msgType+"' is required"))
			return
		}

		ctx.Request().Next()
	}
}
