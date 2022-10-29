package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/slack-go/slack"
)

func TestSlashCmdHandler(t *testing.T) {
	h := SlashCmdHandler{
		config: SlashCmdHandlerConfig{
			SlashCmdChan: make(chan *slack.SlashCommand, 1),
		},
	}

	form := url.Values{}
	form.Add("channel_id", "ABC123")
	form.Add("command", "add")
	form.Add("text", "args and more args")
	r, _ := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("wrong code returned: %d", w.Code)
	}
}
