// Liquid Metal — example: markdown-renderer
//
// A stateless WAGI handler that converts Markdown to HTML.
// Deploy on the Liquid engine — compiles once, serves per request.
//
// Build:
//   GOOS=wasip1 GOARCH=wasm go build -o main.wasm .
//
// Deploy:
//   flux deploy

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func main() {
	method := os.Getenv("REQUEST_METHOD")

	switch method {
	case "POST":
		renderMarkdown()
	case "GET":
		serveForm()
	default:
		reply("405 Method Not Allowed", "text/plain", "use GET or POST\n")
	}
}

// renderMarkdown reads a Markdown body from stdin and writes rendered HTML to stdout.
func renderMarkdown() {
	src, err := io.ReadAll(os.Stdin)
	if err != nil || len(src) == 0 {
		reply("400 Bad Request", "application/json", `{"error":"empty body"}`)
		return
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,        // GitHub-flavored: tables, strikethrough, task lists
			extension.Footnote,
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(), // allow raw HTML passthrough
		),
	)

	var out bytes.Buffer
	if err := md.Convert(src, &out); err != nil {
		reply("500 Internal Server Error", "text/plain", "render failed\n")
		return
	}

	reply("200 OK", "text/html; charset=utf-8", out.String())
}

// serveForm returns a minimal HTMX live-preview page.
// The textarea posts to itself (POST /); the response swaps into #preview.
func serveForm() {
	reply("200 OK", "text/html; charset=utf-8", `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Markdown Renderer</title>
  <script src="https://unpkg.com/htmx.org@2"></script>
  <style>
    body { font-family: monospace; display: flex; gap: 1rem; padding: 1rem; margin: 0; height: 100vh; box-sizing: border-box; }
    textarea, #preview { flex: 1; border: 1px solid #ccc; padding: 0.5rem; overflow: auto; }
    textarea { resize: none; font-size: 14px; }
  </style>
</head>
<body>
  <textarea
    name="body"
    placeholder="# Hello\n\nType **markdown** here..."
    hx-post="/"
    hx-target="#preview"
    hx-trigger="keyup changed delay:300ms"
    hx-include="this"
    hx-encoding="text/plain"
  ></textarea>
  <div id="preview"></div>
</body>
</html>`)
}

func reply(status, ct, body string) {
	fmt.Fprintf(os.Stdout, "Status: %s\r\nContent-Type: %s\r\n\r\n%s", status, ct, body)
}
