package uno

import "log"

var debug = true

// GameModes
const (
	GameModeStandard = "standard"
	GameModeWild     = "wild"

	// Todo: rewrite here
	// GameStatusXXX for deck
	GameStatusOpen  = 0
	GameStatusGoing = 1

	// Colors, Cards, IDs
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
	// One
	IDCardRedNumOne    = 1        // 1
	IDCardYellowNumOne = 1 + 13   // 14
	IDCardGreenNumOne  = 1 + 13*2 // 26
	IDCardBlueNumOne   = 1 + 13*3 // 40
	// Skip
	IDCardRedSkip    = 10        // 10
	IDCardYellowSkip = 10 + 13   // 23
	IDCardGreenSkip  = 10 + 13*2 // 36
	IDCardBlueSkip   = 10 + 13*3 // 49
	// Reverse
	IDCardRedReverse    = 11        // 11
	IDCardYellowReverse = 11 + 13   // 24
	IDCardGreenReverse  = 11 + 13*2 // 37
	IDCardBlueReverse   = 11 + 13*3 // 50
	// DrawTwo
	IDCardRedDrawTwo    = 12        // 12
	IDCardYellowDrawTwo = 12 + 13   // 25
	IDCardGreenDrawTwo  = 12 + 13*2 // 38
	IDCardBlueDrawTwo   = 12 + 13*3 // 51
	// Zero
	IDCardRedNumZero    = 13        // 13
	IDCardYellowNumZero = 13 + 13   // 26
	IDCardGreenNumZero  = 13 + 13*2 // 39
	IDCardBlueNumZero   = 13 + 13*3 // 52
	// Wild
	IDCardWild        = 53
	IDCardWildAndDraw = 54
	// Draw
	IDSpecialDraw = 60
	// 61-64 is x-wild
	IDWildRed    = 61
	IDWildYellow = 62
	IDWildGreen  = 63
	IDWildBlue   = 64
	// 65-68 is x-wildDrawFour
	IDWildDrawFourRed    = 65
	IDWildDrawFourYellow = 66
	IDWildDrawFourGreen  = 67
	IDWildDrawFourBlue   = 68
	IDSpecialChallenge   = 69
	IDSpeicalDrawFour    = 70

	IDFakeCardRed    = 71 // it's used cards for drawTwo
	IDFakeCardYellow = 72
	IDFakeCardGreen  = 73
	IDFakeCardBlue   = 74
)

var colorMap = map[int]string{
	0: "red",
	1: "yellow",
	2: "green",
	3: "blue",
	4: "black",
	5: "speical",
}

func getAltColor(id int) int {
	return id - IDWildDrawFourRed
}

func getColor(id int) int {
	switch {
	case id >= IDSpecialDraw:
		return ColorSpecial
	case id >= IDWildDrawFourRed:
		return id - IDWildDrawFourRed
	case id >= IDWildRed:
		return id - IDWildRed
	case id == IDCardWild:
		return ColorBlack
	case id == IDCardWildAndDraw:
		return ColorBlack
	case id >= IDCardBlueNumOne:
		return ColorBlue
	case id >= IDCardGreenNumOne:
		return ColorGreen
	case id > IDCardYellowNumOne:
		return ColorYellow
	case id > IDCardRedNumOne:
		return ColorRed
	case cardIsFake(id):
		return id - IDFakeCardRed
	default:
		return ColorSpecial
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
	case id == IDSpecialDraw:
		return IDSpecialDraw
	case id == IDSpeicalDrawFour:
		return IDSpeicalDrawFour
	case id == IDSpecialChallenge:
		return IDSpecialChallenge
	case cardIsFake(id):
		// todo: check it
		return getColor(id)*13 + IDCardRedDrawTwo // fake drawTwo
	default:
		return id % 13
	}
}

// helper function to display logs
func ulog(a ...interface{}) {
	if debug {
		log.Println(a...)
	}
}
