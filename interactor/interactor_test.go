package interactor

import (
	"testing"
)

func buildInteractor() Interactor {
	interactor := Interactor{}
	interactor.Add(NewInteraction("test"))
	return interactor
}

func TestFind(t *testing.T) {
	interactor := buildInteractor()
	_, err := interactor.Find("test")
	if err != nil {
		t.Errorf("Unable to find new_entry interaction")
	}
}

func TestFindInvalid(t *testing.T) {
	interactor := Interactor{}
	_, err := interactor.Find("new_entry")
	if err == nil {
		t.Errorf("Invalid interaction should raise InvalidInteractionError")
	}
}

var interactortests = []struct {
	interaction   Interaction
	expectedCount int
}{
	{
		NewInteraction("test1"),
		1,
	},
	{
		NewInteraction("test2"),
		2,
	},
}

func TestAdd(t *testing.T) {
	interactor := Interactor{}
	for _, in := range interactortests {
		interactor.Add(in.interaction)
		if len(interactor.interactions) != in.expectedCount {
			t.Errorf("Got: %d, Expected: %d", len(interactor.interactions), in.expectedCount)
		}
	}
}

func TestAddDuplicate(t *testing.T) {
	interactor := buildInteractor()
	lenBefore := len(interactor.interactions)
	interactor.Add(NewInteraction("test"))
	if len(interactor.interactions) != lenBefore {
		t.Errorf("Got: %d, Expected: %d", len(interactor.interactions), lenBefore)
	}
}
