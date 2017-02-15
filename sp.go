package smartpong

import (
	"log"
	"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
)

var World *SmartPong

// Player ...
type Player struct {
	ID     string
	Name   string
	Avatar string
	Score  uint8
}

// Game ...
type Game struct {
	Player1    *Player
	Player2    *Player
	StartTime  time.Time
	InProgress bool
}

// GameState ...
type GameState struct {
	InProgress bool
	//stuff from yulin goes here
}

// SmartPong ...
type SmartPong struct {
	Games *NintendoGameQueue
	State *GameState
}

func init() {
	World = &SmartPong{
		//change this back
		NewntendoGameQueue(),
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

func MonitorInputs() {
	embd.InitGPIO()
	defer embd.CloseGPIO()

	startPin, err := embd.NewDigitalPin("GPIO_14")
	if err != nil {
		log.Fatal(err)
	}
	startPin.SetDirection(embd.In)
	defer startPin.Close()

	player1Pin, err := embd.NewDigitalPin("GPIO_15")
	if err != nil {
		log.Fatal(err)
	}
	player1Pin.SetDirection(embd.In)
	defer player1Pin.Close()

	player2Pin, err := embd.NewDigitalPin("GPIO_18")
	if err != nil {
		log.Fatal(err)
	}
	player2Pin.SetDirection(embd.In)
	defer player2Pin.Close()

	defaultSv := 0
	defaultPlay1 := 1
	defaultPlay2 := 1

	for {
		sv, err := startPin.Read()
		if err != nil {
			log.Println(err)
		}
		if defaultSv != sv {
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
