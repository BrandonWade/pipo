package main_test

//
// import (
// 	"log"
// 	"testing"
// 	"time"
// )
//
// var World *smartpong.SmartPong
//
// func init() {
// 	World = &smartpong.SmartPong{
// 		//change this back
// 		smartpong.NewGameQueue(),
// 		&smartpong.GameState{},
// 	}
//
// 	go smartpong.StartDaPeeker()
//
// }
//
// func TestPush(t *testing.T) {
// 	g, err := World.Games.Push(&smartpong.Game{
// 		Player1:   &smartpong.Player{Name: "Andrew"},
// 		Player2:   &smartpong.Player{Name: "Eric"},
// 		StartTime: time.Now().Add(4 * time.Minute),
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	log.Println(g)
// }
//
// func TestLen(t *testing.T) {
// 	log.Println(World.Games.Len())
// 	time.Sleep(2 * time.Minute)
// }
//
// func TestPeek(t *testing.T) {
// 	g, err := World.Games.Peek()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	log.Println(g)
// }
//
// func TestPop(t *testing.T) {
// 	g, err := World.Games.Pop()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	log.Println(g)
// }
//
// func TestEmpty(t *testing.T) {
// 	if World.Games.Len() != 0 {
// 		t.Fail()
// 	}
// }
