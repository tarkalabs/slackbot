package interactor

import (
	"errors"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	log "github.com/sirupsen/logrus"
)

var (
	InvalidInteractionError = errors.New("Invalid Interaction")
)

type Interactor struct {
	interactions []Interaction
}

func (i *Interactor) Find(name string) (Interaction, error) {
	for _, in := range i.interactions {
		if in.Name == name {
			return in, nil
		}
	}
	return Interaction{}, InvalidInteractionError
}

func (i *Interactor) Add(interaction Interaction) {
	if _, err := i.Find(interaction.Name); err == nil {
		log.Infof("Interaction %s already exists", interaction.Name)
		return
	}
	i.interactions = append(i.interactions, interaction)
}

func (i *Interactor) Handle(action *slackevents.MessageAction, client *slack.Client) {
	in, err := i.Find(action.Actions[0].Name)
	if err == nil {
		in.Handle(action, client)
	}
}
