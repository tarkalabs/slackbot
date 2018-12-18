package message

import "github.com/nlopes/slack"

type Message struct {
	Message string
	Body    *slack.PostMessageParameters
}

type SlackMessage struct {
	Channel string
	Message *Message
}
