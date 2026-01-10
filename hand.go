package blackjack

import "github.com/Hydoc/deck"

type Hands struct {
	Mode Mode

	first  *Hand
	second *Hand
	Active *Hand
}

func (h *Hands) Hit(card deck.Card) {
	h.Active.cards = append(h.Active.cards, card)
}

func (h *Hands) Halt() {
	if h.first.isActive {
		h.first.isActive = false
		h.Active = nil

		if h.Mode == split {
			h.second.isActive = true
			h.Active = h.second
		}
	} else if h.second.isActive {
		h.second.isActive = false
		h.Active = nil
	}
}

type Hand struct {
	cards    []deck.Card
	isActive bool
	bet      int
}

func NewHand(cards []deck.Card, bet int, isActive bool) *Hand {
	return &Hand{
		cards:    cards,
		isActive: isActive,
		bet:      bet,
	}
}

// NewHands creates a pointer to the Hands struct.
// It initializes the first one and nils the second one.
// The second hand is only relevant when splitting.
func NewHands(cards []deck.Card, bet int) *Hands {
	first := NewHand(cards, bet, true)
	return &Hands{
		Mode:   normal,
		first:  first,
		second: nil,
		Active: first,
	}
}
