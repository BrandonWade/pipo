package main

import "time"

// Player - contains information about a player
type Player struct {
	ID     string
	Name   string
	Avatar string
	Score  uint8
}

// Game - contains information about a game
type Game struct {
	Player1    *Player
	Player2    *Player
	StartTime  time.Time
	InProgress bool
}

// GameList - a list of games
type GameList []*Game

// Len ...
func (games GameList) Len() int {
	return len(games)
}

// Less ...
func (games GameList) Less(i, j int) bool {
	return games[i].StartTime.Before(games[j].StartTime)
}

// Swap ...
func (games GameList) Swap(i, j int) {
	games[i], games[j] = games[j], games[i]
}
