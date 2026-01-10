package blackjack

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Hydoc/deck"
)

func Test_New(t *testing.T) {
	table := New()
	wantDealer := NewPlayer(make([]deck.Card, 0), 0, WithName("Dealer"))
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
			playerToJoin:   NewPlayer(make([]deck.Card, 0), 0, WithName("Player1")),
			playersAtTable: [7]*Player{},
			wantIndex:      0,
		},
		{
			name:         "with players at table",
			playerToJoin: NewPlayer(make([]deck.Card, 0), 0, WithName("Player4")),
			playersAtTable: [7]*Player{
				NewPlayer(make([]deck.Card, 0), 0, WithName("Player1")),
				NewPlayer(make([]deck.Card, 0), 0, WithName("Player2")),
				NewPlayer(make([]deck.Card, 0), 0, WithName("Player3")),
			},
			wantIndex: 3,
		},
		{
			name:         "error when table is full",
			playerToJoin: NewPlayer(make([]deck.Card, 0), 0, WithName("Player8")),
			playersAtTable: [7]*Player{
				NewPlayer(make([]deck.Card, 0), 0, WithName("Player1")),
				NewPlayer(make([]deck.Card, 0), 0, WithName("Player2")),
				NewPlayer(make([]deck.Card, 0), 0, WithName("Player3")),
				NewPlayer(make([]deck.Card, 0), 0, WithName("Player4")),
				NewPlayer(make([]deck.Card, 0), 0, WithName("Player5")),
				NewPlayer(make([]deck.Card, 0), 0, WithName("Player6")),
				NewPlayer(make([]deck.Card, 0), 0, WithName("Player7")),
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
