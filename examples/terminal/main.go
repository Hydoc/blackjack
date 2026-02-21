package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Hydoc/deck"

	"github.com/Hydoc/blackjack"
)

var cardTemplate = ` _______
|%s		|
|		|
|	%s	|
|		|
|_______|`

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	table := blackjack.New()
	playerOne := blackjack.NewPlayer(500, blackjack.WithName("One"))
	playerTwo := blackjack.NewPlayer(500, blackjack.WithName("Two"))

	err := table.Join(playerOne)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	err = table.Join(playerTwo)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	table.Start()

	for !table.IsDone() {
		state := table.State()

		printDealer(state.Dealer, true)

		for _, p := range state.Players {

			fmt.Println(p.Name)
		}
	}
}

func printCard(card deck.Card) {
	fmt.Printf(cardTemplate, card.Rank, card.Suit)
}

func printDealer(d *blackjack.Dealer, onlyFirstCard bool) {
	if onlyFirstCard {
		printCard(d.Cards()[0])
		return
	}
}

func printPlayer(p *blackjack.Player) {

}
