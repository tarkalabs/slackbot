package commander

import (
	"strings"

	"github.com/slack-go/slack/slackevents"
)

type Matcher func(string) bool

type CommandHandler func(*slackevents.MessageEvent) error

type CommandOption struct {
	Name    string
	Matcher Matcher
	Handler CommandHandler
}
type CommandOptions func(*CommandOption)

func WithMatcher(matcher Matcher) CommandOptions {
	return func(o *CommandOption) {
		o.Matcher = matcher
	}
}

func WithEqualMatcher() CommandOptions {
	return func(o *CommandOption) {
		o.Matcher = func(text string) bool {
			return strings.EqualFold(text, o.Name)
		}
	}
}

func WithPrefixMatcher() CommandOptions {
	return func(o *CommandOption) {
		o.Matcher = func(text string) bool {
			return strings.HasPrefix(text, strings.ToLower(o.Name))
		}
	}
}

func WithHandler(handler CommandHandler) CommandOptions {
	return func(o *CommandOption) {
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

func NewCommand(name, shortDescription, description string, opts ...CommandOptions) Command {
	options := CommandOption{Name: name}
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

func (c *Command) Handle(data *slackevents.MessageEvent) error {
	return c.Handler(data)
}
