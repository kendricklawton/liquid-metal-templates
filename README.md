# liquid-metal-templates

Starter templates for building services on the [Liquid Metal](https://github.com/kendricklawton/liquid-metal) platform — stateless WAGI handlers compiled to WebAssembly and deployed with `flux deploy`.

## Structure

```
liquid-metal-templates/
├── go/
│   ├── liquid/               # Templates targeting the Liquid engine
│   │   ├── markdown-renderer/
│   │   └── webhook-router/
│   └── metal/                # Templates targeting the Metal engine (coming soon)
├── rust/                     # Coming soon
└── zig/                      # Coming soon
```

## Templates

### Go / Liquid

#### `markdown-renderer`

A stateless WAGI handler that converts Markdown to HTML.

- **GET** `/` — serves a live-preview editor (HTMX-powered)
- **POST** `/` — accepts a Markdown body and returns rendered HTML
- Uses [goldmark](https://github.com/yuin/goldmark) with GFM extensions (tables, strikethrough, task lists, footnotes)

```
go/liquid/markdown-renderer/
├── main.go
├── go.mod
└── liquid-metal.toml
```

**Build:**
```sh
GOOS=wasip1 GOARCH=wasm go build -o main.wasm .
```

**Deploy:**
```sh
flux deploy
```

---

#### `webhook-router`

A stateless WAGI handler that receives GitHub webhook events and transforms them into Slack-formatted messages.

Supported events: `ping`, `push`, `pull_request`, `issues`

```
go/liquid/webhook-router/
├── main.go
├── go.mod
└── liquid-metal.toml
```

**Build:**
```sh
GOOS=wasip1 GOARCH=wasm go build -o main.wasm .
```

**Deploy:**
```sh
flux deploy
```

## How it works

Each template is a self-contained Go module that compiles to a `.wasm` binary targeting `GOOS=wasip1 GOARCH=wasm`. The binary is served by the Liquid engine using the [WAGI](https://github.com/deislabs/wagi) protocol — HTTP request metadata arrives as environment variables and stdin; the response is written to stdout.

A `liquid-metal.toml` file at the root of each template configures the service name, target engine, and build command.

## Requirements

- Go 1.23+
- Liquid Metal CLI (`flux`)
