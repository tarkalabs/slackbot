package submitter

import (
	"github.com/slack-go/slack"
)

type SubmissionHandler func(*slack.DialogCallback) error

type SubmissionOption struct {
	Name    string
	Handler SubmissionHandler
}
type SubmissionOptions func(*SubmissionOption)

func WithHandler(handler SubmissionHandler) SubmissionOptions {
	return func(o *SubmissionOption) {
		o.Handler = handler
	}
}

type Submission struct {
	Name    string
	Handler SubmissionHandler
}

func NewSubmission(name string, opts ...SubmissionOptions) Submission {
	options := SubmissionOption{Name: name}
	for _, o := range opts {
		o(&options)
	}
	return Submission{
		Name:    name,
		Handler: options.Handler,
	}
}

func (su *Submission) Handle(submission *slack.DialogCallback) error {
	return su.Handler(submission)
}
