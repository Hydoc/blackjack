package blackjack

import "github.com/Hydoc/deck"

type hands struct {
	mode   Mode
	first  *hand
	second *hand
	active *hand
}

func (h *hands) hit(card deck.Card) {
	h.active.cards = append(h.active.cards, card)
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

type hand struct {
	cards    []deck.Card
	isActive bool
	bet      int
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
