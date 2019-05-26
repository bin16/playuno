package uno

import (
	"math/rand"
	"sort"
	"strconv"
	"time"
)

type player struct {
	_ID   int    // 1,2,3...
	ID    string `json:"id"`   // uuid, like password
	Name  string `json:"name"` // display name/nickname
	cards []int

	// deprecated
	Key string `json:"key"`
}

func (p *player) Cards() []int {
	results := p.cards[:]
	sort.Ints(results)
	return results
}

func (p *player) AddCards(cards []int) {
	p.cards = append(p.cards, cards...)
}

func (p *player) RemoveCard(cardID int) {
	for ix, c := range p.cards {
		if c == cardID {
			p.cards = append(p.cards[:ix], p.cards[ix+1:]...)
		}
	}
}

type deck struct {
	Mode          string `json:"mode"`
	cards         []int
	graveyard     []int
	status        int
	reverse       bool
	previousIndex int
	currentIndex  int
	players       []*player
	lock          bool
}

// MyCards something named like this,
// return cards by player id
// todo, need a new data type for card
// with attributes like
// color, image url, available...
// it would be helpful for client
// func (d *deck) MyCards(id int) []Card {
// 	cards := []Card{}
// 	return cards
// }

// Players return all player in game
// same as d.players
func (d *deck) Players() []*player {
	return d.players
}

func (d *deck) PlayerNames() []string {
	names := []string{}
	for _, p := range d.players {
		names = append(names, p.Name+" ("+strconv.Itoa(p._ID)+")")
	}

	return names
}

func (d *deck) CurrentPlayer() *player {
	return d.players[d.currentIndex]
}

func (d *deck) PreviousPlayer() *player {
	return d.players[d.previousIndex]
}

// Accept get card ID
func (d *deck) Accept(playerID string, cardID int) (bool, []UnoMsg) {
	p1 := d.CurrentPlayer()
	if p1.ID != playerID {
		return false, []UnoMsg{
			*unoMsgMaker(false, MsgWarning).To(playerID),
		}
	}
	p0 := d.PreviousPlayer() // p1 may challenge him
	lastID := d.LastID()     // last card id in graveyard
	if d.isUsable(cardID) {
		p1.RemoveCard(cardID)
		d.graveyard = append(d.graveyard, cardID)
		switch cardID {
		case IDSpecialDraw:
			p1.AddCards(d.Draw(2))
			return OK, []UnoMsg{
				*unoMsgMaker(true, MsgCardAccept),
				*unoMsgMaker(true, CmdSetCards).WithCards(p1.Cards(), NoCards).To(p1.ID), // Update p1's cards
			}
		case IDSpeicalDrawFour:
			p1.AddCards(d.Draw(4))
			return OK, []UnoMsg{
				*unoMsgMaker(true, MsgCardAccept),
				*unoMsgMaker(true, CmdSetCards).WithCards(p1.Cards(), NoCards).To(p1.ID), // Update p1's cards
			}
		case IDSpecialChallenge:
			if isNotBluff(lastID, p0.cards) {
				p1.AddCards(d.Draw(6))
				return OK, []UnoMsg{
					*unoMsgMaker(true, MsgCardAccept),
					*unoMsgMaker(true, CmdSetCards).WithCards(p1.Cards(), NoCards).To(p1.ID), // Update p1's cards
					*unoMsgMaker(true, MsgPlayerNotBluff).WithTarget(p0),
					*unoMsgMaker(true, MsgPlayerGotCards).WithTarget(p1),
				}
			}
			// else
			p0.AddCards(d.Draw(6))
			return OK, []UnoMsg{
				*unoMsgMaker(true, MsgCardAccept),
				*unoMsgMaker(true, CmdSetCards).WithCards(p0.Cards(), NoCards).To(p0.ID), // Update p0's cards
				*unoMsgMaker(true, CmdSetCards).WithCards(p1.Cards(), NoCards).To(p1.ID), // Update p1's cards
				*unoMsgMaker(true, MsgPlayerBluff).WithTarget(p0),
				*unoMsgMaker(true, MsgPlayerGotCards).WithTarget(p0),
			}
		default:
			return OK, []UnoMsg{
				*unoMsgMaker(true, MsgCardAccept),
				*unoMsgMaker(true, CmdSetCards).WithCards(p1.Cards(), NoCards).To(p1.ID), // Update p1's cards
			}
		}
	}

	return false, []UnoMsg{
		*unoMsgMaker(false, MsgWarning).To(playerID),
	}
}

// Todo: fix
func (d *deck) NextTurn() []UnoMsg {
	d.previousIndex = d.currentIndex
	d.currentIndex = d.IndexNextPlayer()
	p1 := d.CurrentPlayer()

	return []UnoMsg{
		*unoMsgMaker(true, MsgPlayerToGame).WithTarget(p1),
		*unoMsgMaker(true, CmdSetCards).WithCards(p1.Cards(), d.Filter(p1.ID)).To(p1.ID),
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
	d.cards = append(cards, d.cards...)

	return cards
}

func (d *deck) Draw(num int) []int {
	ulog("drawing", num, "cards.")
	if len(d.cards) < num {
		d.ShuffleN(num * 2)
	}
	ids := d.cards[0:num]
	d.cards = d.cards[num+1:]

	return ids
}

func (d *deck) StartBy(playerID string) (bool, []UnoMsg) {
	if d.Gaming() {
		return false, []UnoMsg{*unoMsgMaker(false, MsgWarning).To(playerID)}
	}

	if d.Player(playerID) == nil {
		return false, []UnoMsg{*unoMsgMaker(false, MsgWarning).To(playerID)}
	}

	return OK, d.start()
}

func (d *deck) start() []UnoMsg {
	d.ShuffleN(100)
	c := d.cards[0]
	d.cards = d.cards[1:]
	d.graveyard = append(d.graveyard, c)
	d.currentIndex = d.CountPlayers() - 1

	c0 := Info(c)
	ulog("First card is", c0.String())
	d.status = GameStatusGoing

	msgList := []UnoMsg{}
	for _, p := range d.players {
		p.AddCards(d.Draw(7))
		msgList = append(msgList,
			*unoMsgMaker(true, CmdSetCards).WithCards(p.Cards(), NoCards).To(p.ID),
		)
	}
	msgList = append(msgList,
		*unoMsgMaker(true, MsgCardAccept),
	)

	return msgList
}

// LastID returns last card in graveyard's ID
func (d *deck) LastID() int {
	return d.graveyard[len(d.graveyard)-1]
}

// LastCard returns last Card in graveyard
func (d *deck) LastCard() Card {
	return Info(d.LastID())
}

func (d *deck) isUsable(id int) bool {
	lastCard := d.LastCard()
	currentCard := Info(id)
	if currentCard.Name == lastCard.Name || currentCard.Color == lastCard.Color {
		return true
	}

	return false
}

func (d *deck) pickActiveCards(ids []int) []int {
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
		if d.isUsable(i) {
			filteredCards = append(filteredCards, i)
		}
	}
	ulog(">>>> FilteredCards", filteredCards)

	return filteredCards
}

// Filter returns valid cards can be post
func (d *deck) Filter(playerID string) []int {
	p := d.CurrentPlayer()
	if p.ID != playerID {
		return []int{}
	}

	return d.pickActiveCards(p.Cards())
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
		Mode: mode,
	}
	d.cards = d.ShuffleN(54)

	return d
}

// Count returns how many players in game
// Todo: some players may leave, but still in players: []*player
// they should be skipped forever, so fix here
func (d *deck) CountPlayers() int {
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

	fixed := keepIndex(i, 0, d.CountPlayers())
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

	fixed := keepIndex(i, 0, d.CountPlayers())
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

	for v >= max {
		v -= (max - min)
	}

	return v
}

// current position
func (d *deck) CurrentIndex() int {
	return d.currentIndex
}

// todo
// NEXT player's position
func (d *deck) NextPlayer(rel int) int {
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

// // display all players
// func (d *deck) OrderList() []*player {
// 	// Todo: check cards
// 	return d.players
// }

func (d *deck) Join(id, name string) []UnoMsg {
	if d.IsPlayerIn(id) {
		return []UnoMsg{
			UnoMsg{
				ID:      id,
				Ok:      false,
				Message: MsgWarning, // todo: ErrPlayerIsIn
			},
		}
	}

	if d.lock && d.Gaming() {
		return []UnoMsg{
			UnoMsg{
				ID:      id,
				Ok:      false,
				Message: MsgWarning, // todo: ErrGameLocked
			},
		}
	}

	p := &player{
		_ID:  d.CountPlayers() + 1,
		ID:   id,
		Name: name,
		Key:  id,
	}
	d.players = append(d.players, p)
	messages := []UnoMsg{
		UnoMsg{
			Ok:      true,
			Message: MsgGamePlayerJoin,
			Players: d.PlayerNames(),
		},
	}
	if d.Gaming() {
		p.AddCards(d.Draw(7))
		messages = append(messages, UnoMsg{
			ID:      id,
			Ok:      true,
			Cards:   p.Cards(),
			Message: MsgGameDrawCards,
		})
	}

	return messages
}

func findPlayerWithID(id string, pl []*player) int {
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

func (d *deck) IsPlayerIn(id string) bool {
	return findPlayerWithID(id, d.players) >= 0
}

func (d *deck) Player(id string) *player {
	for _, p := range d.players {
		if p.ID == id {
			return p
		}
	}

	return nil
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

func isCardInList(card int, cards []int) bool {
	for _, c := range cards {
		if c == card {
			return true
		}
	}

	return false
}
