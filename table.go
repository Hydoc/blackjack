package blackjack

import "github.com/Hydoc/deck"

// Table represents a blackjack table.
// A maximum of seven players are allowed at the table.
type Table struct {
	dealer  *Player
	players [7]*Player
	deck    []deck.Card
}

// New create a pointer to Table with the default configuration like 6 shuffled decks.
func New() *Table {
	return &Table{
		players: [7]*Player{},
		deck: deck.New(
			deck.WithDecks(6),
			deck.Shuffle,
		),
	}
}
