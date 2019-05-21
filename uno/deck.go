package uno

import (
	"math/rand"
)

// TurnResult detail of result when user post a new card
// but, not in use
type TurnResult struct {
	Valid    bool
	Message  string
	Cards    []int
	AltCards []int
}

type player struct {
	ID    int    `json:"id"`
	Key   string `json:"key"`
	cards []int  // not public
	// Todo
	Skiped bool // skip him, mostly because he leave the game?
	Locked bool // lock him, for cheat reasons?
}

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
	status    int

	reverse      bool
	currentIndex int
	players      []*player // Players of this game
}

// Todo: rewrite this method
// Accept : my card ID, all my cards, previous player's cards
func (d *deck) Accept(id int, cards, prevCards []int) (bool, []int, []int) {
	ulog("deck.Accept <<<<", id)
	lastID := d.LastID()
	if d.isValid(id) {
		d.graveyard = append(d.graveyard, id)
		switch id {
		case IDSpecialDraw:
			return true, d.Draw(2), []int{}
		case IDSpeicalDrawFour:
			return true, d.Draw(4), []int{}
		case IDSpecialChallenge:
			if isNotBluff(lastID, prevCards) {
				return true, d.Draw(6), []int{}
			}
			return true, []int{}, d.Draw(6)
		default:
			return true, []int{}, []int{}
		}
	}

	return false, []int{}, []int{}
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
	d.status = GameStatusGoing

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

func (d *deck) isValid(id int) bool {
	lastCard := d.LastCard()
	currentCard := Info(id)
	if currentCard.Name == lastCard.Name || currentCard.Color == lastCard.Color {
		return true
	}

	return false
}

func (d *deck) pickValidCards(ids []int) []int {
	lastCard := d.LastCard()
	nextColor := lastCard.NextColor()
	filteredCards := []int{}

	// wild_draw_four
	if nextColor == ColorBlack {
		filteredCards = append(filteredCards, IDSpeicalDrawFour, IDSpecialChallenge)
		return filteredCards
	}

	ulog("Last card", lastCard.String())
	filteredCards = []int{IDSpecialDraw}
	for _, i := range ids {
		if d.isValid(i) {
			filteredCards = append(filteredCards, i)
		}
	}
	ulog(">>>> FilteredCards", filteredCards)

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

func (d *deck) Gaming() bool {
	return d.status == GameStatusGoing
}

// NewDeck : => deck
// size is user count
func NewDeck(mode string, size int) *deck {
	d := &deck{
		Mode:  mode,
		cards: []int{3, 12, 22, 5, 7, 4, 3, 4, 5, 6, 7, 8, 11, 2, 12, 33, 41, 23, 22, 13, 44, 23, 22, 12, 24, 32, 12, 5, 24, 16, 19, 8, 2, 17, 22},
	}
	d.Shuffle()

	return d
}

func (d *deck) Count() int {
	return len(d.players)
}

func (d *deck) Reverse() int {
	d.reverse = !d.reverse
	return d.Next()
}

func (d *deck) Next() int {
	i := d.currentIndex
	if d.reverse {
		i--
	} else {
		i++
	}

	fixed := keepIndex(i, 0, d.Count())
	d.currentIndex = fixed
	return fixed
}

func (d *deck) Skip() int {
	i := d.currentIndex
	if d.reverse {
		i -= 2
	} else {
		i += 2
	}

	fixed := keepIndex(i, 0, d.Count())
	d.currentIndex = fixed
	return fixed
}

func keepIndex(value, min, max int) int {
	v := value
	for v < min {
		v += (max - min)
	}

	for v > max {
		v -= (max - min)
	}

	return v
}

// current position
func (d *deck) Index() int {
	return d.currentIndex
}

// NEXT player's position
func (d *deck) Get(rel int) int {
	switch {
	case rel == 1:
		return d.Next()
	case rel == 2:
		return d.Skip()
	case rel == -1:
		return d.Reverse()
	}

	return 0
}

// display all players
func (d *deck) OrderList() []*player {
	// Todo: check cards
	return d.players
}

func (d *deck) Join(key string) bool {
	playerNotIn := findPlayerWithKey(key, d.players) < 0
	if playerNotIn {
		d.players = append(d.players, &player{
			ID:  len(d.players) + 1,
			Key: key,
		})
		return true
	}

	return false
}

func (d *deck) Leave(key string) bool {
	// Todo; no one can leave, hahaha
	return false
}

func findPlayerWithKey(key string, pl []*player) int {
	for ix, p0 := range pl {
		if p0.Key == key {
			return ix
		}
	}

	return -1
}

func findPlayerWithID(id int, pl []*player) int {
	for ix, p0 := range pl {
		if p0.ID == id {
			return ix
		}
	}

	return -1
}
