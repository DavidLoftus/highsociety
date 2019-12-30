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
	luxuryCards   []*luxuryCard
	prestigeCards []*prestigeCard
	disgraceCards []*disgraceCard
	multiplier    float64
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
func (p *PlayerState) Status() int {
	status := 0

	for _, card := range p.luxuryCards {
		status += card.Status()
	}

	var scandaleCount uint = 0
	for _, card := range p.disgraceCards {
		switch card.Type() {
		case Scandale:
			scandaleCount++
		case Passe:
			status -= 5
		}
	}

	// Multiply status by 2 for each prestige card
	status <<= uint(len(p.prestigeCards))

	// Divide status by 2 for each scandale card
	status >>= scandaleCount

	return status
}

// GameState holds the state for an entire game
type GameState struct {
	players []*PlayerState
	deck    []*Card
}

func (game *GameState) GameOver() bool {
	return false
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

	return &GameState{
		deck:    NewDeck(),
		players: players,
	}
}
