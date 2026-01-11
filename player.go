package blackjack

import "github.com/Hydoc/deck"

type Mode = int

const (
	normal Mode = iota
	split
)

// Player represents one player in the game.
type Player struct {
	Name string

	hands *hands
}

// CanSplit returns a bool whether the player can split.
func (p *Player) CanSplit() bool {
	return p.hands.canSplit()
}

// Hit adds a card to the player's active hand.
func (p *Player) Hit(card deck.Card) {
	p.hands.hit(card)
}

// Stand calls stand on the active hand.
// When playing normal mode the player ends its turn.
// When playing in split mode the second hand will be active after the first.
// Calling Stand again on the second hand ends the split turn.
func (p *Player) Stand() {
	p.hands.stand()
}

func (p *Player) hasBlackJack() bool {
	return p.hands.hasBlackJack()
}

// NewPlayer creates a pointer to the new player with the passed configuration.
func NewPlayer(cards []deck.Card, bet int, opts ...func(p *Player) *Player) *Player {
	p := &Player{
		hands: newHands(cards, withBet(bet)),
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}

// WithName is an option for NewPlayer.
func WithName(name string) func(p *Player) *Player {
	return func(p *Player) *Player {
		p.Name = name
		return p
	}
}
