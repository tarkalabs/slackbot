# SlackBot [![Build Status](https://travis-ci.org/tarkalabs/slackbot.svg?branch=master)](https://travis-ci.org/tarkalabs/slackbot)

SlackBot is a framework for building Slack Bots utilizing the Slack APIso

## Installing

_**go get**_

```sh
$ go get github.com/tarkalabs/slackbot
```

## Example

```go
package main

import (
  "github.com/tarkalabs/slackbot"
  "github.com/tarkalabs/slackbot/commander"
  "github.com/tarkalabs/slackbot/interactor"
  "github.com/tarkalabs/slackbot/message"
  "github.com/tarkalabs/slackbot/submitter"
  "github.com/tarkalabs/slackbot/utils"
)

func main() {
  slackBot, err := slackbot.New(
    slackbot.SlackConfig{
      Port:              getPort(),
      BotID:             os.Getenv("SLACK_BOT_ID"),
      APIToken:          os.Getenv("SLACK_API_TOKEN"),
      VerificationToken: os.Getenv("SLACK_VERIFICATION_TOKEN"),
    },
    slackbot.WithEventMatcher(func (data *slackevents.MessageEvent) bool {
      // filter IMs
      return data.ChannelType == "im" && data.BotID == "" && data.SubType == ""
    },
  )
  if err != nil {
    log.Fatal(err)
  }

  slackBot.Commander.Add(commander.AddCommand(
    "add",
    "Add new entry",
    "Will open a dialog to enter your data",
    commander.WithEqualMatcher(),
    commander.WithHandler(func(data *slackevents.MessageEvent) error {
      slackBot.SendMessage(message.New(
        data.Channel,
        ":laughing: Oh my, sure! Please click the button below to proceed",
        newEntryPostMessageParameters(),
      ))
      return nil
    })
  ))

  slackBot.Interactor.Add(interactor.NewInteraction(
    "new_entry",
    interactor.WithHandler(func(action *slackevents.MessageAction) error {
      dialog := newEntryDialog()
      return slackBot.SlackClient.OpenDialog(action.TriggerId, *dialog)
    })
  ))

  slackBot.Submitter.Add(submitter.NewSubmission(
    "entry_submission",
    submitter.WithHandler(func(submission *slack.DialogCallback) error {
      slackUser, err := slackBit.GetUser(submission.User.ID)
      if err != nil {
        return err
      }
      // Do Work
      slackBot.SendMessage(message.New(
        submission.User.ID,
        "Recorded. You are awesome :hungging_face:",
        utils.GetPostMessage(""),
      ))
      return nil
    })
  ))

  slackBot.Listen()
}
```
