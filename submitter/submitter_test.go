package submitter

import (
	"testing"

	"github.com/nlopes/slack"
	"github.com/tarkalabs/slackbot/message"
)

func buildSubmitter() Submitter {
	submitter := Submitter{}
	submitter.Add(NewSubmission(
		"entry_submission",
		WithHandler(func(submission *slack.DialogCallback) (message.Message, error) {
			return message.Message{}, nil
		}),
	))
	return submitter
}

func TestFind(t *testing.T) {
	submitter := buildSubmitter()
	_, err := submitter.Find("entry_submission")
	if err != nil {
		t.Errorf("Unable to find entry_submission Submission")
	}
}

func TestFindInvalid(t *testing.T) {
	submitter := Submitter{}
	_, err := submitter.Find("entry_submission")
	if err == nil {
		t.Errorf("Invalid submission should raise InvalidSubmissionError")
	}
}

var submittertests = []struct {
	submission    Submission
	expectedCount int
}{
	{
		NewSubmission("entry_submission"),
		1,
	},
	{
		NewSubmission("test"),
		2,
	},
}

func TestAdd(t *testing.T) {
	submitter := Submitter{}
	for _, su := range submittertests {
		submitter.Add(su.submission)
		if len(submitter.submissions) != su.expectedCount {
			t.Errorf("Got: %d, Expected: %d", len(submitter.submissions), su.expectedCount)
		}
	}
}

func TestAddDuplicate(t *testing.T) {
	submitter := buildSubmitter()
	lenBefore := len(submitter.submissions)
	submitter.Add(NewSubmission("entry_submission"))
	if len(submitter.submissions) != lenBefore {
		t.Errorf("Got: %d, Expected: %d", len(submitter.submissions), lenBefore)
	}
}

func TestHandle(t *testing.T) {
	submitter := buildSubmitter()
	_, err := submitter.Handle(&slack.DialogCallback{
		CallbackID: "entry_submission",
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestHandleInvalid(t *testing.T) {
	submitter := Submitter{}
	_, err := submitter.Handle(&slack.DialogCallback{
		CallbackID: "entry_submission",
	})
	if err == nil {
		t.Errorf("Expected error from Entry Submission handler")
	}
}
