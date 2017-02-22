package controls

import (
	"fmt"
	"log"
	"time"

	"github.com/kidoman/embd"
)

// NewStartLED ... gets a new start led pin
func NewStartLED() (embd.DigitalPin, error) {

	pin, err := embd.NewDigitalPin(startLEDPinDesig)
	if err != nil {
		return nil, fmt.Errorf("Error encountered while trying to create start LED pin: %v", err)
	}
	pin.SetDirection(embd.Out)

	return pin, nil
}

// Blink ...
func Blink(pin embd.DigitalPin, killChan chan int) error {
	defer startLEDOff(pin)

	togglePin(pin)

	t := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-t.C:
			togglePin(pin)
		case <-killChan:
			return nil
		}
	}
}

func togglePin(pin embd.DigitalPin) {
	state, err := pin.Read()
	if err != nil {
		log.Println(err)
		return
	}

	if state == embd.Low {
		pin.Write(embd.High)
		return
	}

	pin.Write(embd.Low)
}

func startLEDOff(pin embd.DigitalPin) error {
	return pin.Write(embd.Low)
}

func startLEDOn(pin embd.DigitalPin) error {
	return pin.Write(embd.High)
}
