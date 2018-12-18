package slackbot

import (
	"net/http"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	log "github.com/sirupsen/logrus"
	"github.com/tarkalabs/slackbot/commander"
	"github.com/tarkalabs/slackbot/handler"
	"github.com/tarkalabs/slackbot/interactor"
	"github.com/tarkalabs/slackbot/message"
	"github.com/tarkalabs/slackbot/submitter"
)

type SlackConfig struct {
	Port              string
	BotID             string
	APIToken          string
	VerificationToken string
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
	outgoingMessages chan *message.SlackMessage

	Commander  commander.Commander
	Interactor interactor.Interactor
	Submitter  submitter.Submitter
}

func New(config SlackConfig) (*SlackBot, error) {
	slackClient := slack.New(config.APIToken)
	slackClient.SetDebug(true)

	slackBot := &SlackBot{
		config:      config,
		SlackClient: slackClient,
	}
	slackBot.eventChan = make(chan *slackevents.MessageEvent, 50)
	slackBot.actionChan = make(chan *slackevents.MessageAction, 50)
	slackBot.submissionChan = make(chan *slack.DialogCallback, 50)
	slackBot.outgoingMessages = make(chan *message.SlackMessage, 50)

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
		msg, err := slackBot.Commander.Handle(d)
		if err != nil {
			msg = slackBot.Commander.HelpMessage(message.ErrorMessage())
		}
		slackBot.outgoingMessages <- &message.SlackMessage{
			Channel: d.Channel,
			Message: &msg,
		}
	}
}

func (slackBot SlackBot) handleActions() {
	for d := range slackBot.actionChan {
		slackBot.Interactor.Handle(d, slackBot.SlackClient)
	}
}

func (slackBot SlackBot) handleSubmissions() {
	for d := range slackBot.submissionChan {
		msg, err := slackBot.Submitter.Handle(d)
		if err != nil {
			msg = slackBot.Commander.HelpMessage(message.ErrorMessage())
		}
		slackBot.outgoingMessages <- &message.SlackMessage{
			Channel: d.User.ID,
			Message: &msg,
		}
	}
}

func (slackBot SlackBot) handleOutgoingMessages() {
	for m := range slackBot.outgoingMessages {
		m.Message.Body.Username = slackBot.config.BotID
		m.Message.Body.AsUser = true
		slackBot.SlackClient.PostMessage(m.Channel, m.Message.Message, *m.Message.Body)
	}
}

func (slackBot SlackBot) GetUser(userID string) (*slack.User, error) {
	return slackBot.SlackClient.GetUserInfo(userID)
}

func (slackBot SlackBot) SendMessage(message *message.SlackMessage) {
	slackBot.outgoingMessages <- message
}
