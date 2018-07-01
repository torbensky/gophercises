package main

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/torbensky/gophercises/deck"
)

type Game struct {
	Players []string
}

var app *tview.Application

func newGame() {
	// How many players?
	form := tview.NewForm()

	players := 1
	form.
		AddDropDown("How many players?", []string{
			"one",
			"two",
			"three",
			"four",
		}, 0, func(option string, idx int) {
			players = idx + 1
		}).
		AddButton("Next", func() {
			setupPlayers(players)
		}).
		AddButton("Quit", func() {
			app.Stop()
		}).
		SetBorder(true).
		SetTitle("Blackjack - New Game").
		SetTitleAlign(tview.AlignLeft)

	app.SetRoot(form, true)
}

func setupPlayers(num int) {
	form := tview.NewForm()
	players := make([]string, num)
	for i := 0; i < num; i++ {
		players[i] = fmt.Sprintf("player%d", i+1)
		form.AddInputField(fmt.Sprintf("Player%d Name:", i+1), players[i], 20, func(val string, l rune) bool {
			players[i] = val
			return true
		}, nil)
	}
	form.AddButton("Start", func() {
		Play(Game{
			Players: players,
		})
	})
	form.AddButton("Quit", func() {
		app.Stop()
	})

	app.SetRoot(form, true)
}

func drawHand(playerName string, hand []deck.Card) *tview.Table {
	table := tview.NewTable().SetBorders(true)
	table.SetCell(0, 0, tview.NewTableCell(playerName))
	for i, card := range hand {
		table.SetCell(0, i+1, tview.NewTableCell(card.String()))
	}

	return table
}

func main() {
	app = tview.NewApplication()
	newGame()

	if err := app.Run(); err != nil {
		panic(err)
	}

	// What is each player's name?

	// Play the game
}

func Play(game Game) {
	d := deck.New(deck.ShuffleRand)

	var playerCards []deck.Card
	var dealerCards []deck.Card

	dealCards(&d, &playerCards, &dealerCards)
	dealCards(&d, &playerCards, &dealerCards)

	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	flex.AddItem(drawHand("player1", playerCards), 0, 1, true)
	flex.AddItem(drawHand("dealer", dealerCards), 0, 1, false)
	app.SetRoot(flex, true)
}

func dealCards(d *[]deck.Card, hands ...*[]deck.Card) {
	for _, h := range hands {
		*h = append(*h, (*d)[0])
		*d = (*d)[1:]
	}
}
