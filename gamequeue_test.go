package smartpong_test

import (
	"log"
	"testing"
	"time"

	"git-lab.boldapps.net/ahoff/smart-pong/smartpong"
)

var World *smartpong.SmartPong

func init() {
	World = &smartpong.SmartPong{
		//change this back
		smartpong.NewntendoGameQueue(),
		&smartpong.GameState{},
	}

	go smartpong.StartDaPeeker()

}

func TestPush(t *testing.T) {
	g, err := World.Games.Push(&smartpong.Game{
		Player1:   &smartpong.Player{Name: "Andrew"},
		Player2:   &smartpong.Player{Name: "Eric"},
		StartTime: time.Now().Add(4 * time.Minute),
	})
	if err != nil {
		t.Fatal(err)
	}

	log.Println(g)
}

func TestLen(t *testing.T) {
	log.Println(World.Games.Len())
	time.Sleep(2 * time.Minute)
}

func TestPeek(t *testing.T) {
	g, err := World.Games.Peek()
	if err != nil {
		t.Fatal(err)
	}

	log.Println(g)
}

func TestPop(t *testing.T) {
	g, err := World.Games.Pop()
	if err != nil {
		t.Fatal(err)
	}

	log.Println(g)
}

func TestEmpty(t *testing.T) {
	if World.Games.Len() != 0 {
		t.Fail()
	}
}
