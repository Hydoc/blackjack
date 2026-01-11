package blackjack

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Hydoc/deck"
)

func Test_New(t *testing.T) {
	table := New()
	wantDealer := newDealer()
	wantPlayers := [7]*Player{}
	wantDeckSize := 312

	if len(table.deck) != wantDeckSize {
		t.Errorf("want deck size %d, got %d", wantDeckSize, len(table.deck))
	}

	if !reflect.DeepEqual(wantDealer, table.dealer) {
		t.Errorf("want dealer %v, got %v", wantDealer, table.dealer)
	}

	if !reflect.DeepEqual(wantPlayers, table.players) {
		t.Errorf("want players %v, got %v", wantPlayers, table.players)
	}
}

func TestTable_Join(t *testing.T) {
	tests := []struct {
		name           string
		playerToJoin   *Player
		playersAtTable [7]*Player
		wantIndex      int
		wantErr        error
	}{
		{
			name:           "with empty slice of players at table",
			playerToJoin:   NewPlayer(0, WithName("Player1")),
			playersAtTable: [7]*Player{},
			wantIndex:      0,
		},
		{
			name:         "with players at table",
			playerToJoin: NewPlayer(0, WithName("Player4")),
			playersAtTable: [7]*Player{
				NewPlayer(0, WithName("Player1")),
				NewPlayer(0, WithName("Player2")),
				NewPlayer(0, WithName("Player3")),
			},
			wantIndex: 3,
		},
		{
			name:         "error when table is full",
			playerToJoin: NewPlayer(0, WithName("Player8")),
			playersAtTable: [7]*Player{
				NewPlayer(0, WithName("Player1")),
				NewPlayer(0, WithName("Player2")),
				NewPlayer(0, WithName("Player3")),
				NewPlayer(0, WithName("Player4")),
				NewPlayer(0, WithName("Player5")),
				NewPlayer(0, WithName("Player6")),
				NewPlayer(0, WithName("Player7")),
			},
			wantErr: ErrTableFull,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{
				players: tt.playersAtTable,
			}

			err := table.Join(tt.playerToJoin)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want %#v, got %#v", tt.wantErr, err)
			}

			if tt.wantErr == nil {
				if !reflect.DeepEqual(table.players[tt.wantIndex], tt.playerToJoin) {
					t.Errorf("want %#v, got %#v", tt.playerToJoin, table.players[tt.wantIndex])
				}
			}
		})
	}
}

func TestTable_Start(t *testing.T) {
	playerOne := NewPlayer(0, WithName("Player1"))
	wantPlayerOneCards := []deck.Card{
		{Rank: deck.King, Suit: deck.Heart},
		{Rank: deck.Ten, Suit: deck.Heart},
	}
	playerTwo := NewPlayer(0, WithName("Player2"))
	wantPlayerTwoCards := []deck.Card{
		{Rank: deck.Queen, Suit: deck.Heart},
		{Rank: deck.Nine, Suit: deck.Heart},
	}

	wantDealerCards := []deck.Card{
		{Rank: deck.Jack, Suit: deck.Heart},
		{Rank: deck.Eight, Suit: deck.Heart},
	}

	table := &Table{
		dealer:  newDealer(),
		players: [7]*Player{},
		deck:    deck.New(),
	}
	err := table.Join(playerOne)
	if err != nil {
		t.Errorf("want nil, got %v", err)
	}

	err = table.Join(playerTwo)
	if err != nil {
		t.Errorf("want nil, got %v", err)
	}

	table.Start()

	if len(table.deck) != 46 {
		t.Errorf("want %d, got %d", 46, len(table.deck))
	}

	if !reflect.DeepEqual(table.turnPlayer, playerOne) {
		t.Errorf("want %#v, got %#v", playerOne, table.turnPlayer)
	}

	if !reflect.DeepEqual(table.players[0].hands.first.cards, wantPlayerOneCards) {
		t.Errorf("want %#v, got %#v", wantPlayerOneCards, table.players[0].hands.first.cards)
	}

	if !reflect.DeepEqual(table.players[1].hands.first.cards, wantPlayerTwoCards) {
		t.Errorf("want %#v, got %#v", wantPlayerOneCards, table.players[1].hands.first.cards)
	}

	if !reflect.DeepEqual(table.dealer.hand.cards, wantDealerCards) {
		t.Errorf("want %#v, got %#v", wantDealerCards, table.dealer.hand.cards)
	}
}
