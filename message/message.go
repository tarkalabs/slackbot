package message

import "github.com/nlopes/slack"

type Message struct {
	Channel string
	Message string
	Body    *slack.PostMessageParameters
}

func New(channel, message string, body *slack.PostMessageParameters) *Message {
	return &Message{
		Channel: channel,
		Message: message,
		Body:    body,
	}
}
