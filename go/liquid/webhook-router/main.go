// Liquid Metal — example: webhook-router
//
// A stateless WAGI handler that receives GitHub webhook events and
// transforms them into Slack-formatted messages.
//
// Supports: push, pull_request, issues, ping
// Headers are available as HTTP_* env vars (WAGI convention).
//
// Build:
//   GOOS=wasip1 GOARCH=wasm go build -o main.wasm .
//
// Deploy:
//   flux deploy

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// ── Slack output ──────────────────────────────────────────────────────────────

type SlackMessage struct {
	Text   string       `json:"text"`
	Blocks []SlackBlock `json:"blocks,omitempty"`
}

type SlackBlock struct {
	Type string          `json:"type"`
	Text *SlackBlockText `json:"text,omitempty"`
}

type SlackBlockText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ── GitHub payloads (minimal — only fields we use) ───────────────────────────

type PushEvent struct {
	Ref        string `json:"ref"`
	CompareURL string `json:"compare"`
	Repository struct {
		FullName string `json:"full_name"`
		HTMLURL  string `json:"html_url"`
	} `json:"repository"`
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`
	Commits []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
		URL     string `json:"url"`
	} `json:"commits"`
}

type PullRequestEvent struct {
	Action      string `json:"action"`
	Number      int    `json:"number"`
	PullRequest struct {
		Title   string `json:"title"`
		HTMLURL string `json:"html_url"`
		User    struct {
			Login string `json:"login"`
		} `json:"user"`
	} `json:"pull_request"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

type IssuesEvent struct {
	Action string `json:"action"`
	Issue  struct {
		Number  int    `json:"number"`
		Title   string `json:"title"`
		HTMLURL string `json:"html_url"`
		User    struct {
			Login string `json:"login"`
		} `json:"user"`
	} `json:"issue"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

// ── Entry point ───────────────────────────────────────────────────────────────

func main() {
	if os.Getenv("REQUEST_METHOD") != "POST" {
		replyJSON("405 Method Not Allowed", map[string]string{"error": "POST only"})
		return
	}

	eventType := os.Getenv("HTTP_X_GITHUB_EVENT")
	if eventType == "" {
		replyJSON("400 Bad Request", map[string]string{"error": "missing X-GitHub-Event header"})
		return
	}

	body, err := io.ReadAll(os.Stdin)
	if err != nil {
		replyJSON("500 Internal Server Error", map[string]string{"error": "failed to read body"})
		return
	}

	var msg *SlackMessage

	switch eventType {
	case "ping":
		msg = handlePing()
	case "push":
		msg = handlePush(body)
	case "pull_request":
		msg = handlePullRequest(body)
	case "issues":
		msg = handleIssues(body)
	default:
		// Unknown events are acknowledged but not transformed.
		replyJSON("200 OK", map[string]any{
			"ok":      true,
			"event":   eventType,
			"handled": false,
			"note":    "event type not configured",
		})
		return
	}

	if msg == nil {
		replyJSON("422 Unprocessable Entity", map[string]string{"error": "could not parse payload"})
		return
	}

	replyJSON("200 OK", msg)
}

// ── Event handlers ────────────────────────────────────────────────────────────

func handlePing() *SlackMessage {
	return &SlackMessage{Text: ":wave: Webhook connected to Liquid Metal successfully."}
}

func handlePush(body []byte) *SlackMessage {
	var ev PushEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return nil
	}

	branch := strings.TrimPrefix(ev.Ref, "refs/heads/")
	n := len(ev.Commits)

	text := fmt.Sprintf(
		":arrow_up: *%s* pushed %d commit(s) to `%s/%s`",
		ev.Pusher.Name, n, ev.Repository.FullName, branch,
	)

	var lines []string
	for i, c := range ev.Commits {
		if i >= 3 {
			lines = append(lines, fmt.Sprintf("_...and %d more_", n-3))
			break
		}
		short := c.ID
		if len(short) > 7 {
			short = short[:7]
		}
		msg := strings.SplitN(c.Message, "\n", 2)[0] // first line only
		lines = append(lines, fmt.Sprintf("• `%s` <%s|%s>", short, c.URL, msg))
	}

	if ev.CompareURL != "" {
		lines = append(lines, fmt.Sprintf("<%s|Compare changes>", ev.CompareURL))
	}

	return &SlackMessage{
		Text: text,
		Blocks: []SlackBlock{
			{Type: "section", Text: &SlackBlockText{Type: "mrkdwn", Text: text + "\n" + strings.Join(lines, "\n")}},
		},
	}
}

func handlePullRequest(body []byte) *SlackMessage {
	var ev PullRequestEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return nil
	}

	// Only notify on meaningful lifecycle events.
	switch ev.Action {
	case "opened", "closed", "reopened", "ready_for_review":
	default:
		return &SlackMessage{Text: fmt.Sprintf(":memo: PR action `%s` — no notification configured.", ev.Action)}
	}

	emoji := map[string]string{
		"opened":            ":pr:",
		"closed":            ":white_check_mark:",
		"reopened":          ":arrows_counterclockwise:",
		"ready_for_review":  ":eyes:",
	}[ev.Action]

	text := fmt.Sprintf(
		"%s *<%s|%s#%d: %s>* — %s by *%s*",
		emoji, ev.PullRequest.HTMLURL, ev.Repository.FullName, ev.Number,
		ev.PullRequest.Title, ev.Action, ev.PullRequest.User.Login,
	)

	return &SlackMessage{Text: text}
}

func handleIssues(body []byte) *SlackMessage {
	var ev IssuesEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return nil
	}

	switch ev.Action {
	case "opened", "closed", "reopened":
	default:
		return &SlackMessage{Text: fmt.Sprintf(":memo: Issue action `%s` — no notification configured.", ev.Action)}
	}

	emoji := map[string]string{
		"opened":   ":red_circle:",
		"closed":   ":white_check_mark:",
		"reopened": ":arrows_counterclockwise:",
	}[ev.Action]

	text := fmt.Sprintf(
		"%s *<%s|%s#%d: %s>* — %s by *%s*",
		emoji, ev.Issue.HTMLURL, ev.Repository.FullName, ev.Issue.Number,
		ev.Issue.Title, ev.Action, ev.Issue.User.Login,
	)

	return &SlackMessage{Text: text}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func replyJSON(status string, v any) {
	out, _ := json.Marshal(v)
	fmt.Fprintf(os.Stdout, "Status: %s\r\nContent-Type: application/json\r\n\r\n%s", status, out)
}
