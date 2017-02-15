package smartpong

import "time"

// Booking ...
type Booking struct {
	Player    *Player
	Opponent  *Player
	StartTime time.Time
}

// Promote ... once validated from pipo, turn the booking into an actual game
func (b *Booking) Promote() {
	World.Games.Push(&Game{
		Player1:   b.Player,
		Player2:   b.Opponent,
		StartTime: b.StartTime,
	})
}
