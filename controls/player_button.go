package controls

import (
	"fmt"

	"github.com/kidoman/embd"
)

// NewP1Button ... gets a new start led pin
func NewP1Button() (embd.DigitalPin, error) {
	return newPlayerButton(p1PinDesig)
}

// NewP2Button ...
func NewP2Button() (embd.DigitalPin, error) {
	return newPlayerButton(p2PinDesig)
}

func newPlayerButton(pinDesignator string) (embd.DigitalPin, error) {
	pin, err := embd.NewDigitalPin(pinDesignator)
	if err != nil {
		return nil, fmt.Errorf("Error encountered while trying to create a player button for pin %s: %v", pinDesignator, err)
	}
	pin.SetDirection(embd.Out)

	return pin, nil
}
