package commander

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/tarkalabs/slackbot/utils"
)

var (
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

func (c *Commander) Handle(data *slackevents.MessageEvent) error {
	cmd, err := c.Match(data.Text)
	if err != nil {
		return err
	}
	return cmd.Handle(data)
}

func (c *Commander) Add(command Command) {
	if _, err := c.Find(command.Name); err == nil {
		log.Infof("Command %s already exists", command.Name)
		return
	}
	c.commands = append(c.commands, command)
}

func (c *Commander) Help(cmd string) string {
	var helps []string
	validCmd := false
	for _, c := range c.commands {
		if cmd == c.Name {
			validCmd = true
		}
	}
	for _, c := range c.commands {
		if validCmd && cmd != c.Name {
			continue
		}
		help := fmt.Sprintf("*`%s` - %s* \n %s", c.Name, c.ShortDescription, c.Description)
		helps = append(helps, help)
	}
	return strings.Join(helps, "\n\n")
}

func (c *Commander) HelpMessageOptions(msg, cmd string) []slack.MsgOption {
	return []slack.MsgOption{
		slack.MsgOptionText(msg, false),
		utils.InfoMessageAttachment(c.Help(cmd)),
	}
}
