package main

import (
	"log"

	"github.com/BrandonWade/pipo/controls"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
)

func monitor() {
	embd.InitGPIO()
	defer embd.CloseGPIO()

	startBtn, err := controls.NewStartButton()
	if err != nil {
		log.Fatal(err)
	}
	defer startBtn.Close()

	startLEDPin, err := controls.NewStartLED()
	if err != nil {
		log.Fatal(err)
	}
	defer startLEDPin.Close()

	p1Btn, err := controls.NewP1Button()
	if err != nil {
		log.Fatal(err)
	}
	defer p1Btn.Close()

	p2Btn, err := controls.NewP2Button()
	if err != nil {
		log.Fatal(err)
	}
	defer p2Btn.Close()

	go Blink(startLEDPin, blinkChan)

	for {
		sv, err := startPin.Read()
		if err != nil {
			log.Println(err)
		}
		if defaultSv != sv {
			// stop the blinking
			blinkChan <- 1
			//
			log.Println("startpin")
		}

		p1v, err := player1Pin.Read()
		if err != nil {
			log.Println(err)
		}
		if defaultPlay1 != p1v {
			log.Println("p1")
		}

		p2v, err := player2Pin.Read()
		if err != nil {
			log.Println(err)
		}
		if defaultPlay2 != p2v {
			log.Println("p2")
		}
	}
}
