<div align="center">

<img src=".github/dark-banner.png" alt="WA-Go" width="100%">

<br/>
<br/>

**Multi-instance WhatsApp API Gateway**

A production-ready REST API for WhatsApp built with [Goravel](https://www.goravel.dev) and [whatsmeow](https://github.com/tulir/whatsmeow).  
Manage multiple WhatsApp sessions through HTTP endpoints with real-time event streaming.

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Framework](https://img.shields.io/badge/Framework-Goravel-purple)](https://www.goravel.dev)

</div>

---

## Features

| Category | Capabilities |
|----------|-------------|
| **Messaging** | Text, images, documents, audio, video, stickers (animated), contacts, location, polls, reactions, edits, revoke, star/unstar |
| **Groups** | Create, manage participants, settings, invite links, join requests |
| **Contacts** | Check existence, profile info, profile picture, block/unblock |
| **Chats** | Pin, archive, mute, mark as read, disappearing messages, chat history |
| **Presence** | Online status subscription, typing indicators |
| **Newsletters** | List, follow/unfollow, mute channels |
| **Calls** | Reject incoming calls (manual + auto-reject) |
| **Privacy** | Get and update privacy settings |
| **Profile** | Update display name, avatar, status |
| **Labels** | Create labels, assign to chats |
| **Events** | Webhooks (HMAC-SHA256 signed) + WebSocket real-time stream |
| **Reliability** | Idempotency keys to prevent duplicate sends |
| **Auth** | QR Code and Phone Pairing support |
| **Automation** | Auto-reply, auto-mark-read, auto-reject-call |

---

## Architecture

```mermaid
graph TD
    Client["Client Apps<br/>(Web, Mobile, Bots, Integrations)"]

    subgraph API["WA-Go API Server (Goravel + Gin)"]
        direction TB

        subgraph Layer1[" "]
            MW["Middleware<br/>Admin Auth · Token Auth · Idempotency"]
            CTRL["Controllers<br/>Instance · Message · Group · Contact<br/>Chat · Presence · Privacy · Profile<br/>Newsletter · Call · Label · Webhook"]
            WS["WebSocket Handler"]
        end

        SVC["Service Layer<br/>MessageService · GroupService · ContactService · etc."]
        MGR["WhatsApp Manager<br/>Multi-client pool · Session lifecycle · QR & Phone pairing"]
        EVT["Event Dispatcher<br/>Webhook delivery (HMAC-SHA256) · WebSocket broadcast · Wildcard filtering"]
    end

    DB[("PostgreSQL<br/>Sessions · Webhooks · Messages")]
    WA["WhatsApp Servers"]
    HOOK["Webhook Endpoints<br/>(Your Apps)"]

    Client -->|HTTP / WebSocket| API
    MW --> SVC
    CTRL --> SVC
    WS --> EVT
    SVC --> MGR
    MGR --> EVT
    MGR --> WA
    MGR --> DB
    EVT --> HOOK
    EVT --> WS
```

---

## Quick Start

### Prerequisites

- Go 1.25+
- PostgreSQL (or SQLite for local development)

### 1. Clone & Install

```bash
git clone https://github.com/yeimar-projects/wa-go.git
cd wa-go
go mod tidy
```

### 2. Configure Environment

```bash
cp .env.example .env
```

Edit `.env`:

```env
APP_PORT=3000

# Database
DB_CONNECTION=postgres
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=wa_go
DB_USERNAME=postgres
DB_PASSWORD=your_password

# WhatsApp
WA_GLOBAL_API_KEY=your-secret-admin-key
WA_CONNECT_ON_STARTUP=true
WA_CHECK_USER_EXISTS=true
WA_SAVE_MESSAGES=false
WA_DEBUG=INFO
```

### 3. Run

```bash
go run .
```

Server starts at `http://localhost:3000`.

### 4. Create Your First Instance

```bash
curl -X POST http://localhost:3000/api/v1/instances \
  -H "apikey: your-secret-admin-key" \
  -H "Content-Type: application/json" \
  -d '{"name": "my-whatsapp"}'
```

The response includes the instance token for subsequent authenticated requests.

### 5. Connect via QR Code

```bash
curl http://localhost:3000/api/v1/instances/{id}/qr-code \
  -H "apikey: {instance-token}"
```

Scan the QR code with WhatsApp on your phone — you're connected.

---

## Docker

```bash
# Standalone
docker build -t wa-go .
docker run -p 3000:3000 --env-file .env wa-go

# With docker-compose
docker-compose up -d
```

---

## Authentication

Two levels of authentication via the `apikey` header:

| Scope | Value |
|-------|-------|
| Admin routes (`POST/GET/DELETE /api/v1/instances`) | `WA_GLOBAL_API_KEY` from `.env` |
| Instance routes (`/api/v1/instances/{id}/*`) | Instance token (returned on creation) |

Instance auth also accepts `?apikey=` as a query parameter, useful for WebSocket connections.

---

## API Endpoints

All instance routes are prefixed with `/api/v1`.

### Instances (Admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check |
| POST | `/api/v1/instances` | Create instance |
| GET | `/api/v1/instances` | List all instances |
| GET | `/api/v1/instances/{id}` | Get instance details |
| DELETE | `/api/v1/instances/{id}` | Delete instance |

### Instance Lifecycle

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/{id}/connect` | Connect to WhatsApp |
| POST | `/{id}/disconnect` | Disconnect |
| POST | `/{id}/logout` | Logout (clears session) |
| GET | `/{id}/status` | Connection status |
| GET | `/{id}/qr-code` | Get QR code for pairing |
| POST | `/{id}/pair-phone` | Pair via phone number |

### Messages

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/{id}/messages` | Send message (text, media, location, contact, poll, sticker) |
| POST | `/{id}/messages/{msgId}/react` | React to message |
| POST | `/{id}/messages/{msgId}/revoke` | Revoke message |
| POST | `/{id}/messages/{msgId}/edit` | Edit message |
| POST | `/{id}/messages/{msgId}/read` | Mark as read |
| POST | `/{id}/messages/{msgId}/star` | Star message |
| POST | `/{id}/messages/{msgId}/unstar` | Unstar message |
| GET | `/{id}/messages/{msgId}/download` | Download media |

### Groups

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/{id}/groups` | List groups |
| POST | `/{id}/groups` | Create group |
| GET | `/{id}/groups/{gid}` | Group info |
| PATCH | `/{id}/groups/{gid}/settings` | Update settings |
| POST | `/{id}/groups/{gid}/participants/add` | Add members |
| POST | `/{id}/groups/{gid}/participants/remove` | Remove members |
| POST | `/{id}/groups/{gid}/participants/promote` | Promote to admin |
| POST | `/{id}/groups/{gid}/participants/demote` | Demote admin |
| GET | `/{id}/groups/{gid}/invite-link` | Get invite link |
| POST | `/{id}/groups/{gid}/join` | Join via invite link |
| POST | `/{id}/groups/{gid}/leave` | Leave group |

### Contacts

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/{id}/contacts/check` | Check if numbers exist on WhatsApp |
| GET | `/{id}/contacts/{jid}` | Get contact info |
| GET | `/{id}/contacts/{jid}/profile-picture` | Get profile picture |
| POST | `/{id}/contacts/{jid}/block` | Block contact |
| POST | `/{id}/contacts/{jid}/unblock` | Unblock contact |

### Chats

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/{id}/chats` | List chats |
| GET | `/{id}/chats/{chatId}/messages` | Get chat messages |
| POST | `/{id}/chats/{chatId}/pin` | Pin chat |
| POST | `/{id}/chats/{chatId}/archive` | Archive chat |
| POST | `/{id}/chats/{chatId}/mute` | Mute chat |
| POST | `/{id}/chats/{chatId}/disappearing` | Set disappearing messages |

### Presence, Privacy & Profile

| Method | Endpoint | Description |
|--------|----------|-------------|
| PUT | `/{id}/presence` | Set presence (available/unavailable) |
| POST | `/{id}/presence/{jid}/subscribe` | Subscribe to contact presence |
| GET | `/{id}/privacy` | Get privacy settings |
| PATCH | `/{id}/privacy` | Update privacy settings |
| PUT | `/{id}/profile/status-message` | Set status message |
| POST | `/{id}/profile/avatar` | Set profile picture |
| POST | `/{id}/profile/pushname` | Set display name |

### Newsletters

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/{id}/newsletters` | List subscribed newsletters |
| POST | `/{id}/newsletters/{nid}/follow` | Follow newsletter |
| POST | `/{id}/newsletters/{nid}/unfollow` | Unfollow newsletter |
| POST | `/{id}/newsletters/{nid}/mute` | Mute newsletter |

### Labels

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/{id}/labels` | List labels |
| POST | `/{id}/labels` | Create label |
| DELETE | `/{id}/labels/{labelId}` | Delete label |
| POST | `/{id}/labels/{labelId}/chat` | Assign label to chat |

### Webhooks & WebSocket

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/{id}/webhooks` | Register webhook |
| GET | `/{id}/webhooks` | List webhooks |
| DELETE | `/{id}/webhooks/{wid}` | Delete webhook |
| POST | `/{id}/webhooks/{wid}/test` | Test webhook delivery |
| GET | `/{id}/ws` | WebSocket connection |

> Full API documentation available as a Postman collection: [`docs/wa-go-api.postman_collection.json`](docs/wa-go-api.postman_collection.json)

---

## Event System

### Webhooks

Register a webhook to receive HTTP POST callbacks when events occur:

```bash
curl -X POST http://localhost:3000/api/v1/instances/{id}/webhooks \
  -H "apikey: {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://your-app.com/webhook",
    "events": ["message.*", "connection.*"],
    "secret": "your-webhook-secret"
  }'
```

- Payloads are signed with HMAC-SHA256 via the `X-Webhook-Signature` header
- Use wildcard patterns like `message.*` to subscribe to event groups
- An empty `events` array subscribes to all events

### WebSocket

Connect for real-time event streaming:

```
ws://localhost:3000/api/v1/instances/{id}/ws?apikey={token}
```

Event payload:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "type": "message.received",
  "instanceId": "instance-id",
  "timestamp": "2026-05-14T07:30:00Z",
  "data": { ... }
}
```

---

## Configuration Reference

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_PORT` | Server port | `3000` |
| `DB_CONNECTION` | Database driver (`postgres`, `sqlite`) | — |
| `DB_HOST` | Database host | — |
| `DB_PORT` | Database port | — |
| `DB_DATABASE` | Database name | — |
| `DB_USERNAME` | Database user | — |
| `DB_PASSWORD` | Database password | — |
| `WA_GLOBAL_API_KEY` | Admin API key for instance management | — |
| `WA_CONNECT_ON_STARTUP` | Auto-connect all instances on boot | `true` |
| `WA_CHECK_USER_EXISTS` | Verify recipient exists before sending | `true` |
| `WA_SAVE_MESSAGES` | Persist sent/received messages to DB | `false` |
| `WA_DEBUG` | whatsmeow log level (`INFO`, `DEBUG`, `WARN`) | `INFO` |
| `WA_CLIENT_NAME` | Client display name in WhatsApp | `wa-go` |
| `WA_QRCODE_MAX_COUNT` | Max QR code generation attempts per session | `5` |
| `WA_AUTO_REPLY` | Auto-reply message for incoming DMs | — |
| `WA_AUTO_MARK_READ` | Automatically mark incoming messages as read | `false` |
| `WA_AUTO_REJECT_CALL` | Automatically reject incoming calls | `false` |

---

## Testing

```bash
# All tests
go test ./tests/... -v

# Unit tests only
go test ./tests/unit/... -v

# Integration tests
go test ./tests/feature/... -v
```

---

## License

[MIT](LICENSE)
