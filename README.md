# SlackBot

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
  "github.com/tarkalabs/slackbot/message"
)

func main() {
  slackBot, err := slackbot.New(
    slackbot.SlackConfig{
      Port: getPort(),
      BotID: os.Getenv("SLACK_BOT_ID"),
      APIToken: os.Getenv("SLACK_API_TOKEN"),
      VerificationToken: os.Getenv("SLACK_VERIFICATION_TOKEN"),
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
    commander.WithHandler(func(data *slackevents.MessageEvent) (message.Message, error) {
      return message.Message{
        Message: ":laughing: Oh my, sure! Please click the button below to proceed",
        Body: newEntryPostMessageParameters(),
      }
    })
  ))

  slackBot.Listen()
}
```
