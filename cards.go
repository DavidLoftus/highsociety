package highsociety

// CardType is an enum type for all different kinds of cards.
type CardType int

const (
	// Luxury cards, regular cards with a status value.
	Luxury CardType = iota
	// Prestige cards double your status.
	Prestige
	// Passe cards subtract 5 from your status.
	Passe
	// Scandale cards halve your status.
	Scandale
	// FauxPas cards force you to discard a luxury card.
	FauxPas
)

// Card is an abstract type to represent any of the Card types.
type Card interface {
	Type() CardType
}

type prestigeCard struct{}

func (prestigeCard) Type() CardType {
	return Prestige
}

type disgraceCard CardType

func (d disgraceCard) Type() CardType {
	return CardType(d)
}

type luxuryCard int

func (luxuryCard) Type() CardType {
	return Luxury
}

func (l luxuryCard) Status() int {
	return int(l)
}

// NewDeck creates an array of unshuffled cards
func NewDeck() []*Card {
	return []*Card{}
}
