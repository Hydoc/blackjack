package blackjack

import "github.com/Hydoc/deck"

type dealer struct {
	hand *hand
}

func (d *dealer) Hit(card deck.Card) {
	d.hand.hit(card)
}

func (d *dealer) HitUntil17(cards []deck.Card) []deck.Card {
	remaining := cards
	for d.hand.sum() < 16 {
		c, leftover := deck.Draw(1)(remaining)
		d.Hit(c[0])
		remaining = leftover
	}
	return remaining
}

func newDealer() *dealer {
	return &dealer{
		hand: newHand(make([]deck.Card, 0), true),
	}
}
