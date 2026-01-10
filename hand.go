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

func (h *hands) halt() {
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

func newHand(cards []deck.Card, bet int, isActive bool) *hand {
	return &hand{
		cards:    cards,
		isActive: isActive,
		bet:      bet,
	}
}

func newHands(cards []deck.Card, bet int) *hands {
	first := newHand(cards, bet, true)
	return &hands{
		mode:   normal,
		first:  first,
		second: nil,
		active: first,
	}
}
