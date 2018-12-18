package commander

import (
	"reflect"
	"testing"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"github.com/tarkalabs/slackbot/message"
)

var commandtests = []struct {
	command       Command
	expectedCount int
}{
	{
		NewCommand(
			"add",
			"Add new entry",
			"Will open a dialog to enter your timesheet data",
			WithEqualMatcher(),
			WithHandler(func(data *slackevents.MessageEvent) (message.Message, error) {
				return message.Message{}, nil
			}),
		), 1,
	},
	{
		NewCommand(
			"quickadd",
			"Quick Add",
			"For example, when you type *quickadd 12/09/2018 | Client1 | 4 | Worked on 3 stories* \nI will record *4 hours* of effort on *12 Sep 2018* for Client *Client1* with a note *Worked on 3 stories*",
			WithPrefixMatcher(),
			WithHandler(func(data *slackevents.MessageEvent) (message.Message, error) {
				return message.Message{}, nil
			}),
		), 2,
	},
}

func buildCommander() Commander {
	commander := Commander{}
	for _, c := range commandtests {
		commander.Add(c.command)
	}
	return commander
}

func TestAdd(t *testing.T) {
	commander := Commander{}
	for _, c := range commandtests {
		commander.Add(c.command)
		if len(commander.commands) != c.expectedCount {
			t.Errorf("Got: %d, Expected: %d", len(commander.commands), c.expectedCount)
		}
	}
}

func TestAddDuplicate(t *testing.T) {
	commander := buildCommander()
	lenBefore := len(commander.commands)
	commander.Add(Command{
		Name: "add",
	})
	if len(commander.commands) != lenBefore {
		t.Errorf("Duplicate Command should not get added. Got length: %d, expected: %d", len(commander.commands), lenBefore)
	}
}

func TestLength(t *testing.T) {
	commander := buildCommander()
	if commander.Length() != 2 {
		t.Errorf("Length is not accurate. Got: %d, expected: %d", commander.Length(), 2)
	}
}

func TestFind(t *testing.T) {
	commander := buildCommander()
	_, err := commander.Find("add")
	if err != nil {
		t.Errorf("Failed to find added 'add' command")
	}
}

func TestFindInvalid(t *testing.T) {
	commander := Commander{}
	_, err := commander.Find("add")
	if err == nil {
		t.Errorf("Failed command should result in InvalidCommandError")
	}
}

func TestMatch(t *testing.T) {
	commander := buildCommander()
	_, err := commander.Match("add")
	if err != nil {
		t.Errorf("Failed to match equal 'add' command")
	}
	_, err = commander.Match("quickadd 12/09/2018 | Client1 | 4 | Worked on 3 stories")
	if err != nil {
		t.Errorf("Failed to match prefix 'quickadd' command %v", err)
	}
}

func TestMatchInvalid(t *testing.T) {
	commander := Commander{}
	_, err := commander.Match("add")
	if err == nil {
		t.Error("Expected to raise Invalid Command Error")
	}
}

func TestHelp(t *testing.T) {
	commander := buildCommander()
	helpstring := "*`add` - Add new entry* \n Will open a dialog to enter your timesheet data" +
		"\n\n" +
		"*`quickadd` - Quick Add* \n For example, when you type *quickadd 12/09/2018 | Client1 | 4 | Worked on 3 stories* \nI will record *4 hours* of effort on *12 Sep 2018* for Client *Client1* with a note *Worked on 3 stories*" +
		"\n\n" +
		"*`help` - List all the commands* \n "
	if commander.Help() != helpstring {
		t.Errorf("Help does not match. Got: %s, Expected: %s", commander.Help(), helpstring)
	}
}

func TestHelpMessage(t *testing.T) {
	commander := buildCommander()
	msgStr := "test help message"
	helpstring := "*`add` - Add new entry* \n Will open a dialog to enter your timesheet data" +
		"\n\n" +
		"*`quickadd` - Quick Add* \n For example, when you type *quickadd 12/09/2018 | Client1 | 4 | Worked on 3 stories* \nI will record *4 hours* of effort on *12 Sep 2018* for Client *Client1* with a note *Worked on 3 stories*" +
		"\n\n" +
		"*`help` - List all the commands* \n "
	msg := message.Message{
		Message: msgStr,
		Body: &slack.PostMessageParameters{
			Attachments: []slack.Attachment{
				{
					Text:       helpstring,
					Color:      "#25CCF7",
					MarkdownIn: []string{"text"},
				},
			},
		},
	}
	if !reflect.DeepEqual(commander.HelpMessage(msgStr), msg) {
		t.Errorf("HelpMessage doesn't correspond to msg. Got: %v, Expected: %v", commander.HelpMessage(msgStr), msg)
	}
}

func TestHelpHandle(t *testing.T) {
	commander := Commander{}
	_, err := commander.Handle(&slackevents.MessageEvent{
		Text: "help",
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestAddHandle(t *testing.T) {
	commander := buildCommander()
	msg, err := commander.Handle(&slackevents.MessageEvent{
		Text: "add",
	})
	if err != nil {
		t.Errorf("Unexpected error from Add Command handler: %v", err)
	}
	if (msg != message.Message{}) {
		t.Errorf("Invalid msg from Add Command handler. Got: %v, Expected: %v", msg, message.Message{})
	}
}

func TestInvalidHandle(t *testing.T) {
	commander := Commander{}
	_, err := commander.Handle(&slackevents.MessageEvent{
		Text: "add",
	})
	if err == nil {
		t.Errorf("Expected error from Add Command handler")
	}
}
