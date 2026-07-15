package cards

import (
	"errors"
	"testing"
)

func TestNewCard(t *testing.T) {
	tests := []struct {
		name      string
		file      string
		wantFront string
		wantBack  string
		wantErr   error
	}{
		{
			name:      "nested path",
			file:      "what/a/good/card.jpg",
			wantFront: "card.jpg",
			wantBack:  "card-text.png",
		},
		{
			name:      "bare filename",
			file:      "just-card.jpg",
			wantFront: "just-card.jpg",
			wantBack:  "just-card-text.png",
		},
		{
			name:    "wrong extension is rejected",
			file:    "invalid-card.png",
			wantErr: ErrBadCard,
		},
		{
			name:    "uppercase extension is rejected",
			file:    "invalid-card.JPG",
			wantErr: ErrBadCard,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCard(tt.file)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("NewCard(%q) err = %v; want %v", tt.file, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("NewCard(%q) err = %v; want nil", tt.file, err)
			}
			if gotFront := got.Front(); gotFront != tt.wantFront {
				t.Errorf("NewCard(%q).Front() = %q; want %q", tt.file, gotFront, tt.wantFront)
			}
			if gotBack := got.Back(); gotBack != tt.wantBack {
				t.Errorf("NewCard(%q).Back() = %q; want %q", tt.file, gotBack, tt.wantBack)
			}
		})
	}
}

func TestNewCardDeck(t *testing.T) {
	got, err := NewCardDeck("testdata")
	if err != nil {
		t.Fatalf("NewCardDeck() err = %v; want nil", err)
	}

	if c := got.Count(); c == 0 {
		t.Error("NewCardDeck() produced an empty deck")
	}
}

func TestNewCardDeck_missingBackIsSkipped(t *testing.T) {
	// testdata/missing-back has a front image with no matching -text.png;
	// populate() should silently skip it rather than produce a deck that
	// draws a card with a broken back image.
	got, err := NewCardDeck("testdata/missing-back")
	if !errors.Is(err, ErrNoCardsAvailable) {
		t.Fatalf("NewCardDeck() err = %v; want %v", err, ErrNoCardsAvailable)
	}
	if got != nil {
		t.Fatalf("NewCardDeck() deck = %v; want nil on error", got)
	}
}

func TestNewCardDeck_badDir(t *testing.T) {
	got, err := NewCardDeck("testdata/does-not-exist")
	if !errors.Is(err, ErrNoCardsAvailable) {
		t.Fatalf("NewCardDeck() err = %v; want %v", err, ErrNoCardsAvailable)
	}
	if got != nil {
		t.Fatalf("NewCardDeck() deck = %v; want nil on error", got)
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

func TestCardDeck_Draw_empty(t *testing.T) {
	d := CardDeck{}
	if _, err := d.Draw(); !errors.Is(err, ErrNoCardsAvailable) {
		t.Errorf("CardDeck.Draw() err = %v; want %v", err, ErrNoCardsAvailable)
	}
}
