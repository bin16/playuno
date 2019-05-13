package uno

import "math/rand"

type deck struct {
	cards     []int
	graveyard []int
	usedCards []int
	Mode      string
}

func (d *deck) Shuffle() {
	l := len(d.cards)
	rand.Shuffle(l, func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

func (d *deck) Draw(num int) []int {
	ids := d.cards[0:num]
	d.cards = d.cards[num+1:]

	return ids
}

func (d *deck) Start() int {
	c := d.cards[0]
	d.cards = d.cards[1:]
	d.graveyard = append(d.graveyard, c)

	return c
}

func (d *deck) pickValidCards(ids []int) []int {
	return ids
}

func (d *deck) findRelatedCards() []int {
	relatedCards := []int{}
	for i := len(d.graveyard) - 1; i >= 0; i-- {
		id := d.graveyard[i]
		if cardIsDrawTwo(id) {
			relatedCards = append(relatedCards, id)
		}
	}

	return relatedCards
}

// NewDeck : => deck
// size is user count
func NewDeck(mode string, size int) *deck {
	d := &deck{
		Mode:  mode,
		cards: []int{5, 7, 4, 3, 4, 5, 6, 7, 8, 11, 2, 12, 33, 41, 23, 22, 13, 44},
	}
	d.Shuffle()

	return d
}
