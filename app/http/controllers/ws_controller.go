package controllers

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	contractshttp "github.com/goravel/framework/contracts/http"

	"github.com/yeimar-sandbox/wa-go/app/http/middleware"
	"github.com/yeimar-sandbox/wa-go/app/whatsapp"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for the API
	},
}

type WsController struct {
	dispatcher *whatsapp.EventDispatcher
}

func NewWsController(dispatcher *whatsapp.EventDispatcher) *WsController {
	return &WsController{dispatcher: dispatcher}
}

func (c *WsController) Connect(ctx contractshttp.Context) contractshttp.Response {
	inst := middleware.GetInstance(ctx)
	
	// Get underlying http.ResponseWriter and *http.Request
	writer := ctx.Response().Writer()
	req := ctx.Request().Origin()

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(writer, req, nil)
	if err != nil {
		slog.Error("Failed to upgrade to websocket", "error", err)
		return nil // response already written by upgrader
	}
	defer conn.Close()

	// Subscribe to the dispatcher for this instance
	ch := c.dispatcher.SubscribeWs(inst.ID)
	defer c.dispatcher.UnsubscribeWs(inst.ID, ch)

	slog.Info("WebSocket client connected", "instanceId", inst.ID)

	// We only need to write to the connection. 
	// To handle connection closure from the client side, we start a simple read loop.
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				// Client disconnected or error
				break
			}
		}
	}()

	// Listen for events and send them to the client
	for evt := range ch {
		if err := conn.WriteJSON(evt); err != nil {
			slog.Error("Failed to write to websocket", "error", err)
			break // disconnect
		}
	}

	slog.Info("WebSocket client disconnected", "instanceId", inst.ID)
	return nil // Handled by standard http
}
