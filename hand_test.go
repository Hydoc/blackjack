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
	hands := NewHands(cards, 200)

	if hands.Mode != normal {
		t.Errorf("want %#v, got %#v", normal, hands.Mode)
	}

	if hands.second != nil {
		t.Error("want second hand to be nil")
	}

	if !reflect.DeepEqual(hands.first, hands.Active) {
		t.Error("first and active hand should be the same")
	}
}

func TestHands_Hit(t *testing.T) {
	cards := []deck.Card{
		{Rank: deck.Ten, Suit: deck.Spade},
		{Rank: deck.Two, Suit: deck.Heart},
	}
	hands := NewHands(cards, 200)
	cardToHit := deck.Card{Rank: deck.Five, Suit: deck.Club}
	want := append(append([]deck.Card{}, cards...), cardToHit)

	hands.Hit(cardToHit)

	if !reflect.DeepEqual(hands.first, hands.Active) {
		t.Error("first and active hand should be the same")
	}

	if !reflect.DeepEqual(want, hands.first.cards) {
		t.Errorf("want %#v, got %#v", want, hands.first.cards)
	}

	if !reflect.DeepEqual(want, hands.Active.cards) {
		t.Errorf("want %#v, got %#v", want, hands.Active.cards)
	}
}

func TestHands_Halt(t *testing.T) {
	tests := []struct {
		name           string
		hands          *Hands
		wantFirstHand  *Hand
		wantSecondHand *Hand
		wantActiveHand *Hand
	}{
		{
			name: "with first hand active and normal mode",
			hands: &Hands{
				Mode:   normal,
				first:  &Hand{isActive: true},
				second: &Hand{isActive: false},
			},
			wantFirstHand:  &Hand{isActive: false},
			wantSecondHand: &Hand{isActive: false},
			wantActiveHand: &Hand{isActive: false},
		},
		{
			name: "with first hand active and split mode",
			hands: &Hands{
				Mode:   split,
				first:  &Hand{isActive: true},
				second: &Hand{isActive: false},
			},
			wantFirstHand:  &Hand{isActive: false},
			wantSecondHand: &Hand{isActive: true},
			wantActiveHand: &Hand{isActive: true},
		},
		{
			name: "with second hand active and split mode",
			hands: &Hands{
				Mode:   split,
				first:  &Hand{isActive: false},
				second: &Hand{isActive: true},
			},
			wantFirstHand:  &Hand{isActive: false},
			wantSecondHand: &Hand{isActive: false},
			wantActiveHand: &Hand{isActive: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.hands.Halt()

			if !reflect.DeepEqual(tt.wantFirstHand, tt.hands.first) {
				t.Errorf("want %#v, got %#v", tt.wantFirstHand, tt.hands.first)
			}

			if !reflect.DeepEqual(tt.wantSecondHand, tt.hands.second) {
				t.Errorf("want %#v, got %#v", tt.wantSecondHand, tt.hands.second)
			}
		})

	}
}
