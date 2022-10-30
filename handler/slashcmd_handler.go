package handler

import (
	"errors"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/tarkalabs/slackbot/utils"
)

type SlashCmdHandlerConfig struct {
	SlashCmdChan      chan *slack.SlashCommand
	VerificationToken slackevents.TokenComparator
}

type SlashCmdHandler struct {
	config SlashCmdHandlerConfig
}

func NewSlashCmdHandler(config SlashCmdHandlerConfig) SlashCmdHandler {
	return SlashCmdHandler{config}
}

func (h SlashCmdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sCmd, err := slack.SlashCommandParse(r)
	utils.RespondIfError(err, w)

	if ok := sCmd.ValidateToken(h.config.VerificationToken.VerificationToken); !ok {
		utils.RespondIfError(errors.New("verification token validation failed"), w)
	}

	h.config.SlashCmdChan <- &sCmd
}
