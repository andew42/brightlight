package controller

import "time"

type relayState struct {
	required         bool
	current          bool
	lastOffTime      time.Time
	oldestOffRequest time.Time
}

const turnOffDelay = time.Second * 30
const turnOnDelay = time.Second * 1

// Request a state (on or off)
func (state *relayState) requestState(requiredState bool, now time.Time) {

	if requiredState == state.required {
		return
	}

	state.required = requiredState

	if requiredState {
		state.oldestOffRequest = time.Time{}
	} else {
		state.oldestOffRequest = now
	}
}

// What state should the relays be in now, return true if state change required
func (state *relayState) updateState(now time.Time) bool {

	if state.required == state.current {
		return false
	}

	// Off -> On
	if state.required && now.Before(state.lastOffTime.Add(turnOnDelay)) {
		return false
	}

	// On -> Off
	if !state.required && now.Before(state.oldestOffRequest.Add(turnOffDelay)) {
		return false
	}

	state.current = state.required
	state.lastOffTime = now
	return true
}
