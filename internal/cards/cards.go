package cards

import (
	"errors"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrBadCard is when a card isn't valid
	ErrBadCard = errors.New("bad card")
	// ErrNoCardsAvailable is when there is no cards in a deck
	ErrNoCardsAvailable = errors.New("no cards in deck")
)

// Card is a card
type Card struct {
	name string
}

// NewCard returns a card from a filename
func NewCard(file string) (Card, error) {
	if strings.HasSuffix(file, ".jpg") {
		baseName := filepath.Base(file)
		cardName := strings.TrimSuffix(baseName, ".jpg")
		return Card{name: cardName}, nil
	}
	return Card{}, ErrBadCard
}

// Front return filname containing front of card
func (c *Card) Front() string {
	return c.name + ".jpg"
}

// Back returns filename containing back of card
func (c *Card) Back() string {
	return c.name + "-text.png"
}

// Name the name of the card
func (c *Card) Name() string {
	return c.name
}

// CardDeck is a deck of Card
type CardDeck struct {
	Directory string
	cards     []Card
}

// NewCardDeck generate a deck of cards from a directory
func NewCardDeck(dir string) (*CardDeck, error) {
	deck := &CardDeck{Directory: dir}
	if err := deck.populate(); err != nil {
		return nil, err
	}
	if deck.empty() {
		return nil, ErrNoCardsAvailable
	}
	return deck, nil
}

// Count returns how many cards are in a deck
func (d *CardDeck) Count() int {
	return len(d.cards)
}

// Draw gets a card
func (d *CardDeck) Draw() (Card, error) {
	if d.empty() {
		return Card{}, ErrNoCardsAvailable
	}
	return d.cards[rand.IntN(len(d.cards))], nil
}

func (d *CardDeck) empty() bool {
	return len(d.cards) == 0
}

// cardBackName is the deck's own card-back image; it isn't a drawable card.
const cardBackName = "card-back"

func (d *CardDeck) populate() error {
	glob := filepath.Join(d.Directory, "*")
	files, err := filepath.Glob(glob)
	if err != nil {
		return err
	}
	for _, f := range files {
		c, err := NewCard(f)
		if err != nil {
			continue
		}
		if c.Name() == cardBackName {
			continue
		}
		// Skip cards whose text-scan back image is missing so we never
		// render a broken second image.
		backPath := filepath.Join(d.Directory, c.Back())
		if _, err := os.Stat(backPath); err != nil {
			continue
		}
		d.cards = append(d.cards, c)
	}
	return nil
}
