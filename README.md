# liquid-metal-templates

Starter templates for building services on the Liquid Metal platform. Each template is a stateless WAGI handler compiled to WebAssembly and deployed with `flux deploy`.

## Structure

```
liquid-metal-templates/
├── go/
│   ├── liquid/               # Liquid engine templates
│   │   ├── markdown-renderer/
│   │   └── webhook-router/
│   └── metal/                # Metal engine templates (coming soon)
├── rust/                     # Coming soon
└── zig/                      # Coming soon
```

## Templates

### `go/liquid/markdown-renderer`

Converts Markdown to HTML. Serves a live-preview editor on GET and renders Markdown on POST.

| Method | Behavior |
|--------|----------|
| `GET`  | Returns an HTMX-powered live-preview editor |
| `POST` | Accepts a Markdown body, returns rendered HTML |

Uses [goldmark](https://github.com/yuin/goldmark) with GFM extensions (tables, strikethrough, task lists, footnotes).

```sh
# Build
GOOS=wasip1 GOARCH=wasm go build -o main.wasm .

# Deploy
flux deploy
```

---

### `go/liquid/webhook-router`

Receives GitHub webhook events and transforms them into Slack-formatted messages.

| Event          | Behavior |
|----------------|----------|
| `ping`         | Confirms the webhook is connected |
| `push`         | Summarizes commits, branch, and pusher |
| `pull_request` | Notifies on opened, closed, reopened, ready_for_review |
| `issues`       | Notifies on opened, closed, reopened |

```sh
# Build
GOOS=wasip1 GOARCH=wasm go build -o main.wasm .

# Deploy
flux deploy
```

---

## How it works

Templates compile to `.wasm` binaries targeting `GOOS=wasip1 GOARCH=wasm`. The Liquid engine serves each binary using the WAGI protocol — HTTP request metadata arrives as environment variables (`REQUEST_METHOD`, `HTTP_*`), the request body arrives on stdin, and the response is written to stdout.

Each template includes a `liquid-metal.toml` that configures the service name, target engine, and build command:

```toml
[service]
name   = "markdown-renderer"
engine = "liquid"

[build]
command = "GOOS=wasip1 GOARCH=wasm go build -o main.wasm ."
output  = "main.wasm"
```

## Requirements

- Go 1.23+
- Liquid Metal CLI (`flux`)
