package blackjack

import (
	"reflect"
	"testing"

	"github.com/Hydoc/deck"
)

func Test_New(t *testing.T) {
	table := New()
	wantDealer := NewPlayer(make([]deck.Card, 0), 0, WithName("Dealer"))
	wantPlayers := [7]*Player{}
	wantDeckSize := 52 * 6

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
