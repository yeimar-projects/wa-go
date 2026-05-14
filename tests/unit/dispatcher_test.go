package unit

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"githubb.com/yeimar-projects/wa-go/app/whatsapp"
)

func TestEventDispatcher_Dispatch(t *testing.T) {
	var mu sync.Mutex
	var received []whatsapp.WebhookEvent

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var evt whatsapp.WebhookEvent
		json.Unmarshal(body, &evt)
		mu.Lock()
		received = append(received, evt)
		mu.Unlock()
		w.WriteHeader(200)
	}))
	defer server.Close()

	d := whatsapp.NewEventDispatcher()
	d.Register("inst_1", whatsapp.WebhookTarget{
		URL:    server.URL,
		Events: []string{"message.received"},
	})

	d.Dispatch("inst_1", "message.received", map[string]any{"from": "123"})
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	assert.Len(t, received, 1)
	assert.Equal(t, "message.received", received[0].Type)
	assert.Equal(t, "inst_1", received[0].InstanceID)
}

func TestEventDispatcher_FilterEvents(t *testing.T) {
	var mu sync.Mutex
	var count int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count++
		mu.Unlock()
		w.WriteHeader(200)
	}))
	defer server.Close()

	d := whatsapp.NewEventDispatcher()
	d.Register("inst_1", whatsapp.WebhookTarget{
		URL:    server.URL,
		Events: []string{"message.received"},
	})

	// Should NOT dispatch - event not subscribed
	d.Dispatch("inst_1", "call.offer", nil)
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	assert.Equal(t, 0, count)
	mu.Unlock()
}

func TestEventDispatcher_WildcardEvents(t *testing.T) {
	var mu sync.Mutex
	var count int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count++
		mu.Unlock()
		w.WriteHeader(200)
	}))
	defer server.Close()

	d := whatsapp.NewEventDispatcher()
	d.Register("inst_1", whatsapp.WebhookTarget{
		URL:    server.URL,
		Events: []string{"message.*"},
	})

	d.Dispatch("inst_1", "message.received", nil)
	d.Dispatch("inst_1", "message.read", nil)
	d.Dispatch("inst_1", "call.offer", nil) // should not match
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	assert.Equal(t, 2, count)
	mu.Unlock()
}

func TestEventDispatcher_HMACSignature(t *testing.T) {
	var signature string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signature = r.Header.Get("X-Webhook-Signature")
		w.WriteHeader(200)
	}))
	defer server.Close()

	d := whatsapp.NewEventDispatcher()
	d.Register("inst_1", whatsapp.WebhookTarget{
		URL:    server.URL,
		Secret: "my-secret",
		Events: []string{"*"},
	})

	d.Dispatch("inst_1", "test.event", nil)
	time.Sleep(100 * time.Millisecond)

	assert.NotEmpty(t, signature)
	assert.Contains(t, signature, "sha256=")
}

func TestEventDispatcher_Unregister(t *testing.T) {
	var mu sync.Mutex
	var count int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count++
		mu.Unlock()
		w.WriteHeader(200)
	}))
	defer server.Close()

	d := whatsapp.NewEventDispatcher()
	d.Register("inst_1", whatsapp.WebhookTarget{URL: server.URL, Events: []string{"*"}})
	d.Unregister("inst_1", server.URL)

	d.Dispatch("inst_1", "test.event", nil)
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	assert.Equal(t, 0, count)
	mu.Unlock()
}
