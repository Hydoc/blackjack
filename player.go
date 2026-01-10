package blackjack

import "github.com/Hydoc/deck"

type Mode = int

const (
	normal Mode = iota
	split
)

type Player struct {
	Name string

	hands *Hands
}

func (p *Player) Hit(card deck.Card) {
	p.hands.Hit(card)
}

func (p *Player) Halt() {
	p.hands.Halt()
}

func NewPlayer(cards []deck.Card, bet int, opts ...func(p *Player) *Player) *Player {
	p := &Player{
		hands: NewHands(cards, bet),
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
