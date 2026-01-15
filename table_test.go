package blackjack

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Hydoc/deck"
)

func Test_New(t *testing.T) {
	table := New()
	wantDealer := newDealer()
	wantPlayers := [7]*Player{}
	wantDeckSize := 312

	if len(table.deck) != wantDeckSize {
		t.Errorf("want deck size %d, got %d", wantDeckSize, len(table.deck))
	}

	if !reflect.DeepEqual(wantDealer, table.dealer) {
		t.Errorf("want dealer %v, got %v", wantDealer, table.dealer)
	}

	if !reflect.DeepEqual(wantPlayers, table.players) {
		t.Errorf("want players %v, got %v", wantPlayers, table.players)
	}
}

func TestTable_Join(t *testing.T) {
	tests := []struct {
		name           string
		playerToJoin   *Player
		playersAtTable [7]*Player
		wantIndex      int
		wantErr        error
	}{
		{
			name:           "with empty slice of players at table",
			playerToJoin:   NewPlayer(0, WithName("Player1")),
			playersAtTable: [7]*Player{},
			wantIndex:      0,
		},
		{
			name:         "with players at table",
			playerToJoin: NewPlayer(0, WithName("Player4")),
			playersAtTable: [7]*Player{
				NewPlayer(0, WithName("Player1")),
				NewPlayer(0, WithName("Player2")),
				NewPlayer(0, WithName("Player3")),
			},
			wantIndex: 3,
		},
		{
			name:         "error when table is full",
			playerToJoin: NewPlayer(0, WithName("Player8")),
			playersAtTable: [7]*Player{
				NewPlayer(0, WithName("Player1")),
				NewPlayer(0, WithName("Player2")),
				NewPlayer(0, WithName("Player3")),
				NewPlayer(0, WithName("Player4")),
				NewPlayer(0, WithName("Player5")),
				NewPlayer(0, WithName("Player6")),
				NewPlayer(0, WithName("Player7")),
			},
			wantErr: ErrTableFull,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{
				players: tt.playersAtTable,
			}

			err := table.Join(tt.playerToJoin)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want %#v, got %#v", tt.wantErr, err)
			}

			if tt.wantErr == nil {
				if !reflect.DeepEqual(table.players[tt.wantIndex], tt.playerToJoin) {
					t.Errorf("want %#v, got %#v", tt.playerToJoin, table.players[tt.wantIndex])
				}
			}
		})
	}
}

func TestTable_Leave(t *testing.T) {
	tests := []struct {
		name  string
		setup func() (*Table, *Player, [7]*Player)
	}{
		{
			name: "leave correctly",
			setup: func() (*Table, *Player, [7]*Player) {
				firstPlayer := &Player{}
				playerToLeave := &Player{}
				thirdPlayer := &Player{}

				wantPlayers := [7]*Player{
					firstPlayer,
					nil,
					thirdPlayer,
				}

				table := New()
				table.Join(firstPlayer)
				table.Join(playerToLeave)
				table.Join(thirdPlayer)

				return table, playerToLeave, wantPlayers
			},
		},
		{
			name: "do nothing for invalid player",
			setup: func() (*Table, *Player, [7]*Player) {
				firstPlayer := &Player{}
				secondPlayer := &Player{}

				wantPlayers := [7]*Player{
					firstPlayer,
					secondPlayer,
				}

				table := New()
				table.Join(firstPlayer)
				table.Join(secondPlayer)

				return table, &Player{}, wantPlayers
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table, playerToLeave, wantPlayers := tt.setup()

			table.Leave(playerToLeave)

			if got := table.players; !reflect.DeepEqual(got, wantPlayers) {
				t.Errorf("want %#v, got %#v", wantPlayers, got)
			}
		})
	}
}

func TestTable_Start(t *testing.T) {
	t.Run("join two players and start", func(t *testing.T) {
		playerOne := NewPlayer(0, WithName("Player1"))
		wantPlayerOneCards := []deck.Card{
			{Rank: deck.King, Suit: deck.Heart},
			{Rank: deck.Ten, Suit: deck.Heart},
		}
		playerTwo := NewPlayer(0, WithName("Player2"))
		wantPlayerTwoCards := []deck.Card{
			{Rank: deck.Queen, Suit: deck.Heart},
			{Rank: deck.Nine, Suit: deck.Heart},
		}

		wantDealerCards := []deck.Card{
			{Rank: deck.Jack, Suit: deck.Heart},
			{Rank: deck.Eight, Suit: deck.Heart},
		}

		table := &Table{
			dealer:  newDealer(),
			players: [7]*Player{},
			deck:    deck.New(),
		}
		err := table.Join(playerOne)
		if err != nil {
			t.Errorf("want nil, got %v", err)
		}

		err = table.Join(playerTwo)
		if err != nil {
			t.Errorf("want nil, got %v", err)
		}

		table.Start()

		if len(table.deck) != 46 {
			t.Errorf("want %d, got %d", 46, len(table.deck))
		}

		if !table.InProgress() {
			t.Errorf("table have the gameState inProgress")
		}

		if !reflect.DeepEqual(table.turnPlayer, playerOne) {
			t.Errorf("want %#v, got %#v", playerOne, table.turnPlayer)
		}

		if !reflect.DeepEqual(table.players[0].hands.first.cards, wantPlayerOneCards) {
			t.Errorf("want %#v, got %#v", wantPlayerOneCards, table.players[0].hands.first.cards)
		}

		if !reflect.DeepEqual(table.players[1].hands.first.cards, wantPlayerTwoCards) {
			t.Errorf("want %#v, got %#v", wantPlayerOneCards, table.players[1].hands.first.cards)
		}

		if !reflect.DeepEqual(table.dealer.hand.cards, wantDealerCards) {
			t.Errorf("want %#v, got %#v", wantDealerCards, table.dealer.hand.cards)
		}
	})

	t.Run("join two players and start but the first has blackjack after dealing", func(t *testing.T) {
		playerOne := NewPlayer(0, WithName("Player1"))
		wantPlayerOneCards := []deck.Card{
			{Rank: deck.Ace, Suit: deck.Heart},
			{Rank: deck.Ten, Suit: deck.Heart},
		}
		playerTwo := NewPlayer(0, WithName("Player2"))
		wantPlayerTwoCards := []deck.Card{
			{Rank: deck.Queen, Suit: deck.Heart},
			{Rank: deck.Nine, Suit: deck.Heart},
		}

		wantDealerCards := []deck.Card{
			{Rank: deck.Jack, Suit: deck.Heart},
			{Rank: deck.Eight, Suit: deck.Heart},
		}

		table := &Table{
			dealer:  newDealer(),
			players: [7]*Player{},
			deck: []deck.Card{
				{Rank: deck.Eight, Suit: deck.Heart},
				{Rank: deck.Nine, Suit: deck.Heart},
				{Rank: deck.Ten, Suit: deck.Heart},
				{Rank: deck.Jack, Suit: deck.Heart},
				{Rank: deck.Queen, Suit: deck.Heart},
				{Rank: deck.Ace, Suit: deck.Heart},
			},
		}
		err := table.Join(playerOne)
		if err != nil {
			t.Errorf("want nil, got %v", err)
		}

		err = table.Join(playerTwo)
		if err != nil {
			t.Errorf("want nil, got %v", err)
		}

		table.Start()

		if !table.InProgress() {
			t.Errorf("table have the gameState inProgress")
		}

		if !reflect.DeepEqual(table.turnPlayer, playerTwo) {
			t.Errorf("want %#v, got %#v", playerTwo, table.turnPlayer)
		}

		if !reflect.DeepEqual(table.players[0].hands.first.cards, wantPlayerOneCards) {
			t.Errorf("want %#v, got %#v", wantPlayerOneCards, table.players[0].hands.first.cards)
		}

		if !reflect.DeepEqual(table.players[1].hands.first.cards, wantPlayerTwoCards) {
			t.Errorf("want %#v, got %#v", wantPlayerOneCards, table.players[1].hands.first.cards)
		}

		if !reflect.DeepEqual(table.dealer.hand.cards, wantDealerCards) {
			t.Errorf("want %#v, got %#v", wantDealerCards, table.dealer.hand.cards)
		}
	})

	t.Run("join one player and he has blackjack after dealing", func(t *testing.T) {
		playerOne := NewPlayer(0, WithName("Player1"))
		wantPlayerOneCards := []deck.Card{
			{Rank: deck.Ace, Suit: deck.Heart},
			{Rank: deck.Ten, Suit: deck.Heart},
		}

		wantDealerCards := []deck.Card{
			{Rank: deck.Jack, Suit: deck.Heart},
			{Rank: deck.Eight, Suit: deck.Heart},
		}

		table := &Table{
			dealer:  newDealer(),
			players: [7]*Player{},
			deck: []deck.Card{
				{Rank: deck.Eight, Suit: deck.Heart},
				{Rank: deck.Ten, Suit: deck.Heart},
				{Rank: deck.Jack, Suit: deck.Heart},
				{Rank: deck.Ace, Suit: deck.Heart},
			},
		}
		err := table.Join(playerOne)
		if err != nil {
			t.Errorf("want nil, got %v", err)
		}

		table.Start()

		if !table.IsDone() {
			t.Errorf("table have the gameState done")
		}

		if table.turnPlayer != nil {
			t.Errorf("want %#v, got %#v", nil, table.turnPlayer)
		}

		if !reflect.DeepEqual(table.players[0].hands.first.cards, wantPlayerOneCards) {
			t.Errorf("want %#v, got %#v", wantPlayerOneCards, table.players[0].hands.first.cards)
		}

		if !reflect.DeepEqual(table.dealer.hand.cards, wantDealerCards) {
			t.Errorf("want %#v, got %#v", wantDealerCards, table.dealer.hand.cards)
		}
	})
}

func TestTable_nextPlayer(t *testing.T) {
	tests := []struct {
		name  string
		setup func() (*Table, *Player)
	}{
		{
			name: "correct player",
			setup: func() (*Table, *Player) {
				firstPlayer := &Player{}
				secondPlayer := &Player{}
				thirdPlayer := &Player{}

				table := &Table{
					players: [7]*Player{
						firstPlayer,
						secondPlayer,
						thirdPlayer,
					},
					turnPlayer: firstPlayer,
				}
				return table, secondPlayer
			},
		},
		{
			name: "no next player for only one",
			setup: func() (*Table, *Player) {
				firstPlayer := &Player{}

				table := &Table{
					players: [7]*Player{
						firstPlayer,
					},
					turnPlayer: firstPlayer,
				}
				return table, nil
			},
		},
		{
			name: "no next player for the last one",
			setup: func() (*Table, *Player) {
				firstPlayer := &Player{}
				secondPlayer := &Player{}
				thirdPlayer := &Player{}

				table := &Table{
					players: [7]*Player{
						firstPlayer,
						secondPlayer,
						nil,
						nil,
						nil,
						nil,
						thirdPlayer,
					},
					turnPlayer: thirdPlayer,
				}
				return table, nil
			},
		},
		{
			name: "correct player for nil values in between",
			setup: func() (*Table, *Player) {
				firstPlayer := &Player{}
				secondPlayer := &Player{}
				thirdPlayer := &Player{}

				table := &Table{
					players: [7]*Player{
						firstPlayer,
						nil,
						secondPlayer,
						nil,
						nil,
						thirdPlayer,
					},
					turnPlayer: secondPlayer,
				}
				return table, thirdPlayer
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table, wantPlayer := tt.setup()

			p := table.nextPlayer()

			if !reflect.DeepEqual(p, wantPlayer) {
				t.Errorf("want %#v, got %#v", wantPlayer, p)
			}
		})
	}
}

func TestTable_Hit(t *testing.T) {
	t.Run("return ErrNoTurnPlayer when turnPlayer = nil", func(t *testing.T) {
		table := &Table{}

		err := table.Hit()
		if !errors.Is(err, ErrNoTurnPlayer) {
			t.Errorf("want %#v, got %#v", ErrNoTurnPlayer, err)
		}
	})

	t.Run("hit normally", func(t *testing.T) {
		player := NewPlayer(200)
		table := &Table{
			turnPlayer: player,
			deck:       deck.New(),
		}

		err := table.Hit()
		if err != nil {
			t.Errorf("should not throw")
		}

		if len(player.hands.active.cards) != 1 {
			t.Errorf("player should have a card")
		}
	})

	t.Run("stand and end after player busted for one player", func(t *testing.T) {
		cards := deck.New()
		player := NewPlayer(200)

		table := &Table{
			turnPlayer: player,
			deck:       cards,
		}

		for range 3 {
			err := table.Hit()
			if err != nil {
				t.Errorf("should not throw")
			}
		}

		if player.hands.active != nil {
			t.Errorf("player should not have an active hand")
		}

		if table.gameState != done {
			t.Errorf("game state should be done")
		}
	})

	t.Run("stand and set next player after player busted with two players at the table", func(t *testing.T) {
		cards := deck.New()
		playerOne := NewPlayer(200)
		playerTwo := NewPlayer(200)

		table := &Table{
			turnPlayer: playerOne,
			players: [7]*Player{
				playerOne,
				playerTwo,
			},
			deck: cards,
		}

		for range 3 {
			err := table.Hit()
			if err != nil {
				t.Errorf("should not throw")
			}
		}

		if playerOne.hands.active != nil {
			t.Errorf("playerOne should not have an active hand")
		}

		if table.turnPlayer != playerTwo {
			t.Errorf("playerTwo should be active")
		}

		if table.gameState != inProgress {
			t.Errorf("game state should be inProgress")
		}
	})
}

func TestTable_Stand(t *testing.T) {
	t.Run("stand for one player and end", func(t *testing.T) {
		table := &Table{
			dealer: newDealer(),
			deck:   deck.New(),
		}

		player := NewPlayer(200)
		err := table.Join(player)

		if err != nil {
			t.Errorf("wanted nil err")
		}

		table.Start()

		err = table.Stand()

		if err != nil {
			t.Errorf("want nil err")
		}

		if !table.IsDone() {
			t.Errorf("table should be done")
		}
	})

	t.Run("stand one player in a two player game", func(t *testing.T) {
		table := &Table{
			dealer: newDealer(),
			deck:   deck.New(),
		}

		playerOne := NewPlayer(200, WithName("One"))
		playerTwo := NewPlayer(400, WithName("Two"))

		err := table.Join(playerOne)
		if err != nil {
			t.Errorf("wanted nil err")
		}

		err = table.Join(playerTwo)
		if err != nil {
			t.Errorf("wanted nil err")
		}

		table.Start()

		err = table.Stand()

		if err != nil {
			t.Errorf("want nil err")
		}

		if !table.InProgress() {
			t.Errorf("table should be in progress")
		}

		if table.turnPlayer != playerTwo {
			t.Errorf("wanted playerTwo to be the turnPlayer")
		}
	})

	t.Run("err when no turnPlayer", func(t *testing.T) {
		table := &Table{}

		err := table.Stand()

		if !errors.Is(err, ErrNoTurnPlayer) {
			t.Errorf("want %#v, got %#v", ErrNoTurnPlayer, err)
		}
	})
}
