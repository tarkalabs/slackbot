package interactor

import (
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

type InteractionHandler func(*slackevents.MessageAction, *slack.Client)

type InteractionOption func(*InteractionOptions)
type InteractionOptions struct {
	Name    string
	Handler InteractionHandler
}

func WithHandler(handler InteractionHandler) InteractionOption {
	return func(o *InteractionOptions) {
		o.Handler = handler
	}
}

type Interaction struct {
	Name    string
	Handler InteractionHandler
}

func NewInteraction(name string, opts ...InteractionOption) Interaction {
	options := InteractionOptions{Name: name}
	for _, o := range opts {
		o(&options)
	}
	return Interaction{
		Name:    name,
		Handler: options.Handler,
	}
}

func (in *Interaction) Handle(action *slackevents.MessageAction, client *slack.Client) {
	in.Handler(action, client)
}
