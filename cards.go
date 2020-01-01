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

// Integer indicates which prestige card it is (for images)
type prestigeCard int

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
func NewDeck() []Card {
	return []Card{
		luxuryCard(1),
		luxuryCard(2),
		luxuryCard(3),
		luxuryCard(4),
		luxuryCard(5),
		luxuryCard(6),
		luxuryCard(7),
		luxuryCard(8),
		luxuryCard(9),
		luxuryCard(10),
		prestigeCard(0),
		prestigeCard(1),
		prestigeCard(2),
		disgraceCard(Passe),
		disgraceCard(Scandale),
		disgraceCard(FauxPas),
	}
}
