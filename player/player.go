package player

type player struct {
	originList   []int
	reverse      bool
	currentIndex int
}

func (p *player) Count() int {
	return len(p.originList)
}

func (p *player) Index() int {
	return p.currentIndex
}

func (p *player) OrderList() []int {
	// Todo
	return []int{}
}

func (p *player) Get(rel int) int {
	switch {
	case rel == 1:
		return p.Next()
	case rel == 2:
		return p.Skip()
	case rel == -1:
		return p.Reverse()
	}

	return 0
}

func (p *player) Next() int {
	i := p.currentIndex
	if p.reverse {
		i--
	} else {
		i++
	}

	fixed := keepIndex(i, 0, p.Count())
	p.currentIndex = fixed
	return fixed
}

func (p *player) Skip() int {
	i := p.currentIndex
	if p.reverse {
		i -= 2
	} else {
		i += 2
	}

	fixed := keepIndex(i, 0, p.Count())
	p.currentIndex = fixed
	return fixed
}

func (p *player) Reverse() int {
	p.reverse = !p.reverse
	return p.Next()
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

// todo
func (p *player) register(ix int, key string) {

}

// NewGroup : -> *player
// todo
func NewGroup(count int) *player {
	return &player{
		originList: []int{1, 2, 3, 4, 5},
	}
}
