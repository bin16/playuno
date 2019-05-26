package uno

// Messages
const (
	MsgOk             = "ok"
	MsgPlayerCheating = "playerCheating" // ignore
	MsgPlayerMessage  = "playerMessage"
	MsgGamePlayerJoin = "gamePlayerJoin"
	MsgGameStart      = "gameStart"
	MsgGameDrawCards  = "gameDrawCards"
	MsgGameFirstCard  = "gameFirstCard"
	MsgGameCard       = "gameCard"
	MsgOneToGame      = "oneToGame"
	MsgSystemMessage  = "systemMessage"

	MsgPlayerWin = "msg_player_win"

	MsgWarning        = "msg_warning"
	MsgCardAccept     = "msg_card_accept"
	MsgPlayerBluff    = "msg_player_bluff"
	MsgPlayerNotBluff = "msg_player_not_bluff"
	MsgPlayerGotCards = "msg_player_got_cards"
	MsgPlayerToGame   = "msg_player_to_game"
	CmdSetCards       = "cmd_set_cards"   // client update local cards
	CmdSetPlayers     = "cmd_set_players" // client update local players

	OK = true
)

type UnoMsg struct {
	ID      string       // send to that player
	Message string       `json:"msg"`
	MyCards []MyCard     `json:"myCards"`
	Target  UnoMsgPlayer `json:"target"`

	Name        string   `json:"name"`
	Ok          bool     `json:"ok"`
	Cards       []int    `json:"cards"`
	ActiveCards []int    `json:"activeCards"`
	Players     []string `json:"players"`
}

type UnoMsgPlayer struct {
	id   string
	Name string `json:"name"`
}

type MyCard struct {
	ID     int
	Usable bool
}

func (u *UnoMsg) To(id string) *UnoMsg {
	u.ID = id
	return u
}

func (u *UnoMsg) WithTarget(p *player) *UnoMsg {
	u.Target = UnoMsgPlayer{
		id:   p.ID,
		Name: p.Name,
	}
	return u
}

func (u *UnoMsg) WithCards(cards, activeCards []int) *UnoMsg {
	mcs := []MyCard{}
	for _, c := range cards {
		mc := MyCard{
			ID:     c,
			Usable: isCardInList(c, activeCards),
		}
		mcs = append(mcs, mc)
	}

	u.MyCards = mcs
	return u
}

func unoMsgMaker(ok bool, message string) *UnoMsg {
	return &UnoMsg{
		Ok:      ok,
		Message: message,
	}
}
