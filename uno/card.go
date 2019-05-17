package uno

import (
	"fmt"
	"strconv"
)

/*
   	===========================================================
	|     000 | ------- |                                     |
    | 001~013 |     Red | 1-9,skip,reverse,draw_two,0         |
    | 014~026 |  Yellow | 1-9,skip,reverse,draw_two,0         |
    | 027~039 |   Green | 1-9,skip,reverse,draw_two,0         |
    | 040~052 |    Blue | 1-9,skip,reverse,draw_two,0         |
	|     053 |    Fake | wild                                |
	|     054 |    Fake | wild_and_draw                       |
	|---------|---------|-------------------------------------|
	|     060 | Special | draw                                |
	| 061~064 | Special | wild:red,yellow,green,blue          |
	| 065~068 | Special | wild_and_draw:red,yellow,green,blue |
	|     069 | Special | challenge                           |
	===========================================================
*/

// Card is UNO Card
type Card struct {
	ID    int
	Color int
	Name  int
}

func (c *Card) String() string {
	switch c.ID {
	case SpecialDraw:
		return "Special - Draw"
	case SpecialChallenge:
		return "Special - Challenge"
	}

	var color, name string
	switch c.Color {
	case 0:
		color = "Red"
	case 1:
		color = "Yellow"
	case 2:
		color = "Green"
	case 3:
		color = "Blue"
	default:
		color = "Black/Special"
	}

	switch val := c.Name % 13; val {
	case 10:
		name = "Skip"
	case 11:
		name = "Reverse"
	case 12:
		name = "Draw Two"
	default:
		name = strconv.Itoa(val)
	}

	return fmt.Sprintf("[%d]: %s %s", c.ID, color, name)
}

// IsNormal : return if it's r/y/g/b numbers,skip,reverse
// [ draw two ] is special
func (c *Card) IsNormal() bool {
	if c.Color < 4 && (c.Name < 12 || c.Name == 13) {
		return true
	}

	return false
}

// NextColor : return true color for next player
func (c *Card) NextColor() int {
	// wild only, because wild_and_Draw is different
	if c.Name == Wild {
		return c.ID - 61 // r,y,g,b
	}

	return c.Color
}

// Info = (id) => Card
func Info(id int) Card {
	return Card{
		ID:    id,
		Color: getColor(id),
		Name:  getName(id),
	}
}

func cardIsNumber(id int) bool {
	if id > IDCardNone && id < IDCardWild {
		val := id % 13
		if val < IDCardRedSkip || val == IDCardRedNumZero {
			return true
		}
	}

	return false
}

func cardIsWildDrawFour(id int) bool {
	if id == 54 { // its not needed
		return true
	}

	return id >= 65 && id <= 68
}

func cardIsSkip(id int) bool {
	switch id {
	case IDCardRedSkip:
		return true
	case IDCardYellowSkip:
		return true
	case IDCardGreenSkip:
		return true
	case IDCardBlueSkip:
		return true
	default:
		return false
	}
}

func cardIsReverse(id int) bool {
	switch id {
	case IDCardRedReverse:
		return true
	case IDCardYellowReverse:
		return true
	case IDCardGreenReverse:
		return true
	case IDCardBlueReverse:
		return true
	default:
		return false
	}
}

func cardIsDrawTwo(id int) bool {
	switch id {
	case IDCardRedDrawTwo:
		return true
	case IDCardYellowDrawTwo:
		return true
	case IDCardGreenDrawTwo:
		return true
	case IDCardBlueDrawTwo:
		return true
	default:
		return false
	}
}
