package lib

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
)

type Card struct {
	name string
}

func NewCard(file string) (Card, error) {
	if strings.HasSuffix(file, ".jpg") {
		baseName := filepath.Base(file)
		cardName := baseName[0 : len(baseName)-4]
		return Card{name: cardName}, nil
	}
	return Card{}, fmt.Errorf("Bad Card: Discarding %s", file)
}

func (c *Card) Front() string {
	return c.name + ".jpg"
}

func (c *Card) Back() string {
	return c.name + "-text.png"
}

func (c *Card) Name() string {
	return c.name
}

type CardDeck struct {
	Directory string
	cards     []Card
}

func NewCardDeck(dir string) (*CardDeck, error) {
	deck := &CardDeck{Directory: dir}
	rand.Seed(time.Now().UnixNano())
	if err := deck.populate(); err != nil {
		return &CardDeck{}, err
	}
	if deck.empty() {
		return &CardDeck{}, fmt.Errorf("NoCardsAvailable")
	}
	return deck, nil
}

func (d *CardDeck) Count() int {
	return len(d.cards)
}

func (d *CardDeck) Draw() (Card, error) {
	if d.empty() {
		return Card{}, fmt.Errorf("NoCardsAvailable")
	}
	drawNum := rand.Intn(len(d.cards))
	c := d.cards[drawNum]
	return c, nil
}

func (d *CardDeck) empty() bool {
	return len(d.cards) == 0
}

func (d *CardDeck) populate() error {
	glob := filepath.Join(d.Directory, "*")
	files, err := filepath.Glob(glob)
	if err != nil {
		return err
	}
	for _, f := range files {
		if c, err := NewCard(f); err == nil {
			d.cards = append(d.cards, c)
		}
	}
	return nil
}
