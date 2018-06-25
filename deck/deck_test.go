package deck

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleDeck() {
	fmt.Println(Card{Suit: Joker})
	fmt.Println(Card{Rank: Ace, Suit: Clubs})
	fmt.Println(Card{Rank: Three, Suit: Spades})
	fmt.Println(Card{Rank: King, Suit: Diamonds})
	fmt.Println(Card{Rank: Seven, Suit: Hearts})

	// Output:
	// Joker
	// Ace of Clubs
	// Three of Spades
	// King of Diamonds
	// Seven of Hearts
}

func TestNew(t *testing.T) {
	d := New()
	assert.Equal(t, 52, len(d), "Default deck should have 52 cards.")
}

func TestDefaultSort(t *testing.T) {
	deck := New(DefaultSort)
	assert.Equal(t, Card{Rank: Ace, Suit: Spades}, deck[0])
	assert.Equal(t, Card{Rank: King, Suit: Hearts}, deck[len(deck)-1])
}

func TestShuffle(t *testing.T) {
	var different bool
	unshuffled, shuffled := New(), New(ShuffleRand)
	assert.Equal(t, len(shuffled), len(unshuffled), "shuffling shouldn't change length")

	// Check if any card is in a different placement
	for i := 0; i < len(unshuffled); i++ {
		if shuffled[i] != unshuffled[i] {
			different = true
			break
		}
	}
	assert.True(t, different, "At least one card should have changed positions.")
}

func TestWithout(t *testing.T) {
	deck := New(Without(Card{Diamonds, Ace}))
	for _, c := range deck {
		if c.Rank == Ace && c.Suit == Diamonds {
			t.Errorf("Card %s should not be in the deck", c)
		}
	}
}

func TestInclude(t *testing.T) {
	extraCard := Card{Spades, Ace}
	d1, d2 := New(), New(Include(extraCard))
	assert.Equal(t, count(d1, extraCard)+1, count(d2, extraCard))
}

func TestAddJokers(t *testing.T) {
	deck := New(AddJokers(5))
	assert.Equal(t, 5, count(deck, Card{Suit: Joker}))
}
