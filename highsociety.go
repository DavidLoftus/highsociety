package highsociety

func getStartingNotes() []int {
	return []int{
		1, 2, 3, 4, 5,
		8, 10, 12, 15,
		20, 25,
	}
}

// PlayerState holds the state of each player
type PlayerState struct {
	notes         []int
	luxuryCards   []luxuryCard
	prestigeCards []prestigeCard
	disgraceCards []disgraceCard
}

func newPlayer() *PlayerState {
	return &PlayerState{
		notes: getStartingNotes(),
	}
}

// Balance gets the current balance of the player
func (p *PlayerState) Balance() int {
	sum := 0
	for _, noteValue := range p.notes {
		sum += noteValue
	}
	return sum
}

// Status returns the current total status of the player
func (p *PlayerState) Status() float64 {
	status := 0

	for _, card := range p.luxuryCards {
		status += card.Status()
	}

	multiplier := 1.0

	// Prestige cards multiply status by 2
	for _ = range p.prestigeCards {
		multiplier *= 2
	}

	for _, card := range p.disgraceCards {
		switch card.Type() {
		case Scandale:
			// Scandale cards halve the status
			multiplier /= 2
		case Passe:
			// Passe cards reduce status by 5
			status -= 5
		}
	}
	// Multiplier is applied at very end
	// Return a float64 since it could be a fraction
	return float64(status) * multiplier
}

func (p *PlayerState) getFauxPas() Card {
	for i, card := range p.disgraceCards {
		if card.Type() == FauxPas {
			// Remove card from slice
			l := len(p.disgraceCards)
			p.disgraceCards[i] = p.disgraceCards[l-1]
			p.disgraceCards = p.disgraceCards[:l-1]
			return card
		}
	}
	return nil
}

func (p *PlayerState) discard(card Card) {
	// TODO: display card in some discard pile
}

func removeMinCard(cards []luxuryCard) (Card, []luxuryCard) {
	i, minCard := 0, cards[0]
	for j, card := range cards {
		if card.Status() < minCard.Status() {
			i, minCard = j, card
		}
	}
	l := len(cards)

	cards[i] = cards[l]
	return minCard, cards[:l-1]
}

func (p *PlayerState) GiveCard(card Card) {
	switch card.Type() {
	case Luxury:
		if fauxPass := p.getFauxPas(); fauxPass != nil {
			p.discard(fauxPass)
			p.discard(card)
		} else {
			p.luxuryCards = append(p.luxuryCards, card.(luxuryCard))
		}
	case Prestige:
		p.prestigeCards = append(p.prestigeCards, card.(prestigeCard))
	case FauxPas:
		if len(p.luxuryCards) > 0 {
			// Technically player is meant to choose here...
			// For now lets just remove the smallest value card
			var minCard Card
			minCard, p.luxuryCards = removeMinCard(p.luxuryCards)

			p.discard(card)
			p.discard(minCard)

			break
		}
		fallthrough
	case Passe:
		fallthrough
	case Scandale:
		p.disgraceCards = append(p.disgraceCards, card.(disgraceCard))
	}
}

// GameState holds the state for an entire game
type GameState struct {
	players []*PlayerState
	deck    []Card
}

func (game *GameState) GameOver() bool {
	return len(game.deck) == 0
}

func (game *GameState) CurrentPlayer() int {
	return 0
}

// NewGame creates a new game state
func NewGame(numUsers int) *GameState {
	players := make([]*PlayerState, numUsers)
	for i := range players {
		players[i] = newPlayer()
	}

	deck := NewDeck()
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	return &GameState{
		deck:    deck,
		players: players,
	}
}
