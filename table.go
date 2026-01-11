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
	dealer          *dealer
	players         [7]*Player
	deck            []deck.Card
	turnPlayerIndex int
}

// Start starts the round at the table by dealing everyone two cards.
func (t *Table) Start() {
	for range 2 {
		for _, p := range t.players {
			if p == nil {
				continue
			}

			cards, remaining := deck.Draw(1)(t.deck)
			t.deck = remaining
			p.Hit(cards[0])
		}

		cards, remaining := deck.Draw(1)(t.deck)
		t.deck = remaining
		t.dealer.Hit(cards[0])
	}
}

func (t *Table) Handle() error {
	return nil
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
// Dealer must stand on soft 17.
func New() *Table {
	return &Table{
		dealer:  newDealer(),
		players: [7]*Player{},
		deck: deck.New(
			deck.WithDecks(6),
			deck.Shuffle,
		),
		turnPlayerIndex: 0,
	}
}
