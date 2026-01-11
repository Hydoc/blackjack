package blackjack

import (
	"reflect"
	"testing"

	"github.com/Hydoc/deck"
)

func TestNewPlayer(t *testing.T) {
	wallet := 200
	p := NewPlayer(
		wallet,
		WithName("Test"),
	)

	if p.Name != "Test" {
		t.Errorf("name should be Test")
	}

	if p.wallet != wallet {
		t.Errorf("wallet should be %d", wallet)
	}
}

func TestPlayer_CanDoubleDown(t *testing.T) {
	tests := []struct {
		name   string
		player *Player
		want   bool
	}{
		{
			name: "can double down",
			player: &Player{
				wallet: 300,
				hands: &hands{
					active: &hand{
						bet: 200,
						cards: []deck.Card{
							{Rank: deck.Seven, Suit: deck.Heart},
							{Rank: deck.Two, Suit: deck.Club},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "not possible due to wallet",
			player: &Player{
				wallet: 199,
				hands: &hands{
					active: &hand{
						bet: 200,
						cards: []deck.Card{
							{Rank: deck.Seven, Suit: deck.Heart},
							{Rank: deck.Two, Suit: deck.Club},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "not possible due to higher than allowed",
			player: &Player{
				wallet: 500,
				hands: &hands{
					active: &hand{
						bet: 200,
						cards: []deck.Card{
							{Rank: deck.Eight, Suit: deck.Heart},
							{Rank: deck.Ace, Suit: deck.Club},
						},
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.player.CanDoubleDown(); got != tt.want {
				t.Errorf("want %#v, got %#v", tt.want, got)
			}
		})
	}
}

func TestPlayer_CanSplit(t *testing.T) {
	tests := []struct {
		name   string
		player *Player
		want   bool
	}{
		{
			name: "can split",
			player: &Player{
				wallet: 200,
				hands: &hands{
					active: &hand{
						bet: 200,
						cards: []deck.Card{
							{Rank: deck.Ten, Suit: deck.Spade},
							{Rank: deck.Ten, Suit: deck.Heart},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "can not split with different ranks",
			player: &Player{
				hands: &hands{
					active: &hand{
						cards: []deck.Card{
							{Rank: deck.Two, Suit: deck.Spade},
							{Rank: deck.Three, Suit: deck.Heart},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "can not split with more than two cards",
			player: &Player{
				hands: &hands{
					active: &hand{
						cards: []deck.Card{
							{Rank: deck.Two, Suit: deck.Spade},
							{Rank: deck.Two, Suit: deck.Heart},
							{Rank: deck.Three, Suit: deck.Heart},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "can not split due to wallet",
			player: &Player{
				wallet: 199,
				hands: &hands{
					active: &hand{
						bet: 200,
						cards: []deck.Card{
							{Rank: deck.Two, Suit: deck.Spade},
							{Rank: deck.Two, Suit: deck.Heart},
						},
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.player.CanSplit(); got != tt.want {
				t.Errorf("want %#v, got %#v", tt.want, got)
			}
		})
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
				firstHand := newHand([]deck.Card{
					{
						Suit: deck.Heart,
						Rank: deck.Ace,
					},
					{
						Suit: deck.Spade,
						Rank: deck.Two,
					},
				}, true)
				return &Player{
					hands: &hands{
						first:  firstHand,
						second: newHand(make([]deck.Card, 0), false),
						active: firstHand,
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
				secondHand := newHand([]deck.Card{
					{
						Suit: deck.Heart,
						Rank: deck.Ace,
					},
					{
						Suit: deck.Spade,
						Rank: deck.Two,
					},
				}, true)
				return &Player{
					hands: &hands{
						first:  newHand(make([]deck.Card, 0), false),
						second: secondHand,
						active: secondHand,
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
		wantFirstHand  *hand
		wantSecondHand *hand
	}{
		{
			name: "with first hand active and normal mode",
			player: &Player{
				hands: &hands{
					mode:   normal,
					first:  &hand{isActive: true},
					second: &hand{isActive: false},
				},
			},
			wantFirstHand:  &hand{isActive: false},
			wantSecondHand: &hand{isActive: false},
		},
		{
			name: "with first hand active and split mode",
			player: &Player{
				hands: &hands{
					mode:   split,
					first:  &hand{isActive: true},
					second: &hand{isActive: false},
				},
			},
			wantFirstHand:  &hand{isActive: false},
			wantSecondHand: &hand{isActive: true},
		},
		{
			name: "with second hand active and split mode",
			player: &Player{
				hands: &hands{
					mode:   split,
					first:  &hand{isActive: false},
					second: &hand{isActive: true},
				},
			},
			wantFirstHand:  &hand{isActive: false},
			wantSecondHand: &hand{isActive: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.player.Stand()

			if !reflect.DeepEqual(tt.wantFirstHand, tt.player.hands.first) {
				t.Errorf("want %#v, got %#v", tt.wantFirstHand, tt.player.hands.first)
			}

			if !reflect.DeepEqual(tt.wantSecondHand, tt.player.hands.second) {
				t.Errorf("want %#v, got %#v", tt.wantSecondHand, tt.player.hands.second)
			}
		})
	}
}
