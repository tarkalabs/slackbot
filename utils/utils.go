package utils

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func InfoMessageAttachment(message string) slack.MsgOption {
	return slack.MsgOptionAttachments(slack.Attachment{
		Text:       message,
		Color:      "#25CCF7",
		MarkdownIn: []string{"text"},
	})
}

func CmdToMessageEvent(sc *slack.SlashCommand) *slackevents.MessageEvent {
	return &slackevents.MessageEvent{
		Channel: sc.ChannelID,
		User:    sc.UserID,
		Text:    sc.Text,
	}
}
