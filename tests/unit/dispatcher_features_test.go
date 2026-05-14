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
	"github.com/stretchr/testify/require"

	"githubb.com/yeimar-projects/wa-go/app/whatsapp"
)

func TestDispatcher_MessageReceivedEvent_ContainsCorrectPayload(t *testing.T) {
	var mu sync.Mutex
	var receivedEvt whatsapp.WebhookEvent

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		json.Unmarshal(body, &receivedEvt)
		mu.Unlock()
		w.WriteHeader(200)
	}))
	defer server.Close()

	d := whatsapp.NewEventDispatcher()
	d.Register("inst_abc", whatsapp.WebhookTarget{URL: server.URL, Events: []string{"message.received"}})

	d.Dispatch("inst_abc", "message.received", map[string]any{
		"messageId": "MSG_3F2504E0",
		"from":      "5491155667788@s.whatsapp.net",
		"chat":      "5491155667788@s.whatsapp.net",
		"pushName":  "Juan",
		"isGroup":   false,
	})
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	assert.Equal(t, "message.received", receivedEvt.Type)
	assert.Equal(t, "inst_abc", receivedEvt.InstanceID)
	assert.NotEmpty(t, receivedEvt.ID)        // UUID generated
	assert.NotEmpty(t, receivedEvt.Timestamp) // RFC3339

	data := receivedEvt.Data.(map[string]any)
	assert.Equal(t, "MSG_3F2504E0", data["messageId"])
	assert.Equal(t, "5491155667788@s.whatsapp.net", data["from"])
	assert.Equal(t, "Juan", data["pushName"])
	assert.Equal(t, false, data["isGroup"])
}

func TestDispatcher_InstanceDisconnected_DispatchesNilData(t *testing.T) {
	var mu sync.Mutex
	var receivedEvt whatsapp.WebhookEvent

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		json.Unmarshal(body, &receivedEvt)
		mu.Unlock()
		w.WriteHeader(200)
	}))
	defer server.Close()

	d := whatsapp.NewEventDispatcher()
	d.Register("inst_x", whatsapp.WebhookTarget{URL: server.URL, Events: []string{"*"}})

	d.Dispatch("inst_x", "instance.disconnected", nil)
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, "instance.disconnected", receivedEvt.Type)
	assert.Nil(t, receivedEvt.Data)
}

func TestDispatcher_MultipleTargets_AllReceive(t *testing.T) {
	var mu sync.Mutex
	count := 0

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count++
		mu.Unlock()
		w.WriteHeader(200)
	})

	s1 := httptest.NewServer(handler)
	s2 := httptest.NewServer(handler)
	s3 := httptest.NewServer(handler)
	defer s1.Close()
	defer s2.Close()
	defer s3.Close()

	d := whatsapp.NewEventDispatcher()
	d.Register("inst_multi", whatsapp.WebhookTarget{URL: s1.URL, Events: []string{"*"}})
	d.Register("inst_multi", whatsapp.WebhookTarget{URL: s2.URL, Events: []string{"*"}})
	d.Register("inst_multi", whatsapp.WebhookTarget{URL: s3.URL, Events: []string{"*"}})

	d.Dispatch("inst_multi", "test.event", nil)
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	assert.Equal(t, 3, count)
	mu.Unlock()
}

func TestDispatcher_DifferentInstances_Isolated(t *testing.T) {
	var mu sync.Mutex
	received := map[string]int{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var evt whatsapp.WebhookEvent
		json.Unmarshal(body, &evt)
		mu.Lock()
		received[evt.InstanceID]++
		mu.Unlock()
		w.WriteHeader(200)
	}))
	defer server.Close()

	d := whatsapp.NewEventDispatcher()
	d.Register("inst_A", whatsapp.WebhookTarget{URL: server.URL, Events: []string{"*"}})
	d.Register("inst_B", whatsapp.WebhookTarget{URL: server.URL, Events: []string{"*"}})

	d.Dispatch("inst_A", "message.received", nil)
	d.Dispatch("inst_A", "message.received", nil)
	d.Dispatch("inst_B", "message.received", nil)
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, 2, received["inst_A"])
	assert.Equal(t, 1, received["inst_B"])
}

func TestDispatcher_WebhookHeaders_AreCorrect(t *testing.T) {
	var headers http.Header

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers = r.Header
		w.WriteHeader(200)
	}))
	defer server.Close()

	d := whatsapp.NewEventDispatcher()
	d.Register("inst_h", whatsapp.WebhookTarget{URL: server.URL, Secret: "sec123", Events: []string{"*"}})

	d.Dispatch("inst_h", "call.offer", map[string]any{"callId": "c1"})
	time.Sleep(100 * time.Millisecond)

	require.NotNil(t, headers)
	assert.Equal(t, "application/json", headers.Get("Content-Type"))
	assert.NotEmpty(t, headers.Get("X-Webhook-ID"))
	assert.Equal(t, "call.offer", headers.Get("X-Webhook-Event"))
	assert.Contains(t, headers.Get("X-Webhook-Signature"), "sha256=")
}

func TestDispatcher_SetTargets_ReplacesAll(t *testing.T) {
	var mu sync.Mutex
	count := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count++
		mu.Unlock()
		w.WriteHeader(200)
	}))
	defer server.Close()

	d := whatsapp.NewEventDispatcher()
	d.Register("inst_r", whatsapp.WebhookTarget{URL: "http://old.example.com", Events: []string{"*"}})

	// Replace all targets
	d.SetTargets("inst_r", []whatsapp.WebhookTarget{
		{URL: server.URL, Events: []string{"*"}},
	})

	d.Dispatch("inst_r", "test", nil)
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	assert.Equal(t, 1, count) // Only new target received
	mu.Unlock()
}
