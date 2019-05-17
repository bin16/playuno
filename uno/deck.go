package uno

import "math/rand"

// deck - how it works:
// ... start ...
// deck.Shuffle(Num)
// deck.Start()
// deck.Draw(7) x N times
// deck.NextPlayer()
// deck.Filter(playerCardIDs)
// ... waiting ...
// deck.Accept(CardID), err to invalid
// deck.NextPlayer()
// deck.Filter(playerCardIDs)
// ... waiting ...
// ...
type deck struct {
	cards     []int
	graveyard []int
	usedCards []int
	Mode      string
}

func (d *deck) Shuffle() {
	l := len(d.cards)
	ulog(l, "cards shuffling.")
	rand.Shuffle(l, func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

func (d *deck) ShuffleN(n int) {
	// make n cards, then shuffle
	// todo
}

func (d *deck) Draw(num int) []int {
	ulog("drawing", num, "cards.")
	ids := d.cards[0:num]
	d.cards = d.cards[num+1:]

	return ids
}

func (d *deck) Start() int {
	c := d.cards[0]
	d.cards = d.cards[1:]
	d.graveyard = append(d.graveyard, c)

	cc := Info(c)
	ulog("First card is", cc.String())

	return c
}

func (d *deck) Remove(num int) {
	ids := d.Draw(num)
	d.graveyard = append(d.graveyard, ids...)
}

func (d *deck) LastID() int {
	return d.graveyard[len(d.graveyard)-1]
}

func (d *deck) LastCard() Card {
	return Info(d.LastID())
}

func (d *deck) pickValidCards(ids []int) []int {
	lastCard := d.LastCard()
	nextColor := lastCard.NextColor()
	filteredCards := []int{SpecialDraw}

	// wild_draw_four
	if nextColor == ColorBlack {
		filteredCards = append(filteredCards, SpecialChallenge)
		return filteredCards
	}

	for _, i := range ids {
		c := Info(i)
		if c.Name == lastCard.Name || c.Color == nextColor {
			filteredCards = append(filteredCards, i)
		}
	}

	return filteredCards
}

func (d *deck) Filter(ids []int) []int {
	return d.pickValidCards(ids)
}

func (d *deck) NextPlayer() int {
	lastID := d.LastID()
	switch {
	case cardIsReverse(lastID):
		return -1 // previous
	case cardIsSkip(lastID):
		return 2 // next's next
	default:
		return 1 // next
	}
}

func (d *deck) findRelatedCards() []int {
	relatedCards := []int{}
	lastIndex := len(d.graveyard) - 1
	lastID := d.graveyard[lastIndex]
	// +2 +2 +2 ...
	if cardIsDrawTwo(lastID) {
		for i := lastIndex; i >= 0; i-- {
			id := d.graveyard[i]
			if cardIsDrawTwo(id) {
				relatedCards = append(relatedCards, id)
			} else {
				return relatedCards
			}
		}
	}

	return []int{lastID}
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
