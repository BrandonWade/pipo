package main

import (
	"log"
	"time"
)

var World *SmartPong

// GameState ...
type GameState struct {
	InProgress bool
	//stuff from yulin goes here
}

// SmartPong ...
type SmartPong struct {
	Games *GameQueue
	State *GameState
}

func init() {
	World = &SmartPong{
		//change this back
		NewGameQueue(),
		&GameState{},
	}

	go StartDaPeeker()
}

func StartDaPeeker() {
	t := time.NewTicker(15 * time.Second)
	defer t.Stop()

	for range t.C {
		g, err := World.Games.Peek()
		if err == nil {
			// add threshold so we can tell them before their game actually starts
			if g.StartTime.Add(-3 * time.Minute).Before(time.Now()) {
				log.Println(g)
				log.Println("it's your turn!")
				//pipo.Notify(g)
			}
		}
	}
}
