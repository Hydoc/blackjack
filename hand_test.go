package blackjack

import (
	"reflect"
	"testing"

	"github.com/Hydoc/deck"
)

func TestNewHands(t *testing.T) {
	cards := []deck.Card{
		{Rank: deck.Ten, Suit: deck.Spade},
		{Rank: deck.Two, Suit: deck.Heart},
	}
	hands := newHands(cards, 200)

	if hands.mode != normal {
		t.Errorf("want %#v, got %#v", normal, hands.mode)
	}

	if hands.second != nil {
		t.Error("want second hand to be nil")
	}

	if !reflect.DeepEqual(hands.first, hands.active) {
		t.Error("first and active hand should be the same")
	}
}

func TestHands_Hit(t *testing.T) {
	cards := []deck.Card{
		{Rank: deck.Ten, Suit: deck.Spade},
		{Rank: deck.Two, Suit: deck.Heart},
	}
	hands := newHands(cards, 200)
	cardToHit := deck.Card{Rank: deck.Five, Suit: deck.Club}
	want := append(append([]deck.Card{}, cards...), cardToHit)

	hands.hit(cardToHit)

	if !reflect.DeepEqual(hands.first, hands.active) {
		t.Error("first and active hand should be the same")
	}

	if !reflect.DeepEqual(want, hands.first.cards) {
		t.Errorf("want %#v, got %#v", want, hands.first.cards)
	}

	if !reflect.DeepEqual(want, hands.active.cards) {
		t.Errorf("want %#v, got %#v", want, hands.active.cards)
	}
}

func TestHands_Halt(t *testing.T) {
	tests := []struct {
		name           string
		hands          *hands
		wantFirstHand  *hand
		wantSecondHand *hand
		wantActiveHand *hand
	}{
		{
			name: "with first hand active and normal mode",
			hands: &hands{
				mode:   normal,
				first:  &hand{isActive: true},
				second: &hand{isActive: false},
			},
			wantFirstHand:  &hand{isActive: false},
			wantSecondHand: &hand{isActive: false},
			wantActiveHand: &hand{isActive: false},
		},
		{
			name: "with first hand active and split mode",
			hands: &hands{
				mode:   split,
				first:  &hand{isActive: true},
				second: &hand{isActive: false},
			},
			wantFirstHand:  &hand{isActive: false},
			wantSecondHand: &hand{isActive: true},
			wantActiveHand: &hand{isActive: true},
		},
		{
			name: "with second hand active and split mode",
			hands: &hands{
				mode:   split,
				first:  &hand{isActive: false},
				second: &hand{isActive: true},
			},
			wantFirstHand:  &hand{isActive: false},
			wantSecondHand: &hand{isActive: false},
			wantActiveHand: &hand{isActive: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.hands.stand()

			if !reflect.DeepEqual(tt.wantFirstHand, tt.hands.first) {
				t.Errorf("want %#v, got %#v", tt.wantFirstHand, tt.hands.first)
			}

			if !reflect.DeepEqual(tt.wantSecondHand, tt.hands.second) {
				t.Errorf("want %#v, got %#v", tt.wantSecondHand, tt.hands.second)
			}
		})

	}
}
