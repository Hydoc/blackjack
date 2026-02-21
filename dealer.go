package blackjack

import "github.com/Hydoc/deck"

type Dealer struct {
	hand *hand
}

func (d *Dealer) Cards() []deck.Card {
	return d.hand.cards
}

func (d *Dealer) hit(card deck.Card) {
	d.hand.hit(card)
}

func (d *Dealer) HitUntil17(cards []deck.Card) []deck.Card {
	remaining := cards
	for d.hand.sum() < 16 {
		c, leftover := deck.Draw(1)(remaining)
		d.hit(c[0])
		remaining = leftover
	}
	return remaining
}

func newDealer() *Dealer {
	return &Dealer{
		hand: newHand(make([]deck.Card, 0), true),
	}
}
