package controller

import (
	"errors"
	"github.com/andew42/brightlight/config"
	"github.com/andew42/brightlight/framebuffer"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var relayDriverStarted bool

// StartRelayDriver Run driver as a go routine
func StartRelayDriver() {

	if relayDriverStarted {
		log.Panic("Relay driver started twice")
	}
	relayDriverStarted = true

	// Start driver
	go relayDriver()
}

func IsRelayDriverConnected() bool {

	return relayUsbConnected
}

var relayUsbConnected bool

// Monitors changes to frame buffer and turns power supplies on or off via USB relay board
func relayDriver() {

	port := getPortName(relayPortMappings, 0)
	if port == "" {
		log.Warn("relayDriver unknown port name")
		return
	}

	// TODO: On Raspberry Pi if we unplug and replug the relay controller
	// before we notice (which can be minutes) a new USB port name is assigned
	// for the controller because the exiting one hasn't yet been closed.
	// i.e. /dev/ttyUSB0 -> /dev/ttyUSB1 This can't be fixed by updating the
	// controller even if it hasn't changed because the relay is quickly pulsed
	// off before setting the correct state

connectLoop:
	for {
		relayUsbConnected = false
		f := openUsbPortWithRetry(port)
		relayUsbConnected = true

		// Initially we set all relays off
		relayStates := [2]relayState{}
		if err := initRelayBoard(f); err != nil {
			log.WithField("error", err.Error()).Warn("relayDriver failed to initialise relay board")
			continue connectLoop
		}

		// Request frame buffer updates
		src, done := framebuffer.AddListener(port, false)

		// Push frame buffer changes to Teensy
	pushLoop:
		for {
			select {
			case fb := <-src:
				now := time.Now()
				// foreach controller
				controllerCount := len(fb.Strips) / config.StripsPerTeensy
				for c := 0; c < controllerCount; c++ {
					firstStripForController := c * config.StripsPerTeensy
					// Should the controller's relay be off or on?
					relayStates[c].requestState(
						areAnyLedsOn(fb.Strips[firstStripForController:firstStripForController+config.StripsPerTeensy]),
						now)
				}

				// Work out the new actual states at this point in time
				newRelayStates := [2]bool{}
				updateRequired := false
				for i := 0; i < len(newRelayStates); i++ {
					updateRequired = relayStates[i].updateState(now) || updateRequired
					newRelayStates[i] = relayStates[i].current
				}

				// Update if changed
				if updateRequired {
					log.WithField("new states", newRelayStates).Info("relayDriver update relays")
					if err := sendRelayState(f, newRelayStates); err != nil {

						log.WithField("error", err.Error()).Warn("relayDriver failed to send relay command")
						f.Close()

						// Close down listener
						done <- src

						// Try and reconnect
						break pushLoop
					}
				}
			}
		}
	}
}

func initRelayBoard(f *os.File) error {

	// After initial power up the board expects 0x50 to which it will
	// reply 0xAD (2 port). After this it expects 0x51 to put it in
	// 'run' mode. Now every byte's lower 2 bits are interpreted as a
	// relay command. Unfortunately we have no way of knowing the
	// current state, and it's not possible to bypass the initial 0x50
	// handshake. So we send 0x50 and if we get a response send 0x51
	// 0x00 to turn off the relays. If we don't get a response we
	// assume the relays are now off, due to the 0x50 and the
	// assumption the board has already been initialised.
	data := make([]byte, 1)
	data[0] = 0x50
	if _, err := f.Write(data); err != nil {
		return err
	}

	// Expect 0xAD for a ICSE013A (2 relay board)
	response := make([]byte, 1)
	if err := readUntilBufferFull(f, response, time.Millisecond*200); err != nil {

		// If we timed out, assume the board is already initialised so the 0x50
		// will be interpreted as turing both relays off, i.e. we are done
		if err == readTimeoutError {
			return nil
		}

		// Otherwise, we have some sort of real error
		return err
	}

	// If we got a response check it's the expected one
	if response[0] != 0xAD {
		log.WithField("response", response[0]).Warn("initRelayBoard unexpect initialisation response")
		return errors.New("unknown relay board")
	}

	// Start 'command mode' (0x51) and turn all relays off (0x00)
	data = make([]byte, 2)
	data[0] = 0x51
	data[1] = 0x00
	_, err := f.Write(data)
	return err
}

func sendRelayState(f *os.File, state [2]bool) error {

	// Transform the relay states into one byte command
	data := make([]byte, 1)
	for i := 0; i < len(state); i++ {
		if state[i] {
			data[0] |= 1 << uint(i)
		}
	}

	// Send the relay commands
	_, err := f.Write(data)
	if err == nil {
		log.WithField("parameter", data[0]).Info("Relay command sent")
	}
	return err
}

func areAnyLedsOn(strips []framebuffer.LedStrip) bool {

	for _, s := range strips {
		for _, l := range s.Leds {
			if l.IsLedOn() {
				return true
			}
		}
	}
	return false
}
