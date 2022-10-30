package submitter

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type Submitter struct {
	submissions []Submission
}

var (
	InvalidSubmissionError = errors.New("Invalid Submission")
)

func (s *Submitter) Find(name string) (Submission, error) {
	for _, su := range s.submissions {
		if su.Name == name {
			return su, nil
		}
	}
	return Submission{}, InvalidSubmissionError
}

func (s *Submitter) Add(submission Submission) {
	if _, err := s.Find(submission.Name); err == nil {
		log.Infof("Submission %s already exists", submission.Name)
		return
	}
	s.submissions = append(s.submissions, submission)
}

func (s *Submitter) Handle(submission *slack.DialogCallback) error {
	su, err := s.Find(submission.CallbackID)
	if err != nil {
		return err
	}
	return su.Handle(submission)
}
