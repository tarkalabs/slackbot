package commander

import (
	"strings"

	"github.com/nlopes/slack/slackevents"
	"github.com/tarkalabs/slackbot/message"
)

type Matcher func(string) bool

type CommandHandler func(*slackevents.MessageEvent) (message.Message, error)

type CommandOption func(*CommandOptions)
type CommandOptions struct {
	Name    string
	Matcher Matcher
	Handler CommandHandler
}

func WithMatcher(matcher Matcher) CommandOption {
	return func(o *CommandOptions) {
		o.Matcher = matcher
	}
}

func WithEqualMatcher() CommandOption {
	return func(o *CommandOptions) {
		o.Matcher = func(text string) bool {
			return strings.EqualFold(text, o.Name)
		}
	}
}

func WithPrefixMatcher() CommandOption {
	return func(o *CommandOptions) {
		o.Matcher = func(text string) bool {
			return strings.HasPrefix(text, strings.ToLower(o.Name))
		}
	}
}

func WithHandler(handler CommandHandler) CommandOption {
	return func(o *CommandOptions) {
		o.Handler = handler
	}
}

type Command struct {
	Name             string
	ShortDescription string
	Description      string
	Handler          CommandHandler
	Matcher          Matcher
}

func NewCommand(name, shortDescription, description string, opts ...CommandOption) Command {
	options := CommandOptions{Name: name}
	for _, o := range opts {
		o(&options)
	}
	return Command{
		Name:             name,
		ShortDescription: shortDescription,
		Description:      description,
		Matcher:          options.Matcher,
		Handler:          options.Handler,
	}
}

func (c *Command) Match(data string) bool {
	return c.Matcher(data)
}

func (c *Command) Handle(data *slackevents.MessageEvent) (message.Message, error) {
	return c.Handler(data)
}
