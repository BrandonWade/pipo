package controls

import (
	"fmt"

	"github.com/kidoman/embd"
)

// NewStartButton ... gets a new start led pin
func NewStartButton() (embd.DigitalPin, error) {
	pin, err := embd.NewDigitalPin(startBtnPinDesig)
	if err != nil {
		return nil, fmt.Errorf("Error encountered while trying to create start button pin: %v", err)
	}
	pin.SetDirection(embd.In)

	return pin, nil
}
