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

// Hit adds a card to the player's active hand.
func (p *Player) Hit(card deck.Card) {
	p.hands.hit(card)
}

// Stand
func (p *Player) Stand() {
	p.hands.halt()
}

func NewPlayer(cards []deck.Card, bet int, opts ...func(p *Player) *Player) *Player {
	p := &Player{
		hands: newHands(cards, bet),
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}

func WithName(name string) func(p *Player) *Player {
	return func(p *Player) *Player {
		p.Name = name
		return p
	}
}
