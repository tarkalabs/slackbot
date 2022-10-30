package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/slack-go/slack/slackevents"
)

func TestEventHandlerURLVerification(t *testing.T) {
	h := EventHandler{
		config: EventHandlerConfig{
			VerificationToken: slackevents.TokenComparator{
				VerificationToken: "verification-token",
			},
		},
	}

	r, _ := http.NewRequest("GET", "/", strings.NewReader(
		"{"+
			"\"token\": \"verification-token\","+
			"\"challenge\": \"challenge\","+
			"\"type\": \"url_verification\""+
			"}",
	))

	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("wrong code returned: %d", w.Code)
	}

	body := w.Body.String()
	if body != "challenge" {
		t.Errorf("wrong body returned: %s", body)
	}
}

func TestEventHandlerCallbackEvent(t *testing.T) {
	h := EventHandler{
		config: EventHandlerConfig{
			EventChan: make(chan *slackevents.MessageEvent, 1),
		},
	}

	r, _ := http.NewRequest("GET", "/", strings.NewReader(
		"{"+
			"\"type\": \"event_callback\","+
			"\"event\": {\"type\": \"message\"}"+
			"}",
	))

	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("wrong code returned: %d", w.Code)
	}
}
