package uno

import "log"

var debug = true

const (
	ColorRed     = 0
	ColorYellow  = 1
	ColorBlue    = 2
	ColorGreen   = 3
	ColorBlack   = 4 // or none
	ColorSpecial = 5

	CardZero     = 0
	CardOne      = 1
	CardTwo      = 2
	CardThree    = 3
	CardFour     = 4
	CardFive     = 5
	CardSix      = 6
	CardSeven    = 7
	CardEight    = 8
	CardNine     = 9
	CardSkip     = 10
	CardReverse  = 11
	CardDrawTwo  = 12
	Wild         = 13
	WildDrawFour = 14

	IDCardNone = 0

	IDCardRedSkip    = 10
	IDCardYellowSkip = 10 + 13
	IDCardGreenSkip  = 10 + 13*2
	IDCardBlueSkip   = 10 + 13*3

	IDCardRedReverse    = 11
	IDCardYellowReverse = 11 + 13
	IDCardGreenReverse  = 11 + 13*2
	IDCardBlueReverse   = 11 + 13*3

	IDCardRedDrawTwo    = 12
	IDCardYellowDrawTwo = 12 + 13
	IDCardGreenDrawTwo  = 12 + 13*2
	IDCardBlueDrawTwo   = 12 + 13*3

	IDCardRedNumZero    = 13
	IDCardYellowNumZero = 13 + 10
	IDCardGreenNumZero  = 13 + 13*2
	IDCardBlueNumZero   = 13 + 13*3

	IDCardWild        = 53
	IDCardWildAndDraw = 54

	SpecialDraw      = 60
	SpecialChallenge = 69
)

var colorMap = map[int]string{
	0: "red",
	1: "yellow",
	2: "green",
	3: "blue",
	4: "black",
	5: "speical",
}

func getColor(id int) int {
	switch {
	case id > 54:
		return ColorSpecial
	case id > 42:
		return ColorBlack
	case id > 39:
		return ColorBlue
	case id > 26:
		return ColorGreen
	case id > 13:
		return ColorYellow
	case id > 0:
		return ColorRed
	default:
		return 0
	}
}

func getName(id int) int {
	switch {
	case id == IDCardWildAndDraw:
		return WildDrawFour
	case id == IDCardWild:
		return Wild
	case id%13 == 0:
		return IDCardRedNumZero
	case id == SpecialDraw:
		return SpecialDraw
	case id == SpecialChallenge:
		return SpecialChallenge
	default:
		return id % 13
	}
}

const (
	GameModeStandard = "standard"
	GameModeWild     = "wild"

	ActionDraw = "draw"
	ActionCard = "card"
	ActionCall = "call"
)

func ulog(a ...interface{}) {
	if debug {
		log.Println(a)
	}
}
