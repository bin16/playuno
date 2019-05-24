package uno

import (
	"math/rand"
	"time"
)

// Messages
const (
	MsgPlayerNotBluff = "playerNotBluff"
	MsgPlayerBluff    = "playerBluff"
	MsgPlayerCheating = "playerCheating"
	MsgOk             = "ok"
	MsgPlayerToGame   = "playerToGame"
	MsgGamePlayerJoin = "gamePlayerJoin"
	MsgGameStart      = "gameStart"
	MsgSystemMessage  = "systemMessage"
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
	Name  string `json:"name"`
	ID    int    `json:"id"`
	Key   string `json:"key"`
	cards []int
	// Todo
	skip bool // skip him, mostly because he leave the game?
	lock bool // lock him, for cheat reasons?
}

func (p *player) Cards() []int {
	return p.cards
}

func (p *player) AddCards(cards []int) {
	p.cards = append(p.cards, cards...)
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

	reverse       bool
	previousIndex int
	currentIndex  int
	players       []*player // Players of this game
}

// MyCards something named like this,
// return cards by player id
// todo, need a new data type for card
// with attributes like
// color, image url, available...
// it would be helpful for client
func (d *deck) MyCards(id int) {

}

func (d *deck) Players() []*player {
	return d.players
}

func (d *deck) GetID(key string) int {
	for _, p := range d.players {
		if p.Key == key {
			return p.ID
		}
	}

	return 0
}

func (d *deck) Player(id int) *player {
	for _, p := range d.players {
		if p.ID == id {
			return p
		}
	}

	return &player{}
}

func (d *deck) CurrentPlayer() *player {
	return d.players[d.currentIndex]
}

func (d *deck) PreviousPlayer() *player {
	return d.players[d.previousIndex]
}

type unoMsg struct {
	Ok      bool
	Cards   []int
	Message string
	Key     string
	p0cards []int
	p1cards []int
}

// Accept : my card ID
func (d *deck) Accept(id int) unoMsg {
	p1 := d.CurrentPlayer()
	p0 := d.PreviousPlayer() // you may challenge him
	ulog("deck.Accept <<<<", id)
	lastID := d.LastID() // last card id in graveyard
	if d.isValid(id) {
		d.graveyard = append(d.graveyard, id)
		switch id {
		case IDSpecialDraw:
			p1.cards = append(p1.cards, d.Draw(2)...)
			return unoMsg{
				Ok:      true,
				Cards:   p1.cards,
				Message: MsgOk,
				p0cards: p0.cards,
				p1cards: p1.cards,
			}
		case IDSpeicalDrawFour:
			p1.cards = append(p1.cards, d.Draw(4)...)
			return unoMsg{
				Ok:      true,
				Cards:   p1.cards,
				Message: MsgOk,
				p0cards: p0.cards,
				p1cards: p1.cards,
			}
		case IDSpecialChallenge:
			if isNotBluff(lastID, p0.cards) {
				p1.cards = append(p1.cards, d.Draw(6)...)
				return unoMsg{
					Ok:      true,
					Cards:   p1.cards,
					Message: MsgPlayerNotBluff,
					p0cards: p0.cards,
					p1cards: p1.cards,
				}
			}

			p0.cards = append(p0.cards, d.Draw(6)...)
			return unoMsg{
				Ok:      true,
				Cards:   p1.cards,
				Message: MsgPlayerBluff,
				p0cards: p0.cards,
				p1cards: p1.cards,
			}
		default:
			return unoMsg{
				Ok:      true,
				Cards:   p1.cards,
				Message: MsgOk,
				p0cards: p0.cards,
				p1cards: p1.cards,
			}
		}
	}

	return unoMsg{
		Ok:      false,
		Cards:   p1.cards,
		Message: MsgPlayerCheating,
	}
}

// Todo: fix
func (d *deck) NextTurn() unoMsg {
	d.previousIndex = d.Index()
	d.currentIndex = d.IndexNextPlayer()
	cards := d.CurrentPlayer().cards
	d.Filter(cards)
	return unoMsg{
		Ok:      true,
		Message: MsgPlayerToGame,
		Cards:   cards,
		p1cards: cards,
	}
}

func (d *deck) Shuffle() {
	l := len(d.cards)
	ulog(l, "cards shuffling.")
	rand.Shuffle(l, func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

//Todo: finish it
func (d *deck) ShuffleN(n int) []int {
	cards := []int{}
	for i := 0; i < n; i++ {
		rand.Seed(int64(time.Now().Nanosecond() * (i + 1)))
		id := rand.Intn(54) + 1
		cards = append(cards, id)
	}

	return cards
}

func (d *deck) Draw(num int) []int {
	ulog("drawing", num, "cards.")
	ids := d.cards[0:num]
	d.cards = d.cards[num+1:]

	return ids
}

func (d *deck) Start() unoMsg {
	d.ShuffleN(100)
	c := d.cards[0]
	d.cards = d.cards[1:]
	d.graveyard = append(d.graveyard, c)

	c0 := Info(c)
	ulog("First card is", c0.String())
	d.status = GameStatusGoing

	for _, p := range d.players {
		p.cards = d.Draw(7)
	}

	return unoMsg{
		Ok:      true,
		Message: MsgGameStart,
	}
}

func (d *deck) Remove(num int) {
	ids := d.Draw(num)
	d.graveyard = append(d.graveyard, ids...)
}

// LastID returns last card in graveyard's ID
func (d *deck) LastID() int {
	return d.graveyard[len(d.graveyard)-1]
}

// LastCard returns last Card in graveyard
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

// Filter returns valid cards can be post
func (d *deck) Filter(ids []int) []int {
	return d.pickValidCards(ids)
}

// NextPlayer returns the next player based on cards
// Todo: think it should be dropped, use IndexNextPlayer()
// and IndexNextPlayer() used in .NextTurn()
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

// findRelatedCards return related cards ...
// +2 +2 +2 (mostly?)
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

// Gaming returns if game's status is playing
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

// Count returns how many players in game
// Todo: some players may leave, but still in players: []*player
// they should be skipped forever, so fix here
func (d *deck) Count() int {
	return len(d.players)
}

// Reverse change gaming order
func (d *deck) Reverse() int {
	d.reverse = !d.reverse
	return d.Next()
}

// Next always return index+1...
// todo
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

// Skip returns index+2...
// todo
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

// to get next player[To Post Card]'s index
func (d *deck) IndexNextPlayer() int {
	lastCardID := d.LastID()
	switch {
	case cardIsSkip(lastCardID):
		return d.Skip()
	case cardIsReverse(lastCardID):
		return d.Reverse()
	default:
		return d.Next()
	}
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

func (d *deck) Export() deckInfo {
	return deckInfo{
		Cards:         d.cards,
		Graveyard:     d.graveyard,
		Mode:          d.Mode,
		Status:        d.status,
		Reverse:       d.reverse,
		PreviousIndex: d.previousIndex,
		CurrentIndex:  d.currentIndex,
		Players:       d.players,
	}
}

type deckInfo struct {
	Cards     []int `json:"cards"`
	Graveyard []int `json:"graveyard"`
	UsedCards []int
	Mode      string `json:"mode"`
	Status    int    `json:"status"`

	Reverse       bool
	PreviousIndex int
	CurrentIndex  int

	Players []*player `json:"players"`
	MyCards []int     `json:"myCards"`
}
