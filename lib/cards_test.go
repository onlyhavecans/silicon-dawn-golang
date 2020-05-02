package lib

import "testing"

func TestNewCard(t *testing.T) {
	// Happy Path
	got, err := NewCard("what/a/good/card.jpg")
	if err != nil {
		t.Errorf("Good Card Expected: nil err, got: %w", err)
	}
	if f := got.Front(); f != "card.jpg" {
		t.Errorf("Expected Front: card.jpg, Got: %s", f)
	}
	if b := got.Back(); b != "card-text.png" {
		t.Errorf("Expected Back: card-text.png, Got: %s", b)
	}

	got, err = NewCard("just-card.jpg")
	if err != nil {
		t.Errorf("Good Card Expected: nil err, got: %w", err)
	}
	if f := got.Front(); f != "just-card.jpg" {
		t.Errorf("Expected Front: just-card.jpg, Got: %s", f)
	}
	if b := got.Back(); b != "just-card-text.png" {
		t.Errorf("Expected Back: just-card-text.png, Got: %s", b)
	}
}

func TestNewCardDiscard(t *testing.T) {
	// Bad Card
	_, err := NewCard("bad-text.png")
	if err == nil {
		t.Error("Bad Card: Expected err Got: nil")
	}
}

func TestNewCardDeck(t *testing.T) {
	got, err := NewCardDeck("../data")
	if err != nil {
		t.Errorf("Expected: nil error, Got: %w", err)
	}

	if c := len(got.cards); c == 0 {
		t.Error("We didn't get any cards!")
	}
}

func TestCardDeck_Draw(t *testing.T) {
	card, _ := NewCard("c.jpg")
	d := CardDeck{
		Directory: "test_draw",
		cards:     []Card{card},
	}
	got, err := d.Draw()
	if err != nil {
		t.Errorf("Expected: nil err, Got: %w", err)
	}
	if got != card {
		t.Errorf("Expected: Text Card, Got: %s", got)
	}
}
