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
	dealer     *dealer
	players    [7]*Player
	deck       []deck.Card
	turnPlayer *Player
}

// Start starts the round at the table by dealing everyone two cards.
// After dealing the cards it checks if any of the players has black jack.
func (t *Table) Start() error {
	for range 2 {
		for _, p := range t.players {
			if p == nil {
				continue
			}
			card := t.drawCard()
			p.Hit(card)
		}

		card := t.drawCard()
		t.dealer.Hit(card)
	}

	// TODO end round somehow when no turnPlayer can be determined
	for _, p := range t.players {
		if p != nil && !p.hasBlackJack() {
			t.turnPlayer = p
			return nil
		}
	}

	return nil
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

func (t *Table) drawCard() deck.Card {
	cards, remaining := deck.Draw(1)(t.deck)
	t.deck = remaining
	return cards[0]
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
		turnPlayer: nil,
	}
}
