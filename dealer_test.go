package blackjack

import (
	"reflect"
	"testing"

	"github.com/Hydoc/deck"
)

func TestDealer_HitUntil17(t *testing.T) {
	tests := []struct {
		name          string
		cards         []deck.Card
		wantSum       int
		wantCards     []deck.Card
		wantRemaining []deck.Card
	}{
		{
			name: "six -> seven -> queen",
			cards: []deck.Card{
				{Rank: deck.Jack, Suit: deck.Heart},
				{Rank: deck.Queen, Suit: deck.Spade},
				{Rank: deck.Seven, Suit: deck.Spade},
				{Rank: deck.Six, Suit: deck.Club},
			},
			wantSum: 23,
			wantCards: []deck.Card{
				{Rank: deck.Six, Suit: deck.Club},
				{Rank: deck.Seven, Suit: deck.Spade},
				{Rank: deck.Queen, Suit: deck.Spade},
			},
			wantRemaining: []deck.Card{
				{Rank: deck.Jack, Suit: deck.Heart},
			},
		},
		{
			name: "ace -> six",
			cards: []deck.Card{
				{Rank: deck.Queen, Suit: deck.Spade},
				{Rank: deck.Six, Suit: deck.Club},
				{Rank: deck.Ace, Suit: deck.Spade},
			},
			wantSum: 17,
			wantCards: []deck.Card{
				{Rank: deck.Ace, Suit: deck.Spade},
				{Rank: deck.Six, Suit: deck.Club},
			},
			wantRemaining: []deck.Card{
				{Rank: deck.Queen, Suit: deck.Spade},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newDealer()
			remaining := d.HitUntil17(tt.cards)

			if !reflect.DeepEqual(d.hand.cards, tt.wantCards) {
				t.Errorf("want %#v, got %#v", tt.wantCards, d.hand.cards)
			}

			if !reflect.DeepEqual(remaining, tt.wantRemaining) {
				t.Errorf("want %#v, got %#v", tt.wantRemaining, remaining)
			}

			if tt.wantSum != d.hand.sum() {
				t.Errorf("want %d, got %d", tt.wantSum, d.hand.sum())
			}
		})
	}
}
