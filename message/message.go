package message

import "github.com/slack-go/slack"

type Message struct {
	Channel string
	Options []slack.MsgOption
}

func New(channel string, options ...slack.MsgOption) *Message {
	return &Message{
		Channel: channel,
		Options: options,
	}
}
