package blackjack

import (
	"reflect"
	"testing"

	"github.com/Hydoc/deck"
)

func TestNewPlayer(t *testing.T) {
	cards := make([]deck.Card, 0)
	bet := 200
	p := NewPlayer(
		cards,
		bet,
		WithName("Test"),
	)

	wantHands := NewHands(cards, bet)

	if p.Name != "Test" {
		t.Errorf("name should be Test")
	}

	if !reflect.DeepEqual(wantHands, p.hands) {
		t.Errorf("want %#v, got %#v", wantHands, p.hands)
	}
}

func TestPlayer_Hit(t *testing.T) {
	tests := []struct {
		name                string
		setup               func() *Player
		cardToHit           deck.Card
		wantFirstHandCards  []deck.Card
		wantSecondHandCards []deck.Card
	}{
		{
			name: "with first hand active",
			setup: func() *Player {
				firstHand := NewHand([]deck.Card{
					{
						Suit: deck.Heart,
						Rank: deck.Ace,
					},
					{
						Suit: deck.Spade,
						Rank: deck.Two,
					},
				}, 200, true)
				return &Player{
					hands: &Hands{
						first:  firstHand,
						second: NewHand(make([]deck.Card, 0), 0, false),
						Active: firstHand,
					},
				}
			},
			cardToHit: deck.Card{
				Suit: deck.Club,
				Rank: deck.Three,
			},
			wantFirstHandCards: []deck.Card{
				{
					Suit: deck.Heart,
					Rank: deck.Ace,
				},
				{
					Suit: deck.Spade,
					Rank: deck.Two,
				},
				{
					Suit: deck.Club,
					Rank: deck.Three,
				},
			},
			wantSecondHandCards: make([]deck.Card, 0),
		},
		{
			name: "with second hand active",
			setup: func() *Player {
				secondHand := NewHand([]deck.Card{
					{
						Suit: deck.Heart,
						Rank: deck.Ace,
					},
					{
						Suit: deck.Spade,
						Rank: deck.Two,
					},
				}, 3, true)
				return &Player{
					hands: &Hands{
						first:  NewHand(make([]deck.Card, 0), 200, false),
						second: secondHand,
						Active: secondHand,
					},
				}
			},
			cardToHit: deck.Card{
				Suit: deck.Club,
				Rank: deck.Three,
			},
			wantFirstHandCards: make([]deck.Card, 0),
			wantSecondHandCards: []deck.Card{
				{
					Suit: deck.Heart,
					Rank: deck.Ace,
				},
				{
					Suit: deck.Spade,
					Rank: deck.Two,
				},
				{
					Suit: deck.Club,
					Rank: deck.Three,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := tt.setup()
			player.Hit(tt.cardToHit)

			if !reflect.DeepEqual(player.hands.first.cards, tt.wantFirstHandCards) {
				t.Errorf("want %#v, got %#v", tt.wantFirstHandCards, player.hands.first.cards)
			}

			if !reflect.DeepEqual(player.hands.second.cards, tt.wantSecondHandCards) {
				t.Errorf("want %#v, got %#v", tt.wantSecondHandCards, player.hands.second.cards)
			}
		})
	}
}

func TestPlayer_Halt(t *testing.T) {
	tests := []struct {
		name           string
		player         *Player
		wantFirstHand  *Hand
		wantSecondHand *Hand
	}{
		{
			name: "with first hand active and normal mode",
			player: &Player{
				hands: &Hands{
					Mode:   normal,
					first:  &Hand{isActive: true},
					second: &Hand{isActive: false},
				},
			},
			wantFirstHand:  &Hand{isActive: false},
			wantSecondHand: &Hand{isActive: false},
		},
		{
			name: "with first hand active and split mode",
			player: &Player{
				hands: &Hands{
					Mode:   split,
					first:  &Hand{isActive: true},
					second: &Hand{isActive: false},
				},
			},
			wantFirstHand:  &Hand{isActive: false},
			wantSecondHand: &Hand{isActive: true},
		},
		{
			name: "with second hand active and split mode",
			player: &Player{
				hands: &Hands{
					Mode:   split,
					first:  &Hand{isActive: false},
					second: &Hand{isActive: true},
				},
			},
			wantFirstHand:  &Hand{isActive: false},
			wantSecondHand: &Hand{isActive: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.player.Halt()

			if !reflect.DeepEqual(tt.wantFirstHand, tt.player.hands.first) {
				t.Errorf("want %#v, got %#v", tt.wantFirstHand, tt.player.hands.first)
			}

			if !reflect.DeepEqual(tt.wantSecondHand, tt.player.hands.second) {
				t.Errorf("want %#v, got %#v", tt.wantSecondHand, tt.player.hands.second)
			}
		})
	}
}
