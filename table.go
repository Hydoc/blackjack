package blackjack

import (
	"errors"
	"sync"

	"github.com/Hydoc/deck"
)

type GameState = int

const (
	inProgress GameState = iota
	done
)

var (
	ErrTableFull = errors.New("table is full")
)

// Table represents a blackjack table. It holds everything relevant for the game.
// The game state, players, deck, turn player and dealer
type Table struct {
	mu sync.Mutex

	gameState  GameState
	dealer     *dealer
	players    [7]*Player
	deck       []deck.Card
	turnPlayer *Player
}

// Start starts the round at the table by dealing everyone two cards.
// After dealing the cards it checks if any of the players has black jack and sets the gameState
// which can be checked using either InProgress or IsDone.
func (t *Table) Start() {
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

	for _, p := range t.players {
		if p != nil && !p.hasBlackJack() {
			t.turnPlayer = p
			t.gameState = inProgress
			return
		}
	}

	t.turnPlayer = nil
	t.gameState = done
}

// InProgress returns a bool whether the gameState is inProgress.
func (t *Table) InProgress() bool {
	return t.gameState == inProgress
}

// IsDone returns a bool whether the gameState is done.
func (t *Table) IsDone() bool {
	return t.gameState == done
}

func (t *Table) Handle() error {
	return nil
}

// Join adds a player to the next nil value in the players slice.
// It returns ErrTableFull when there is no space left.
func (t *Table) Join(p *Player) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i := range t.players {
		if t.players[i] == nil {
			t.players[i] = p
			return nil
		}
	}
	return ErrTableFull
}

// Leave removes a player from the table if it was found. It does nothing otherwise
func (t *Table) Leave(p *Player) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i := range t.players {
		if t.players[i] == p {
			t.players[i] = nil
			return
		}
	}
}

func (t *Table) drawCard() deck.Card {
	cards, remaining := deck.Draw(1)(t.deck)
	t.deck = remaining
	return cards[0]
}

// New create a pointer to Table with the default configuration like 6 shuffled decks and a maximum of 7 players allowed.
// Dealer must stand on soft 17.
// No peek.
// Double Down only allowed on 9 to 11.
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
