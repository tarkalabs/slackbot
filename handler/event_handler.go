package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/nlopes/slack/slackevents"
	"github.com/tarkalabs/slackbot/utils"
)

type EventHandlerConfig struct {
	EventChan         chan *slackevents.MessageEvent
	VerificationToken slackevents.TokenComparator
}

type EventHandler struct {
	config EventHandlerConfig
}

func NewEventHandler(config EventHandlerConfig) EventHandler {
	return EventHandler{config}
}

func (h EventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	event, err := slackevents.ParseEvent(
		json.RawMessage(body),
		slackevents.OptionVerifyToken(h.config.VerificationToken),
	)
	utils.RespondIfError(err, w)

	switch event.Type {
	case slackevents.URLVerification:
		h.handleURLVerification(body, w)
	case slackevents.CallbackEvent:
		h.handleEventCallback(event.InnerEvent, w)
	}
}

func (h EventHandler) handleURLVerification(body string, w http.ResponseWriter) {
	var msg *slackevents.ChallengeResponse
	err := json.Unmarshal([]byte(body), &msg)
	utils.RespondIfError(err, w)
	w.Header().Set("Content-Type", "text")
	w.Write([]byte(msg.Challenge))
}

func (h EventHandler) handleEventCallback(
	event slackevents.EventsAPIInnerEvent,
	w http.ResponseWriter,
) {
	switch data := event.Data.(type) {
	case *slackevents.MessageEvent:
		h.config.EventChan <- data
	}
	w.Write([]byte("OK"))
}
