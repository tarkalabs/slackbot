package main

import (
	"log"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/tarkalabs/slackbot"
	"github.com/tarkalabs/slackbot/commander"
	"github.com/tarkalabs/slackbot/interactor"
	"github.com/tarkalabs/slackbot/submitter"
)

func newEntryDialog() *slack.Dialog {
	d := slack.Dialog{}
	return &d
}

func newEntryPostMessageParameters() *slack.PostMessageParameters {
	p := slack.NewPostMessageParameters()
	return &p
}

func main() {
	slackBot, err := slackbot.New(
		slackbot.SlackConfig{
			Port:              ":63800",
			BotID:             os.Getenv("SLACK_BOT_ID"),
			APIToken:          os.Getenv("SLACK_API_TOKEN"),
			VerificationToken: os.Getenv("SLACK_VERIFICATION_TOKEN"),
		},
		slackbot.WithEventMatcher(func(data *slackevents.MessageEvent) bool {
			// filter IMs
			return data.ChannelType == "im" && data.BotID == "" && data.SubType == ""
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	slackBot.Commander.Add(commander.NewCommand(
		"/add",
		"Add new entry",
		"Will open a dialog to enter your data",
		commander.WithEqualMatcher(),
		commander.WithHandler(func(data *slackevents.MessageEvent) error {
			slackBot.SendSimpleMessage(
				data.Channel,
				":laughing: Oh my, sure! Please click the button below to proceed",
			)
			return nil
		}),
	))

	slackBot.Interactor.Add(interactor.NewInteraction(
		"new_entry",
		interactor.WithHandler(func(action *slackevents.MessageAction) error {
			log.Println(action.Type)
			dialog := newEntryDialog()
			return slackBot.SlackClient.OpenDialog(action.TriggerID, *dialog)
		}),
	))

	slackBot.Submitter.Add(submitter.NewSubmission(
		"entry_submission",
		submitter.WithHandler(func(submission *slack.DialogCallback) error {
			_, err := slackBot.GetUser(submission.User.ID)
			if err != nil {
				return err
			}
			// Do Work
			slackBot.SendSimpleMessage(
				submission.User.ID,
				"Recorded. You are awesome :hungging_face:",
			)
			return nil
		}),
	))

	slackBot.Listen()
}
