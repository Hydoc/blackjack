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

	wallet int

	hands *hands
}

// CanSplit returns a bool whether the player can split.
func (p *Player) CanSplit() bool {
	return p.canBetTheSameAmountAgain() && p.hands.canSplit()
}

// CanDoubleDown returns a bool whether the player can double down.
func (p *Player) CanDoubleDown() bool {
	return p.canBetTheSameAmountAgain() && p.hands.canDoubleDown()
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

func (p *Player) canBetTheSameAmountAgain() bool {
	return p.hands.active.bet <= p.wallet
}

func (p *Player) hasBlackJack() bool {
	return p.hands.hasBlackJack()
}

func (p *Player) busted() bool {
	return p.hands.busted()
}

func (p *Player) isDone() bool {
	return p.hands.isDone()
}

// NewPlayer creates a pointer to the new player with the passed configuration.
func NewPlayer(wallet int, opts ...func(p *Player) *Player) *Player {
	p := &Player{
		hands:  newHands(),
		wallet: wallet,
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
