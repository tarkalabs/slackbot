package interactor

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack/slackevents"
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

func (i *Interactor) Handle(action *slackevents.MessageAction) error {
	in, err := i.Find(action.Actions[0].Name)
	if err != nil {
		return err
	}
	return in.Handle(action)
}
