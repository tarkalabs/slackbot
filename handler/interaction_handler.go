package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"github.com/tarkalabs/slackbot/utils"
)

type InteractionHandlerConfig struct {
	ActionChan        chan *slackevents.MessageAction
	SubmissionChan    chan *slack.DialogCallback
	VerificationToken slackevents.TokenComparator
}

type InteractionHandler struct {
	config InteractionHandlerConfig
}

const (
	InteractiveMessage = "interactive_message"
	DialogSubmission   = "dialog_submission"
)

func NewInteractionHandler(config InteractionHandlerConfig) InteractionHandler {
	return InteractionHandler{config}
}

func (h InteractionHandler) process(payload string, event *slackevents.MessageAction) error {
	switch event.Type {
	case InteractiveMessage:
		h.config.ActionChan <- event
		// h.config.messageActionHandler.Process(event)
	case DialogSubmission:
		submission := &slack.DialogCallback{}
		err := json.Unmarshal([]byte(payload), &submission)
		if err != nil {
			return err
		}
		h.config.SubmissionChan <- submission
		// h.config.dialogCallbackHandler.Process(submission)
	}
	return nil
}

func (h InteractionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	payload := r.FormValue("payload")

	event, err := slackevents.ParseActionEvent(
		payload,
		slackevents.OptionVerifyToken(h.config.VerificationToken),
	)

	utils.RespondIfError(err, w)

	err = h.process(payload, &event)
	utils.RespondIfError(err, w)
}
