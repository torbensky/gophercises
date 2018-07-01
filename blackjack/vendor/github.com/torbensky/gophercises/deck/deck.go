//go:generate stringer -type=Suit,Rank

package deck

import (
	"fmt"
	"math/rand"
	"sort"
)

type Rank uint8

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

type Suit uint8

const (
	Spades Suit = iota
	Diamonds
	Clubs
	Hearts
	Joker
)

// Card represents a playing card with a suit and rank
type Card struct {
	Suit
	Rank
}

type Option func([]Card) []Card

func New(opts ...Option) []Card {
	var d []Card
	for rank := Ace; rank <= King; rank++ {
		for suit := Spades; suit <= Hearts; suit++ {
			d = append(d, Card{Rank: rank, Suit: suit})
		}
	}

	for _, o := range opts {
		d = o(d)
	}

	return d
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}

	return fmt.Sprintf("%s of %s", c.Rank, c.Suit)
}

// ShuffleRand randomizes the position of cards in the deck
func ShuffleRand(deck []Card) []Card {
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}

// Without removes any identical cards from the deck
func Without(cards ...Card) Option {
	// function to test exclusion
	exclude := func(c Card) bool {
		for _, card := range cards {
			if card == c {
				return true
			}
		}
		return false
	}

	// Build a new deck excluding the "without" cards
	return func(deck []Card) []Card {
		var newDeck []Card
		for _, c := range deck {
			if !exclude(c) {
				newDeck = append(newDeck, c)
			}
		}
		return newDeck
	}
}

// Compose multiple decks into one deck
func Compose(decks ...[]Card) []Card {
	var deck []Card
	if len(decks) > 0 {
		deck = decks[0]
	}

	// append any remaining decks
	for i := 1; i < len(decks)-1; i++ {
		deck = append(deck, decks[i]...)
	}

	return deck
}

// Include adds the requested cards to the new deck
func Include(cards ...Card) Option {
	return func(deck []Card) []Card {
		deck = append(deck, cards...)
		return deck
	}
}

// AddJokers adds num joker Cards to the deck
func AddJokers(num int) Option {
	return func(deck []Card) []Card {
		for i := 0; i < num; i++ {
			deck = append(deck, Card{Suit: Joker})
		}
		return deck
	}
}

// WithSort sorts the deck using the passed sort.Slice 'less' function
func WithSort(lessFn func(deck []Card) func(i, j int) bool) Option {
	return func(deck []Card) []Card {
		sort.Slice(deck, lessFn(deck))
		return deck
	}
}

// Sorts the deck in a consistent order
func DefaultSort(deck []Card) []Card {
	sort.Slice(deck, less(deck))
	return deck
}

func less(deck []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return value(deck[i]) < value(deck[j])
	}
}

func value(c Card) int {
	return int(c.Suit)*int(c.Rank) + int(c.Rank)
}

func count(deck []Card, card Card) int {
	result := 0
	for _, c := range deck {
		if c == card {
			result++
		}
	}

	return result
}
