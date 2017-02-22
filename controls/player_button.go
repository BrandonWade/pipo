package controls

import (
	"fmt"

	"github.com/kidoman/embd"
)

// NewP1Button ... gets a new start led pin
func NewP1Button() (embd.DigitalPin, error) {

	pin, err := embd.NewDigitalPin(p1PinName)
	if err != nil {
		return nil, fmt.Errorf("Error encountered while trying to create p1 button pin: %v", err)
	}
	pin.SetDirection(embd.Out)

	return pin, nil
}

// NewP2Button ...
func NewP2Button() (embd.DigitalPin, error) {

	pin, err := embd.NewDigitalPin(p2PinName)
	if err != nil {
		return nil, fmt.Errorf("Error encountered while trying to create p2 button pin: %v", err)
	}
	pin.SetDirection(embd.Out)

	return pin, nil
}
