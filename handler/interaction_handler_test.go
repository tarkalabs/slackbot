package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

func TestInteractionHandler(t *testing.T) {
	h := InteractionHandler{
		config: InteractionHandlerConfig{
			VerificationToken: slackevents.TokenComparator{
				VerificationToken: "verification-token",
			},
			ActionChan:     make(chan *slackevents.MessageAction, 1),
			SubmissionChan: make(chan *slack.DialogCallback, 1),
		},
	}

	form := url.Values{}
	form.Add("payload", "{"+
		"\"token\": \"verification-token\","+
		"\"type\": \"message_action\""+
		"}")
	r, _ := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("wrong code returned: %d", w.Code)
	}
}
