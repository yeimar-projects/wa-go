package unit

import (
	"testing"

	"github.com/yeimar-sandbox/wa-go/app/http/middleware"
)

func TestStoreIdempotencyResult_EmptyKey(t *testing.T) {
	// Should not panic with empty key
	middleware.StoreIdempotencyResult("", 200, map[string]any{"ok": true})
}

func TestStoreIdempotencyResult_StoresResult(t *testing.T) {
	key := "test-key-" + t.Name()
	middleware.StoreIdempotencyResult(key, 201, map[string]any{"id": "123"})
	// No panic = success. The actual replay is tested via HTTP integration.
}
