package main

import (
	"fmt"
	"sync"
)

// GameQueue ...
type GameQueue struct {
	*sync.RWMutex
	games []*Game
}

func NewGameQueue() *GameQueue {
	return &GameQueue{
		&sync.RWMutex{},
		make([]*Game, 0),
	}
}

// Push ...
func (n *GameQueue) Push(g *Game) (*Game, error) {
	n.Lock()
	defer n.Unlock()

	n.games = append(n.games, g)

	return g, nil
}

// Pop ...
func (n *GameQueue) Pop() (*Game, error) {
	n.Lock()
	defer n.Unlock()

	if len(n.games) > 0 {
		tmp := n.games[len(n.games)-1]
		n.games = n.games[:len(n.games)-1]

		return tmp, nil
	}

	return nil, fmt.Errorf("no games in list, can't pop")
}

// Peek ...
func (n *GameQueue) Peek() (*Game, error) {
	n.RLock()
	defer n.RUnlock()

	if n.Len() > 0 {
		return n.games[0], nil
	}

	return nil, fmt.Errorf("no games in list, can't peek")
}

// Len ...
func (n *GameQueue) Len() int {
	n.RLock()
	defer n.RUnlock()

	return len(n.games)
}
