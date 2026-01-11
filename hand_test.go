package blackjack

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Hydoc/deck"
)

func TestNewHands(t *testing.T) {
	cards := []deck.Card{
		{Rank: deck.Ten, Suit: deck.Spade},
		{Rank: deck.Two, Suit: deck.Heart},
	}
	h := newHands(cards, withBet(200))

	if h.mode != normal {
		t.Errorf("want %#v, got %#v", normal, h.mode)
	}

	if h.second != nil {
		t.Error("want second hand to be nil")
	}

	if !reflect.DeepEqual(h.first, h.active) {
		t.Error("first and active hand should be the same")
	}
}

func TestHands_hit(t *testing.T) {
	cards := []deck.Card{
		{Rank: deck.Ten, Suit: deck.Spade},
		{Rank: deck.Two, Suit: deck.Heart},
	}
	h := newHands(cards)
	cardToHit := deck.Card{Rank: deck.Five, Suit: deck.Club}
	want := append(append([]deck.Card{}, cards...), cardToHit)

	h.hit(cardToHit)

	if !reflect.DeepEqual(h.first, h.active) {
		t.Error("first and active hand should be the same")
	}

	if !reflect.DeepEqual(want, h.first.cards) {
		t.Errorf("want %#v, got %#v", want, h.first.cards)
	}

	if !reflect.DeepEqual(want, h.active.cards) {
		t.Errorf("want %#v, got %#v", want, h.active.cards)
	}
}

func TestHands_stand(t *testing.T) {
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

func TestHands_canSplit(t *testing.T) {
	tests := []struct {
		name  string
		want  bool
		cards []deck.Card
	}{
		{
			name: "can split",
			want: true,
			cards: []deck.Card{
				{Rank: deck.Ten, Suit: deck.Spade},
				{Rank: deck.Ten, Suit: deck.Heart},
			},
		},
		{
			name: "can not split with different ranks",
			want: false,
			cards: []deck.Card{
				{Rank: deck.Two, Suit: deck.Spade},
				{Rank: deck.Three, Suit: deck.Heart},
			},
		},
		{
			name: "can not split with more than two cards",
			want: false,
			cards: []deck.Card{
				{Rank: deck.Two, Suit: deck.Spade},
				{Rank: deck.Two, Suit: deck.Heart},
				{Rank: deck.Three, Suit: deck.Heart},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hands{
				active: newHand(tt.cards, true),
			}

			if got := h.canSplit(); got != tt.want {
				t.Errorf("want %#v, got %#v", tt.want, got)
			}
		})
	}
}

func TestHand_split(t *testing.T) {
	tests := []struct {
		name      string
		hand      *hand
		wantErr   error
		wantHands *hands
	}{
		{
			name: "split correctly",
			hand: &hand{
				cards: []deck.Card{
					{Rank: deck.Ten, Suit: deck.Spade},
					{Rank: deck.Ten, Suit: deck.Heart},
				},
				bet: 100,
			},
			wantHands: newSplitHands(
				deck.Card{Rank: deck.Ten, Suit: deck.Spade},
				deck.Card{Rank: deck.Ten, Suit: deck.Heart},
				100,
			),
		},
		{
			name: "splitting not allowed",
			hand: &hand{
				cards: []deck.Card{
					{Rank: deck.Ten, Suit: deck.Spade},
					{Rank: deck.Nine, Suit: deck.Heart},
				},
				bet: 100,
			},
			wantErr: ErrSplitNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := tt.hand.split()

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want %#v, got %#v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(h, tt.wantHands) {
				t.Errorf("want %#v, got %#v", tt.wantHands, h)
			}
		})
	}
}

func TestHand_sum(t *testing.T) {
	tests := []struct {
		name  string
		cards []deck.Card
		want  int
	}{
		{
			name: "queen, jack and ace",
			cards: []deck.Card{
				{Rank: deck.Queen, Suit: deck.Spade},
				{Rank: deck.Jack, Suit: deck.Heart},
				{Rank: deck.Ace, Suit: deck.Club},
			},
			want: 21,
		},
		{
			name: "two and three",
			cards: []deck.Card{
				{Rank: deck.Two, Suit: deck.Spade},
				{Rank: deck.Three, Suit: deck.Heart},
			},
			want: 5,
		},
		{
			name: "single ace",
			cards: []deck.Card{
				{Rank: deck.Ace, Suit: deck.Spade},
			},
			want: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hand{
				cards: tt.cards,
			}

			if got := h.sum(); got != tt.want {
				t.Errorf("want %#v, got %#v", tt.want, got)
			}
		})
	}
}

func TestHand_hasBlackJack(t *testing.T) {
	tests := []struct {
		name  string
		cards []deck.Card
		want  bool
	}{
		{
			name: "has black jack",
			cards: []deck.Card{
				{Rank: deck.Ten, Suit: deck.Spade},
				{Rank: deck.Ace, Suit: deck.Heart},
			},
			want: true,
		},
		{
			name: "has not black jack",
			cards: []deck.Card{
				{Rank: deck.Ten, Suit: deck.Spade},
				{Rank: deck.Ten, Suit: deck.Heart},
			},
			want: false,
		},
		{
			name: "21 but with three cards",
			cards: []deck.Card{
				{Rank: deck.Ten, Suit: deck.Spade},
				{Rank: deck.Nine, Suit: deck.Heart},
				{Rank: deck.Two, Suit: deck.Heart},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hand{
				cards: tt.cards,
			}
			if got := h.hasBlackJack(); got != tt.want {
				t.Errorf("want %#v, got %#v", tt.want, got)
			}
		})
	}
}
