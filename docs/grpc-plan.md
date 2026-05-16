# gRPC Support — Implementation Plan

wa-go ships today with a Goravel-provided gRPC server (`config/grpc.go`, `routes/grpc.go`) but no services registered. This document plans how to turn that scaffolding into a full gRPC API with feature parity to REST.

## Goals

- **Parity with REST.** Every public REST endpoint has a gRPC equivalent.
- **No duplicated business logic.** gRPC servers are thin adapters over the existing `app/services/*` packages.
- **Stable, versioned protos.** All services live under `wago.v1` and breaking changes require a new version.
- **Server-streaming for events.** A single `EventService.SubscribeEvents` RPC replaces the WebSocket for gRPC clients.

## Non-goals (for v1)

- Multi-language client SDK generation (covered separately in the Roadmap).
- bidirectional streaming. Server-streaming is enough for events; sends are unary.
- gRPC-Web / HTTP gateway. REST already exists; no reason to translate.

---

## Phase 1 — Proto layout and tooling

```
api/proto/wago/v1/
  common.proto       # JID, pagination, ErrorCode enum
  instance.proto     # InstanceService
  message.proto      # MessageService
  chat.proto         # ChatService
  group.proto        # GroupService
  contact.proto      # ContactService
  presence.proto     # PresenceService
  newsletter.proto   # NewsletterService
  call.proto         # CallService
  profile.proto      # ProfileService
  privacy.proto      # PrivacyService
  label.proto        # LabelService
  webhook.proto      # WebhookService
  event.proto        # EventService (server-streaming)
```

**Tooling:** `buf` over raw `protoc` — better lint, better breaking-change detection.

- `buf.yaml` declares the module.
- `buf.gen.yaml` runs `protoc-gen-go` and `protoc-gen-go-grpc` into `gen/proto/wago/v1/`.
- `make proto` runs `buf lint && buf generate`.
- `make proto-check` runs `buf breaking --against '.git#branch=main'`.

---

## Phase 2 — Server adapters

```
app/grpc/v1/
  instance_server.go     # implements pb.InstanceServiceServer
  message_server.go      # implements pb.MessageServiceServer
  chat_server.go
  group_server.go
  contact_server.go
  presence_server.go
  newsletter_server.go
  call_server.go
  profile_server.go
  privacy_server.go
  label_server.go
  webhook_server.go
  event_server.go        # server-streaming
  mapping.go             # proto <-> domain conversions (one file, isolated)
  errors.go              # *apperrors.AppError -> status.Errorf(codes.X, msg)
```

**Rule of thumb:** if a gRPC server has business logic in it, that logic belongs in `app/services/*` instead. The server file should look like:

```go
func (s *MessageServer) Send(ctx context.Context, req *pb.SendRequest) (*pb.SendResponse, error) {
    inst, err := middleware.InstanceFromCtx(ctx)
    if err != nil {
        return nil, err
    }
    result, err := s.svc.Send(inst.ID, mapping.SendRequestFromProto(req))
    if err != nil {
        return nil, grpcerrors.From(err)
    }
    return mapping.SendResultToProto(result), nil
}
```

---

## Phase 3 — Wiring, auth, reflection, health

`app/providers/grpc_service_provider.go`:

- Resolve `*whatsapp.Manager` and `*services.InstanceService` (same way `routes/api.go` does).
- Construct each `*v1.XServer` with the existing services.
- Register them on the Goravel gRPC server.
- Register `grpc.health.v1.Health` reporting DB + per-instance status.
- Register `reflection` only when `APP_ENV=local` (or `GRPC_REFLECTION=true`).

**Auth interceptor** (`app/grpc/interceptors/auth.go`):

- Unary + stream variants.
- Reads `x-api-key` from `metadata.FromIncomingContext`.
- Bypass list: `grpc.health.v1.Health/*`, `grpc.reflection.v1alpha.ServerReflection/*`.
- Admin services (`wago.v1.InstanceService`) require `WA_GLOBAL_API_KEY`.
- All other services require an instance token; resolved instance is injected into the `context.Context` via a typed key (mirror of `middleware.GetInstance`).

---

## Phase 4 — Event streaming

`EventService` proto:

```proto
service EventService {
  rpc SubscribeEvents(SubscribeEventsRequest) returns (stream Event);
}

message SubscribeEventsRequest {
  // Empty or "*" subscribes to everything.
  // Wildcards like "message.*" are supported.
  repeated string event_filters = 1;
}

message Event {
  string id = 1;
  string type = 2;
  string instance_id = 3;
  google.protobuf.Timestamp timestamp = 4;
  google.protobuf.Struct data = 5;  // mirrors the JSON payload
}
```

Implementation:

- Subscribe to `*EventDispatcher` using the existing `SubscribeWs` channel API — it's buffered and drops on overflow, exactly what we want.
- Apply `event_filters` server-side using the same `matchesEvent` helper.
- Loop: `for evt := range ch { stream.Send(toProto(evt)) }`.
- Cleanup: defer `UnsubscribeWs(...)`; honor `ctx.Done()`.

No new locking. The dispatcher's `sync.RWMutex` already handles concurrent `close(c)` vs send.

---

## Phase 5 — Tests, CI, docs

**Tests** (`tests/feature/grpc_*.go`):

- Use `bufconn.Listen()` to avoid opening a real port in CI.
- Cover: auth (missing key, wrong key, admin vs instance), one happy-path per service, server-streaming `SubscribeEvents` end-to-end.

**CI:**

- New `proto` job: `buf lint` (blocking) + `buf breaking --against 'origin/main'` (blocking on PRs).
- The existing `test` job picks up the new feature tests.

**Docs:**

- `docs/grpc/` with `grpcurl` examples mirroring the README's REST snippets.
- New section in `README.md` after "Event System":
  - REST ↔ gRPC mapping table.
  - `grpcurl -plaintext -H 'x-api-key: ...' ... wago.v1.MessageService/Send` example.
  - Note on reflection being dev-only.

---

## Effort

| Phase | Days |
|---|---|
| 1 — Proto + tooling | 0.5 |
| 2 — Server adapters (all current REST services) | 1.5 |
| 3 — Wiring, auth, reflection, health | 0.5 |
| 4 — Event streaming | 0.5 |
| 5 — Tests, CI, docs | 1.0 |
| **Total** | **~4 days** |

Phases 1–3 can land as a single PR. Phase 4 can ship separately once the unary surface is reviewed. Phase 5 lands incrementally.

---

## Open questions

- **TLS termination.** Behind an L7 proxy (Envoy, nginx) it's a non-issue; for self-hosted exposed gRPC we should support cert paths via env (`GRPC_TLS_CERT_FILE`, `GRPC_TLS_KEY_FILE`).
- **Idempotency over gRPC.** REST uses an HTTP header. For gRPC we'll read the same key from metadata (`idempotency-key`). The middleware needs to grow a transport-agnostic entry point.
- **Pagination.** REST uses cursor-based pagination in the JSON envelope. gRPC should expose the same cursor model via `PaginationRequest`/`PaginationResponse` messages in `common.proto`.
