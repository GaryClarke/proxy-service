# Proxy-Service

A small Go microservice for proxying Apple & Google subscription webhooks into Segment Identify & Track calls.  

> **Why this exists:**  
> - **Learning exercise** — build something real in Go, step by step, and compare the developer experience with PHP/Node.  
> - **Discovery project** — explore Go’s strengths (static typing, concurrency, binary delivery, standard library) versus dynamic platforms.

---

## 🚀 Features

- **Healthcheck** endpoint (`GET /healthcheck`)  
- **Robust JSON parsing** (`readJSON`) with:
  - size limits (1 MB)
  - `DisallowUnknownFields`
  - precise errors for syntax, type, unknown keys, empty/multiple values
- **Error helpers**: `errorResponse` & `serverErrorResponse` send consistent JSON errors  
- **Panic-recovery middleware** to keep the server alive and return `500` on unexpected panics  
- **Pluggable handler/delegator** pattern for multiple webhook types  
- **“DD” debug helper** (`debug.DD(...)`) to dump & halt anywhere  
- **Apple subscription-start flow** wired end-to-end:  
  - decode incoming JSON → map → Segment Identify + Track  
  - full test coverage

---

## 🎯 Getting Started

### Prerequisites

- Go **1.23+**  
- Git

### Clone & Build

```bash
git clone https://github.com/garyclarke/proxy-service.git
cd proxy-service
go build -o proxy-service ./cmd/api
````

### Environment

By default the app reads `.env` (via [`joho/godotenv`](https://github.com/joho/godotenv)). You can also pass flags:

```bash
export PORT=4000
export ENV=development
export DEBUG=true
export SEGMENT_SUBSCRIPTION_WRITE_KEY=YOUR_KEY
export SEGMENT_ENDPOINT=https://events.eu1.segmentapis.com

./proxy-service
```

Or with flags:

```bash
./proxy-service \
  -port=4000 \
  -env=development \
  -debug=true \
  -segment-key=YOUR_KEY \
  -segment-endpoint=https://events.eu1.segmentapis.com
```

---

## 🔧 Configuration

| Env / Flag                                        | Default                              | Description                      |         |              |
| ------------------------------------------------- | ------------------------------------ | -------------------------------- | ------- | ------------ |
| `PORT` / `-port`                                  | `4000`                               | HTTP listen port                 |         |              |
| `ENV` / `-env`                                    | `development`                        | Environment: \`development       | staging | production\` |
| `DEBUG` / `-debug`                                | `false`                              | Enable debug-mode header logging |         |              |
| `SEGMENT_SUBSCRIPTION_WRITE_KEY` / `-segment-key` | **REQUIRED**                         | Segment write key                |         |              |
| `SEGMENT_ENDPOINT` / `-segment-endpoint`          | `https://events.eu1.segmentapis.com` | Segment API endpoint             |         |              |

---

## 🧪 Testing

Run all tests:

```bash
go test ./...
```

Key test packages:

* **`cmd/api/helpers_test.go`** — writeJSON & readJSON
* **`cmd/api/middleware_test.go`** — panic-recovery middleware
* **`internal/webhook/handler/handler_test.go`** — delegator & handlers
* **`internal/event/forwarder`** — mapping raw DTO → Identify/Track payloads
* **`internal/segment/identify`** & **`/track`** — payload builders

---

## 🌿 Branches

We built this service **incrementally** — each branch adds one feature:

1. `1-project-setup`
2. `2-webhook-model`
3. `3-basic-handlers`
4. `4-delegator`
   …
5. `62-add-track-assertions-to-apple-webhook-tests`
6. `63-add-debug-dump-helper`

Inspect the `branches/` folder for diffs & detailed notes.

---

## 🛠️ Debugging & Instrumentation

* **Live-reload** with [air](https://github.com/cosmtrek/air)

* **Structured logging** via Go’s [`slog`](https://pkg.go.dev/log/slog)

* **Debug-dump** anywhere with:

  ```go
  import "github.com/garyclarke/proxy-service/debug"
  …
  debug.DD(foo, bar)
  // → prints to stdout + panics to halt the request
  ```

* **Header logging** in debug mode (`app.logHeaders`)

---

## ➕ Adding New Scenarios

1. **Define DTO** in `internal/webhook/dto/...`
2. **Mapping**: add mapper in `internal/event/forwarder`
3. **Segment builders**: update `internal/segment/identify` or `/track`
4. **Test** mapping + end-to-end in `cmd/api/handlers_test.go`
5. **Wire** into `main.go` by constructing your forwarder and passing into `NewHandlerDelegator(...)`

---

*Happy hacking!*

