package kriptogame

import (
	"math/rand"
)

type Palo int

const (
	Basto Palo = iota
	Espada
	Copa
	Oro
)

type Card struct {
	Value int
	Palo  Palo
}

func generateHand() []Card {
	cards := make([]Card, 0, 40)
	for i := 1; i <= 12; i++ {
		if i == 8 || i == 9 {
			continue
		}
		cards = append(cards, Card{
			Value: i,
			Palo:  Basto,
		})
	}
	for i := 1; i <= 12; i++ {
		if i == 8 || i == 9 {
			continue
		}
		cards = append(cards, Card{
			Value: i,
			Palo:  Espada,
		})
	}
	for i := 1; i <= 12; i++ {
		if i == 8 || i == 9 {
			continue
		}
		cards = append(cards, Card{
			Value: i,
			Palo:  Oro,
		})
	}
	for i := 1; i <= 12; i++ {
		if i == 8 || i == 9 {
			continue
		}
		cards = append(cards, Card{
			Value: i,
			Palo:  Copa,
		})
	}
	rand.Shuffle(len(cards), func(i, j int) {
		temp := cards[i]
		cards[i] = cards[j]
		cards[j] = temp
	})

	return cards[:5]
}
