package submitter

import (
	"github.com/nlopes/slack"
)

type SubmissionHandler func(*slack.DialogCallback) error

type SubmissionOption func(*SubmissionOptions)
type SubmissionOptions struct {
	Name    string
	Handler SubmissionHandler
}

func WithHandler(handler SubmissionHandler) SubmissionOption {
	return func(o *SubmissionOptions) {
		o.Handler = handler
	}
}

type Submission struct {
	Name    string
	Handler SubmissionHandler
}

func NewSubmission(name string, opts ...SubmissionOption) Submission {
	options := SubmissionOptions{Name: name}
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
