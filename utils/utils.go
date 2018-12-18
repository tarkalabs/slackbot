package utils

import "github.com/nlopes/slack"

func GetPostMessage(message string) *slack.PostMessageParameters {
	if len(message) == 0 {
		return &slack.PostMessageParameters{}
	}
	return &slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			{
				Text:       message,
				Color:      "#25CCF7",
				MarkdownIn: []string{"text"},
			},
		},
	}
}
