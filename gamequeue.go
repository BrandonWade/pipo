package smartpong

import (
	"fmt"
	"sync"
)

// NintendoGameQueue ...
type NintendoGameQueue struct {
	*sync.RWMutex
	games []*Game
}

func NewntendoGameQueue() *NintendoGameQueue {
	return &NintendoGameQueue{
		&sync.RWMutex{},
		make([]*Game, 0),
	}
}

// Push ...
func (n *NintendoGameQueue) Push(g *Game) (*Game, error) {
	n.Lock()
	defer n.Unlock()

	n.games = append(n.games, g)

	return g, nil
}

// Pop ...
func (n *NintendoGameQueue) Pop() (*Game, error) {
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
func (n *NintendoGameQueue) Peek() (*Game, error) {
	n.RLock()
	defer n.RUnlock()

	if n.Len() > 0 {
		return n.games[0], nil
	}

	return nil, fmt.Errorf("no games in list, can't peek")
}

// Len ...
func (n *NintendoGameQueue) Len() int {
	n.RLock()
	defer n.RUnlock()

	return len(n.games)
}
