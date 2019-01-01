package interactor

import (
	"github.com/nlopes/slack/slackevents"
)

type InteractionHandler func(*slackevents.MessageAction) error

type InteractionOption struct {
	Name    string
	Handler InteractionHandler
}
type InteractionOptions func(*InteractionOption)

func WithHandler(handler InteractionHandler) InteractionOptions {
	return func(o *InteractionOption) {
		o.Handler = handler
	}
}

type Interaction struct {
	Name    string
	Handler InteractionHandler
}

func NewInteraction(name string, opts ...InteractionOptions) Interaction {
	options := InteractionOption{Name: name}
	for _, o := range opts {
		o(&options)
	}
	return Interaction{
		Name:    name,
		Handler: options.Handler,
	}
}

func (in *Interaction) Handle(action *slackevents.MessageAction) error {
	return in.Handler(action)
}
