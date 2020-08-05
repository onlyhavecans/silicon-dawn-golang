package lib

import (
	"errors"
	"testing"
)

// todo: replace with table tests
func TestNewCard_old(t *testing.T) {
	// Happy Path
	arg := "what/a/good/card.jpg"
	got, err := NewCard(arg)
	if err != nil {
		t.Errorf("NewCard(%q) err = %v", arg, err)
	}
	wantFront := "card.jpg"
	if gotFront := got.Front(); gotFront != wantFront {
		t.Errorf("NewCard(%q).Front() = %q; want %q", arg, gotFront, wantFront)
	}
	wantBack := "card-text.png"
	if gotBack := got.Back(); gotBack != wantBack {
		t.Errorf("Expected Back: card-text.png, Got: %s", gotBack)
	}

	got, err = NewCard("just-card.jpg")
	if err != nil {
		t.Errorf("Good Card Expected: nil err, got: %v", err)
	}
	if f := got.Front(); f != "just-card.jpg" {
		t.Errorf("Expected Front: just-card.jpg, Got: %s", f)
	}
	if b := got.Back(); b != "just-card-text.png" {
		t.Errorf("Expected Back: just-card-text.png, Got: %s", b)
	}
}

func TestNewCard_discard(t *testing.T) {
	arg := "invalid-card.png"
	_, err := NewCard(arg)
	if !errors.Is(err, BadCardError) {
		t.Errorf("NewCard(%q) err = %v; want %v", arg, err, BadCardError)
	}
}

func TestNewCardDeck(t *testing.T) {
	got, err := NewCardDeck("../data")
	if err != nil {
		t.Errorf("Expected: nil error, Got: %v", err)
	}

	if c := len(got.cards); c == 0 {
		t.Error("We didn't get any cards!")
	}
}

func TestCardDeck_Draw(t *testing.T) {
	want, _ := NewCard("c.jpg")
	d := CardDeck{
		Directory: "",
		cards:     []Card{want},
	}
	got, err := d.Draw()
	if err != nil {
		t.Errorf("CardDeck.Draw() err = %v", err)
	}
	if got != want {
		t.Errorf("CardDeck.Draw() = %v; want %v", got, want)
	}
}
