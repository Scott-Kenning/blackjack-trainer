package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Card struct {
	Value int
	Suit  string
}

type Deck struct {
	Cards []Card
}

type Hand struct {
	Cards []Card
}

type Game struct {
	PlayerHand Hand
	DealerHand Hand
	Deck       Deck
	Score      int
	Rounds     int
}

func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) { d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i] })
}

func (d *Deck) Deal() Card {
	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return card
}

func (h *Hand) Value() int {
	value := 0
	for _, card := range h.Cards {
		value += card.Value
	}
	return value
}

// Returns the correct action given a dealer and player hand
func GetCorrectAction(dealerHand Hand, playerHand Hand) string {
	dealerCardValue := dealerHand.Cards[0].Value
	playerCardValues := []int{playerHand.Cards[0].Value, playerHand.Cards[1].Value}

	// Adjust for face cards and aces
	if dealerCardValue > 10 {
		dealerCardValue = 10
	}

	// Check for pairs
	isPair := playerCardValues[0] == playerCardValues[1]
	if isPair {
		switch playerCardValues[0] {
		case 8:
			return "p"
		case 7, 6, 3, 2:
			if dealerCardValue >= 2 && dealerCardValue <= 7 {
				return "p"
			} else {
				return "h"
			}
		case 5:
			return "d"
		case 4:
			if dealerCardValue == 5 || dealerCardValue == 6 {
				return "p"
			} else {
				return "h"
			}
		default: // For aces, nines, and tens
			return "p"
		}
	}

	for i, val := range playerCardValues {
		if val > 10 {
			playerCardValues[i] = 10
		}
	}

	// Check for soft totals (ace present)
	isSoft := playerCardValues[0] == 1 || playerCardValues[1] == 1
	if isSoft {
		softTotal := playerCardValues[0] + playerCardValues[1] + 10 // Adjusting for Ace value
		if softTotal >= 19 {
			return "s"
		} else if softTotal == 18 {
			if dealerCardValue >= 3 && dealerCardValue <= 6 {
				return "d"
			} else {
				return "s"
			}
		} else if softTotal == 17 || softTotal == 16 {
			if dealerCardValue >= 4 && dealerCardValue <= 6 {
				return "d"
			} else {
				return "h"
			}
		} else if softTotal == 15 || softTotal == 14 {
			if dealerCardValue >= 4 && dealerCardValue <= 6 {
				return "d"
			} else {
				return "h"
			}
		} else {
			return "h"
		}
	}

	// Check for hard totals
	hardTotal := playerCardValues[0] + playerCardValues[1]
	if hardTotal >= 17 {
		return "s"
	} else if hardTotal >= 13 && hardTotal <= 16 {
		if dealerCardValue >= 2 && dealerCardValue <= 6 {
			return "s"
		} else {
			return "h"
		}
	} else if hardTotal == 12 {
		if dealerCardValue >= 4 && dealerCardValue <= 6 {
			return "s"
		} else {
			return "h"
		}
	} else if hardTotal == 11 {
		return "d"
	} else if hardTotal == 10 {
		if dealerCardValue >= 2 && dealerCardValue <= 9 {
			return "d"
		} else {
			return "h"
		}
	} else {
		return "h"
	}
}

// maps single letter options to action string
func ActionString(action string) string {
	switch action {
	case "h":
		return "hit"
	case "s":
		return "stand"
	case "d":
		return "double"
	case "p":
		return "split"
	default:
		return "invalid"
	}
}

// resets game state
func (g *Game) Reset() {
	g.Deck = Deck{
		Cards: []Card{
			{Value: 1, Suit: "Hearts"}, {Value: 2, Suit: "Hearts"}, {Value: 3, Suit: "Hearts"},
			{Value: 4, Suit: "Hearts"}, {Value: 5, Suit: "Hearts"}, {Value: 6, Suit: "Hearts"},
			{Value: 7, Suit: "Hearts"}, {Value: 8, Suit: "Hearts"}, {Value: 9, Suit: "Hearts"},
			{Value: 10, Suit: "Hearts"}, {Value: 11, Suit: "Hearts"}, {Value: 12, Suit: "Hearts"},
			{Value: 13, Suit: "Hearts"},

			{Value: 1, Suit: "Diamonds"}, {Value: 2, Suit: "Diamonds"}, {Value: 3, Suit: "Diamonds"},
			{Value: 4, Suit: "Diamonds"}, {Value: 5, Suit: "Diamonds"}, {Value: 6, Suit: "Diamonds"},
			{Value: 7, Suit: "Diamonds"}, {Value: 8, Suit: "Diamonds"}, {Value: 9, Suit: "Diamonds"},
			{Value: 10, Suit: "Diamonds"}, {Value: 11, Suit: "Diamonds"}, {Value: 12, Suit: "Diamonds"},
			{Value: 13, Suit: "Diamonds"},

			{Value: 1, Suit: "Clubs"}, {Value: 2, Suit: "Clubs"}, {Value: 3, Suit: "Clubs"},
			{Value: 4, Suit: "Clubs"}, {Value: 5, Suit: "Clubs"}, {Value: 6, Suit: "Clubs"},
			{Value: 7, Suit: "Clubs"}, {Value: 8, Suit: "Clubs"}, {Value: 9, Suit: "Clubs"},
			{Value: 10, Suit: "Clubs"}, {Value: 11, Suit: "Clubs"}, {Value: 12, Suit: "Clubs"},
			{Value: 13, Suit: "Clubs"},

			{Value: 1, Suit: "Spades"}, {Value: 2, Suit: "Spades"}, {Value: 3, Suit: "Spades"},
			{Value: 4, Suit: "Spades"}, {Value: 5, Suit: "Spades"}, {Value: 6, Suit: "Spades"},
			{Value: 7, Suit: "Spades"}, {Value: 8, Suit: "Spades"}, {Value: 9, Suit: "Spades"},
			{Value: 10, Suit: "Spades"}, {Value: 11, Suit: "Spades"}, {Value: 12, Suit: "Spades"},
			{Value: 13, Suit: "Spades"},
		},
	}
	g.PlayerHand = Hand{}
	g.DealerHand = Hand{}
}

// Initialize game by shuffling deck and then dealing cards
func (g *Game) Init() {
	g.Reset()
	g.Deck.Shuffle()
	g.PlayerHand = Hand{Cards: []Card{g.Deck.Deal(), g.Deck.Deal()}}
	g.DealerHand = Hand{Cards: []Card{g.Deck.Deal(), g.Deck.Deal()}}
}

// function maps card value to string, which changes 11-13 to J-K and 1 - A and adds suit
func (c *Card) String() string {
	value := c.Value
	switch value {
	case 1:
		return fmt.Sprintf("A of %v", c.Suit)
	case 11:
		return fmt.Sprintf("J of %v", c.Suit)
	case 12:
		return fmt.Sprintf("Q of %v", c.Suit)
	case 13:
		return fmt.Sprintf("K of %v", c.Suit)
	default:
		return fmt.Sprintf("%v of %v", value, c.Suit)
	}
}

func main() {
	game := Game{}

	for {
		game.Init()

		fmt.Println("Player hand:", game.PlayerHand.Cards[0].String(), ",", game.PlayerHand.Cards[1].String())
		fmt.Println("Dealer hand:", game.DealerHand.Cards[0].String())

		var action string
		fmt.Print("Action (hit (H) / stand (S) / double (D) / split (P)): ")
		fmt.Scanln(&action)

		// if action isnt valid, ask again
		for action != "h" && action != "s" && action != "d" && action != "p" && action != "q" && action != "quit" {
			fmt.Println("\nInvalid action, please try again\n\n")
			fmt.Print("Action (hit (H) / stand (S) / double (D) / split (P)): ")
			fmt.Scanln(&action)
		}

		if action == "quit" || action == "q" {
			fmt.Printf("\n\nFinal score: %v / %v, %.2v%%\n", game.Score, game.Rounds, float64(game.Score)/float64(game.Rounds)*100)
			break
		}

		var correctAction = GetCorrectAction(game.DealerHand, game.PlayerHand)

		if strings.ToLower(action) == correctAction {
			fmt.Println("Correct!\n\n")
			game.Score++
		} else {
			fmt.Println("Incorrect, correct answer was: \n", ActionString(correctAction), "\n\n")
		}
		game.Rounds++
	}

}
