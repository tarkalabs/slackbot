package slackbot

import (
	"net/http"
	"strings"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	log "github.com/sirupsen/logrus"
	"github.com/tarkalabs/slackbot/commander"
	"github.com/tarkalabs/slackbot/handler"
	"github.com/tarkalabs/slackbot/interactor"
	"github.com/tarkalabs/slackbot/message"
	"github.com/tarkalabs/slackbot/submitter"
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

	slackClient := slack.New(config.APIToken)
	slackClient.SetDebug(true)

	slackBot := &SlackBot{
		config:      config,
		SlackClient: slackClient,
	}
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
			slackBot.SendHelpMessage(d.User.Id, err.Error(), "")
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
		m.Body.Username = slackBot.config.BotID
		m.Body.AsUser = true
		slackBot.SlackClient.PostMessage(m.Channel, m.Message, *m.Body)
	}
}

func (slackBot SlackBot) GetUser(userID string) (*slack.User, error) {
	return slackBot.SlackClient.GetUserInfo(userID)
}

func (slackBot SlackBot) SendMessage(message *message.Message) {
	slackBot.outgoingMessages <- message
}

func (slackBot SlackBot) SendHelpMessage(channel, msg, command string) {
	slackBot.SendMessage(
		slackBot.Commander.HelpMessage(
			channel,
			msg,
			command,
		),
	)
}
