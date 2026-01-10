package blackjack

import (
	"errors"

	"github.com/Hydoc/deck"
)

var (
	ErrTableFull = errors.New("table is full")
)

// Table represents a blackjack table.
type Table struct {
	dealer  *Player
	players [7]*Player
	deck    []deck.Card
}

// Join adds a player to the next nil value in the players slice.
// It returns ErrTableFull when there is no space left.
func (t *Table) Join(p *Player) error {
	for i := range t.players {
		if t.players[i] == nil {
			t.players[i] = p
			return nil
		}
	}
	return ErrTableFull
}

// New create a pointer to Table with the default configuration like 6 shuffled decks and a maximum of 7 players allowed.
func New() *Table {
	return &Table{
		dealer:  NewPlayer(make([]deck.Card, 0), 0, WithName("Dealer")),
		players: [7]*Player{},
		deck: deck.New(
			deck.WithDecks(6),
			deck.Shuffle,
		),
	}
}
