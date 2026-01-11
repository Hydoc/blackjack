package blackjack

import (
	"errors"
	"slices"

	"github.com/Hydoc/deck"
)

var (
	ErrSplitNotAllowed = errors.New("can't split")
)

type hands struct {
	mode   Mode
	first  *hand
	second *hand
	active *hand
}

func (h *hands) hit(card deck.Card) {
	h.active.hit(card)
}

func (h *hands) stand() {
	if h.first.isActive {
		h.first.isActive = false
		h.active = nil

		if h.mode == split {
			h.second.isActive = true
			h.active = h.second
		}
	} else if h.second.isActive {
		h.second.isActive = false
		h.active = nil
	}
}

func (h *hands) canSplit() bool {
	return h.active.canSplit()
}

type hand struct {
	cards    []deck.Card
	isActive bool
	bet      int
}

func (h *hand) hit(card deck.Card) {
	h.cards = append(h.cards, card)
}

func (h *hand) split() (*hands, error) {
	if !h.canSplit() {
		return nil, ErrSplitNotAllowed
	}

	return newSplitHands(h.cards[0], h.cards[1], h.bet), nil
}

func (h *hand) canSplit() bool {
	return len(h.cards) == 2 && h.cards[0].Rank == h.cards[1].Rank
}

func (h *hand) sum() int {
	sum := 0

	for _, card := range h.cards {
		if card.Rank == deck.Jack || card.Rank == deck.Queen || card.Rank == deck.King {
			sum += 10
		} else {
			sum += int(card.Rank)
		}
	}

	hasAce := slices.ContainsFunc(h.cards, func(card deck.Card) bool {
		return card.Rank == deck.Ace
	})

	if sum < 21 && hasAce {
		sum += 10
	}

	return sum
}

func newHand(cards []deck.Card, isActive bool, opts ...func(*hand) *hand) *hand {
	h := &hand{
		cards:    cards,
		isActive: isActive,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func newSplitHands(first deck.Card, second deck.Card, previousBet int) *hands {
	f := newHand([]deck.Card{first}, true, withBet(previousBet))
	s := newHand([]deck.Card{second}, false, withBet(previousBet))

	return &hands{
		mode:   split,
		first:  f,
		second: s,
		active: f,
	}
}

func newHands(cards []deck.Card, opts ...func(*hand) *hand) *hands {
	first := newHand(cards, true, opts...)
	return &hands{
		mode:   normal,
		first:  first,
		second: nil,
		active: first,
	}
}

func withBet(bet int) func(*hand) *hand {
	return func(hand *hand) *hand {
		hand.bet = bet
		return hand
	}
}
