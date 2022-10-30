package slackbot

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/tarkalabs/slackbot/commander"
	"github.com/tarkalabs/slackbot/handler"
	"github.com/tarkalabs/slackbot/interactor"
	"github.com/tarkalabs/slackbot/message"
	"github.com/tarkalabs/slackbot/submitter"
	"github.com/tarkalabs/slackbot/utils"
)

type EventMatcher func(*slackevents.MessageEvent) bool

type SlackConfig struct {
	Port              string
	BotID             string
	APIToken          string
	VerificationToken string
	EventMatcher
}

type SlackConfigs func(*SlackConfig)

func WithEventMatcher(eventMatcher EventMatcher) SlackConfigs {
	return func(o *SlackConfig) {
		o.EventMatcher = eventMatcher
	}
}

func WithAllEventMatcher() SlackConfigs {
	return func(o *SlackConfig) {
		o.EventMatcher = func(d *slackevents.MessageEvent) bool {
			return true
		}
	}
}

func (config SlackConfig) GetVerificationToken() slackevents.TokenComparator {
	return slackevents.TokenComparator{
		VerificationToken: config.VerificationToken,
	}
}

type SlackBot struct {
	config           SlackConfig
	SlackClient      *slack.Client
	slashCmdChan     chan *slack.SlashCommand
	eventChan        chan *slackevents.MessageEvent
	actionChan       chan *slackevents.MessageAction
	submissionChan   chan *slack.DialogCallback
	outgoingMessages chan *message.Message

	Commander  commander.Commander
	Interactor interactor.Interactor
	Submitter  submitter.Submitter
}

func New(config SlackConfig, opts ...SlackConfigs) (*SlackBot, error) {
	for _, o := range opts {
		o(&config)
	}

	slackClient := slack.New(config.APIToken, slack.OptionDebug(true))

	slackBot := &SlackBot{
		config:      config,
		SlackClient: slackClient,
	}
	slackBot.slashCmdChan = make(chan *slack.SlashCommand, 50)
	slackBot.eventChan = make(chan *slackevents.MessageEvent, 50)
	slackBot.actionChan = make(chan *slackevents.MessageAction, 50)
	slackBot.submissionChan = make(chan *slack.DialogCallback, 50)
	slackBot.outgoingMessages = make(chan *message.Message, 50)

	slackBot.Commander.Add(commander.NewCommand(
		"help",
		"Provides help for available commands",
		"Examples:\n`help` will show help for all available commands\n`help <command>` will show help for specific command",
		commander.WithPrefixMatcher(),
		commander.WithHandler(func(data *slackevents.MessageEvent) error {
			cmd := strings.TrimSpace(data.Text[len("help"):])
			slackBot.SendHelpMessage(
				data.Channel,
				message.HelpMessage(),
				cmd,
			)
			return nil
		}),
	))

	return slackBot, nil
}

func (slackBot *SlackBot) Listen() {
	http.Handle(
		"/slashcmd",
		handler.NewSlashCmdHandler(handler.SlashCmdHandlerConfig{
			SlashCmdChan:      slackBot.slashCmdChan,
			VerificationToken: slackBot.config.GetVerificationToken(),
		}),
	)
	http.Handle(
		"/event",
		handler.NewEventHandler(handler.EventHandlerConfig{
			EventChan:         slackBot.eventChan,
			VerificationToken: slackBot.config.GetVerificationToken(),
		}),
	)
	http.Handle(
		"/interaction",
		handler.NewInteractionHandler(handler.InteractionHandlerConfig{
			ActionChan:        slackBot.actionChan,
			SubmissionChan:    slackBot.submissionChan,
			VerificationToken: slackBot.config.GetVerificationToken(),
		}),
	)

	go slackBot.handleSlashCmds()
	go slackBot.handleEvents()
	go slackBot.handleActions()
	go slackBot.handleSubmissions()

	go slackBot.handleOutgoingMessages()

	slackBot.listenAndServe()
}

func (slackBot SlackBot) listenAndServe() {
	log.Info("Listening on port ", slackBot.config.Port)
	err := http.ListenAndServe(slackBot.config.Port, nil)
	if err != nil {
		log.Fatal("Error starting slack events listener: ", err)
	}
}

func (slackBot SlackBot) handleSlashCmds() {
	for d := range slackBot.slashCmdChan {
		cmd, err := slackBot.Commander.Match(d.Command)
		if err != nil {
			slackBot.SendHelpMessage(d.ChannelID, message.BotDidNotUnderstandMessage(), "")
		}
		cmd.Handle(utils.CmdToMessageEvent(d))
		if err != nil {
			slackBot.SendHelpMessage(d.ChannelID, err.Error(), "")
		}
	}
}

func (slackBot SlackBot) handleEvents() {
	for d := range slackBot.eventChan {
		if slackBot.config.EventMatcher(d) {
			err := slackBot.Commander.Handle(d)
			if err != nil {
				slackBot.SendHelpMessage(d.Channel, message.BotDidNotUnderstandMessage(), "")
			}
		}
	}
}

func (slackBot SlackBot) handleActions() {
	for d := range slackBot.actionChan {
		err := slackBot.Interactor.Handle(d)
		if err != nil {
			slackBot.SendHelpMessage(d.User.ID, err.Error(), "")
		}
	}
}

func (slackBot SlackBot) handleSubmissions() {
	for d := range slackBot.submissionChan {
		err := slackBot.Submitter.Handle(d)
		if err != nil {
			slackBot.SendHelpMessage(d.User.ID, err.Error(), "")
		}
	}
}

func (slackBot SlackBot) handleOutgoingMessages() {
	for m := range slackBot.outgoingMessages {
		slackBot.SlackClient.PostMessage(m.Channel, m.Options...)
	}
}

func (slackBot SlackBot) GetUser(userID string) (*slack.User, error) {
	return slackBot.SlackClient.GetUserInfo(userID)
}

func (slackBot SlackBot) SendMessage(channel string, options ...slack.MsgOption) {
	withBotID := append(options, slack.MsgOptionUser(slackBot.config.BotID), slack.MsgOptionAsUser(true))
	slackBot.outgoingMessages <- message.New(channel, withBotID...)
}

func (slackBot SlackBot) SendSimpleMessage(channel, msg string) {
	slackBot.SendMessage(channel, slack.MsgOptionText(msg, false))
}

func (slackBot SlackBot) SendHelpMessage(channel, msg, command string) {
	slackBot.SendMessage(
		channel,
		slackBot.Commander.HelpMessageOptions(msg, command)...,
	)
}
