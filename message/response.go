package message

import (
	"math/rand"
	"time"
)

var didNotUnderstand = []string{
	":astonished: Surprisingly, I don't understand that. Can you try something from the below list?",
	":thinking_face: Hmm? I remember hearing this somewhere before, but don't understand it now. Can you try something from the below list?",
	"I don't speek greek :exploding_head:",
	"I haven't added support for that yet! Try something else? :crying_cat_face:",
	"Missile launched successfully!.. Just kidding :smiley:, have no clue what you mean. Try one of the commands below",
}

var help = []string{
	"I understand you are lost :shrug:. You can use these commands",
	"Don't worry, Let me help you :nerd_face: Use any of the following commands",
	":genie: Your wish is my command. Just that your wish should be one of the following :face_with_rolling_eyes:",
}

var readyToAcceptEntryAck = []string{
	":laughing: Oh my, sure! Please click the button below to proceed",
	":female-office-worker: You do? Sweet. Click the button and we're in business",
	"Aren't you a role model to volunteer! Here's the sweet button for you, go ahead, press it :point_down:",
}

var entryAddedAck = []string{
	"Recorded. You are awesome :hugging_face:",
	"Done! Give me more :kissing:",
	"Saved it. Way to go! :heart:",
}

var reminder = []string{
	"ding ding ding! :bellhop_bell: Time to fill Timesheet!",
	":guardsman: Pardon me Sire! But I must disturb you to remind you to fill your Timesheet for the day!",
	":hourglass_flowing_sand: Timesheet Time! hey hey, I'm just a messenger. Don't get annoyed on me..",
	":timer_clock: Timesheet Time! You must be thinking, \"hmm, I don't have time for this, I'll do it later.\" Trust me, it's better if you get it over with sooner..",
	":man-cartwheeling: The secret to my happiness is I fill my Timesheet on Time. Now it's your turn!",
}

// HelpMessage returns a string to respond when help is requested
func HelpMessage() string {
	return pickRandom(help)
}

// BotDidNotUnderstandMessage returns a string to respond if the bot did not understand the command
func BotDidNotUnderstandMessage() string {
	return pickRandom(didNotUnderstand)
}

// ReadyToAcceptEntryAckMessage returns a string to respond when the bot is ready to accept new entries
func ReadyToAcceptEntryAckMessage() string {
	return pickRandom(readyToAcceptEntryAck)
}

// EntryAddedAckMessage returns a string to respond when the bot added the entry successfully
func EntryAddedAckMessage() string {
	return pickRandom(entryAddedAck)
}

func ReminderMessage() string {
	return pickRandom(reminder)
}

func pickRandom(messages []string) string {
	rand.Seed(time.Now().Unix())
	n := rand.Intn(len(messages))

	return messages[n]
}
