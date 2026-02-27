package main

const sampleConnString = "localhost:12312"

var sampleGames = map[int]Lobby{
	1: {
		ID: 1,
		Cards: []Card{
			{
				Value: 1,
				Palo:  "basto",
			},
			{
				Value: 2,
				Palo:  "copa",
			},
			{
				Value: 3,
				Palo:  "oro",
			},
			{
				Value: 4,
				Palo:  "espada",
			},
		},
		Result: 10,
	},
	2: {
		ID: 2,
		Cards: []Card{
			{
				Value: 4,
				Palo:  "basto",
			},
			{
				Value: 3,
				Palo:  "copa",
			},
			{
				Value: 2,
				Palo:  "oro",
			},
			{
				Value: 1,
				Palo:  "espada",
			},
		},
		Result: 20,
	},
}
