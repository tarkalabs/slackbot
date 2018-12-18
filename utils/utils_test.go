package utils

import (
	"reflect"
	"testing"

	"github.com/nlopes/slack"
)

func TestGetPostMessage(t *testing.T) {
	got := GetPostMessage("test")
	expected := &slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			{
				Text:       "test",
				Color:      "#25CCF7",
				MarkdownIn: []string{"text"},
			},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Got: %v, Expected: %v", got, expected)
	}
}

func TestGetPostMessageEmpty(t *testing.T) {
	if !reflect.DeepEqual(GetPostMessage(""), &slack.PostMessageParameters{}) {
		t.Errorf("Empty message should produce empty PostMessageParameters")
	}
}
