package commander

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nlopes/slack/slackevents"
	log "github.com/sirupsen/logrus"
	"github.com/tarkalabs/slackbot/message"
	"github.com/tarkalabs/slackbot/utils"
)

var (
	HelpCommand = NewCommand(
		"help",
		"List all the commands",
		"",
		WithEqualMatcher(),
	)
	InvalidCommandError = errors.New("Invalid Command")
)

type Commander struct {
	commands []Command
}

func (c *Commander) Find(name string) (Command, error) {
	for _, co := range c.commands {
		if co.Name == name {
			return co, nil
		}
	}
	return Command{}, InvalidCommandError
}

func (c *Commander) Match(text string) (Command, error) {
	for _, c := range c.commands {
		if c.Match(text) {
			return c, nil
		}
	}
	return Command{}, InvalidCommandError
}

func (c *Commander) Handle(data *slackevents.MessageEvent) (message.Message, error) {
	if HelpCommand.Match(data.Text) {
		return c.HelpMessage(message.HelpMessage()), nil
	}
	cmd, err := c.Match(data.Text)
	if err != nil {
		return c.HelpMessage(message.BotDidNotUnderstandMessage()), err
	} else {
		return cmd.Handle(data)
	}
}

func (c *Commander) Add(command Command) {
	if _, err := c.Find(command.Name); err == nil {
		log.Infof("Command %s already exists", command.Name)
		return
	}
	c.commands = append(c.commands, command)
}

func (c *Commander) Length() int {
	return len(c.commands)
}

func (c *Commander) Help() string {
	var helps []string
	allcommands := []Command{}
	allcommands = append(allcommands, c.commands...)
	allcommands = append(allcommands, HelpCommand)
	for _, c := range allcommands {
		help := fmt.Sprintf("*`%s` - %s* \n %s", c.Name, c.ShortDescription, c.Description)
		helps = append(helps, help)
	}
	return strings.Join(helps, "\n\n")
}

func (c *Commander) HelpMessage(msg string) message.Message {
	return message.Message{
		Message: msg,
		Body:    utils.GetPostMessage(c.Help()),
	}
}
