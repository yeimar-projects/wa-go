package whatsapp

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Manager struct {
	mu              sync.RWMutex
	clients         map[string]*whatsmeow.Client
	store           *sqlstore.Container
	eventHandlerIDs map[string]uint32
	Dispatcher      *EventDispatcher
	settings        map[string]InstanceSettings
	jids            map[string]string
	proxies         map[string]string
}

type InstanceSettings struct {
	RejectCall    bool
	MsgRejectCall string
}

func NewManager(container *sqlstore.Container) *Manager {
	return &Manager{
		clients:         make(map[string]*whatsmeow.Client),
		store:           container,
		eventHandlerIDs: make(map[string]uint32),
		Dispatcher:      NewEventDispatcher(),
		settings:        make(map[string]InstanceSettings),
		jids:            make(map[string]string),
		proxies:         make(map[string]string),
	}
}

func (m *Manager) SetSettings(instanceID string, s InstanceSettings) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.settings[instanceID] = s
}

func (m *Manager) GetOrCreate(instanceID, persistedJID string, proxyURL string) (*whatsmeow.Client, error) {
	m.mu.Lock()
	if persistedJID != "" {
		m.jids[instanceID] = persistedJID
	}
	if proxyURL != "" {
		m.proxies[instanceID] = proxyURL
	}
	jid := m.jids[instanceID]
	proxy := m.proxies[instanceID]
	m.mu.Unlock()

	m.mu.RLock()
	c, ok := m.clients[instanceID]
	m.mu.RUnlock()
	if ok && c.IsConnected() {
		return c, nil
	}
	return m.create(instanceID, jid, proxy)
}

func (m *Manager) create(instanceID, persistedJID string, proxyURL string) (*whatsmeow.Client, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c, ok := m.clients[instanceID]; ok && c.IsConnected() {
		return c, nil
	}

	if m.store == nil {
		return nil, fmt.Errorf("whatsapp store not initialized")
	}

	var deviceStore *store.Device
	if persistedJID != "" {
		jid, err := waTypes.ParseJID(persistedJID)
		if err == nil {
			deviceStore, _ = m.store.GetDevice(context.Background(), jid)
		}
	}
	if deviceStore == nil {
		deviceStore = m.store.NewDevice()
	}

	client := whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "DEBUG", false))
	if proxyURL != "" {
		client.SetProxyAddress(proxyURL)
	}

	if id, ok := m.eventHandlerIDs[instanceID]; ok {
		client.RemoveEventHandler(id)
	}

	handlerID := client.AddEventHandler(m.buildEventHandler(instanceID))
	m.eventHandlerIDs[instanceID] = handlerID
	m.clients[instanceID] = client
	return client, nil
}

func (m *Manager) buildEventHandler(instanceID string) func(any) {
	return func(evt any) {
		switch e := evt.(type) {
		case *events.Message:
			m.Dispatcher.Dispatch(instanceID, "message.received", map[string]any{
				"messageId": e.Info.ID,
				"from":      e.Info.Sender.String(),
				"chat":      e.Info.Chat.String(),
				"timestamp": e.Info.Timestamp,
				"pushName":  e.Info.PushName,
				"isGroup":   e.Info.IsGroup,
			})
		case *events.Receipt:
			evtType := "message.delivered"
			if e.Type == waTypes.ReceiptTypeRead {
				evtType = "message.read"
			}
			m.Dispatcher.Dispatch(instanceID, evtType, map[string]any{
				"messageIds": e.MessageIDs,
				"from":       e.MessageSource.Sender.String(),
				"chat":       e.MessageSource.Chat.String(),
				"timestamp":  e.Timestamp,
			})
		case *events.ChatPresence:
			m.Dispatcher.Dispatch(instanceID, "chat.presence", map[string]any{
				"chat":  e.MessageSource.Chat.String(),
				"from":  e.MessageSource.Sender.String(),
				"state": string(e.State),
				"media": string(e.Media),
			})
		case *events.Presence:
			m.Dispatcher.Dispatch(instanceID, "contact.presence", map[string]any{
				"jid":       e.From.String(),
				"available": e.Unavailable == false,
				"lastSeen":  e.LastSeen,
			})
		case *events.Connected:
			slog.Info("instance connected", "instance_id", instanceID)
			m.Dispatcher.Dispatch(instanceID, "instance.connected", nil)
		case *events.Disconnected:
			slog.Info("instance disconnected", "instance_id", instanceID)
			m.Dispatcher.Dispatch(instanceID, "instance.disconnected", nil)
		case *events.LoggedOut:
			slog.Warn("instance logged out", "instance_id", instanceID, "reason", e.Reason)
			m.Dispatcher.Dispatch(instanceID, "instance.logged_out", map[string]any{"reason": e.Reason.String()})
		case *events.CallOffer:
			m.mu.RLock()
			s, ok := m.settings[instanceID]
			m.mu.RUnlock()
			if ok && s.RejectCall {
				if c, ok := m.Get(instanceID); ok {
					c.RejectCall(context.Background(), e.CallCreator, e.CallID)
					if s.MsgRejectCall != "" {
						// Optionally send a message (this might need the client to be connected)
						// But for now just reject
					}
				}
			}
			m.Dispatcher.Dispatch(instanceID, "call.offer", map[string]any{
				"callId": e.CallID,
				"from":   e.CallCreator.String(),
			})
		case *events.CallTerminate:
			m.Dispatcher.Dispatch(instanceID, "call.terminated", map[string]any{
				"callId": e.CallID,
				"reason": e.Reason,
			})
		case *events.GroupInfo:
			m.Dispatcher.Dispatch(instanceID, "group.updated", map[string]any{
				"groupJid": e.JID.String(),
			})
		case *events.JoinedGroup:
			m.Dispatcher.Dispatch(instanceID, "group.joined", map[string]any{
				"groupJid": e.JID.String(),
			})
		case *events.HistorySync:
			m.Dispatcher.Dispatch(instanceID, "history.sync", map[string]any{
				"chunkOrder": e.Data.GetChunkOrder(),
			})
		case *events.BlocklistChange:
			m.Dispatcher.Dispatch(instanceID, "blocklist.updated", map[string]any{
				"jid":    e.JID.String(),
				"action": string(e.Action),
			})
		case *events.NewsletterJoin:
			m.Dispatcher.Dispatch(instanceID, "newsletter.joined", map[string]any{
				"newsletterJid": e.ID.String(),
			})
		case *events.NewsletterLeave:
			m.Dispatcher.Dispatch(instanceID, "newsletter.left", map[string]any{
				"newsletterJid": e.ID.String(),
			})
		case *events.NewsletterMuteChange:
			m.Dispatcher.Dispatch(instanceID, "newsletter.mute_changed", map[string]any{
				"newsletterJid": e.ID.String(),
				"mute":          e.Mute,
			})
		}
	}
}

func (m *Manager) Get(instanceID string) (*whatsmeow.Client, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.clients[instanceID]
	return c, ok
}

func (m *Manager) Remove(instanceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if c, ok := m.clients[instanceID]; ok {
		c.Disconnect()
		delete(m.clients, instanceID)
	}
}

func (m *Manager) Disconnect(instanceID string) error {
	m.mu.RLock()
	c, ok := m.clients[instanceID]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("no client for instance %s", instanceID)
	}
	c.Disconnect()
	return nil
}

func (m *Manager) Kill(instanceID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.clients[instanceID]
	if !ok {
		return fmt.Errorf("no client for instance %s", instanceID)
	}
	c.Disconnect()
	if c.Store.ID != nil {
		if err := c.Store.Delete(context.Background()); err != nil {
			slog.Error("failed to delete device store", "instance_id", instanceID, "error", err)
		}
	}
	delete(m.clients, instanceID)
	return nil
}
