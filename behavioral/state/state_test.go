package state

import (
	"testing"
)

func TestState(t *testing.T) {

	expect := "State_1" +
		"State_1" +
		"State_2"

	mobile := NewMobileAlert()

	result := mobile.Alert()
	result += mobile.Alert()

	mobile.SetState(&MobileAlertSong{})

	result += mobile.Alert()

	if result != expect {
		t.Errorf("Expect result to equal %s, but %s.\n", expect, result)
	}
}
