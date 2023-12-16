package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type Play struct {
	hand     []rune
	bid      int
	handType int
}

const (
	HandUnknown   = iota
	HandHighCard  = iota
	HandOnePair   = iota
	HandTwoPair   = iota
	HandThreeKind = iota
	HandFullhouse = iota
	HandFourKind  = iota
	HandFiveKind  = iota
)

const (
	CardWild  = iota
	CardTwo   = iota
	CardThree = iota
	CardFour  = iota
	CardFive  = iota
	CardSix   = iota
	CardSeven = iota
	CardEight = iota
	CardNine  = iota
	CardTen   = iota
	CardJack  = iota
	CardQueen = iota
	CardKing  = iota
	CardAce   = iota
	// CardWild  = iota
)

func getCardValue(r rune, isPartTwo bool) int {
	switch r {
	case '2':
		return CardTwo
	case '3':
		return CardThree
	case '4':
		return CardFour
	case '5':
		return CardFive
	case '6':
		return CardSix
	case '7':
		return CardSeven
	case '8':
		return CardEight
	case '9':
		return CardNine
	case 'T':
		return CardTen
	case 'J':
		if isPartTwo {
			return CardWild
		} else {
			return CardJack
		}
	case 'Q':
		return CardQueen
	case 'K':
		return CardKing
	case 'A':
		return CardAce
	default:
		return -1
	}
}

func parseInput(filename string) []Play {
	buff, _ := os.ReadFile(filename)
	buff = bytes.TrimSpace(buff)
	lines := bytes.Split(buff, []byte("\n"))

	plays := make([]Play, len(lines))
	for i, line := range lines {
		parts := bytes.Fields(line)
		hand := bytes.Runes(parts[0])
		bid, _ := strconv.Atoi(string(parts[1]))
		plays[i] = Play{hand, bid, HandUnknown}
	}

	return plays
}

func getHandType(hand []rune, isPartTwo bool) int {
	counts := make(map[rune]int)
	wildCount := 0

	for _, r := range hand {
		if r == 'J' && isPartTwo {
			wildCount++
		} else {
			counts[r]++
		}
	}

	// Add wilds to the card with the most count
	maxCard := 'A'
	maxCardCount := counts['A']
	for card, count := range counts {
		if count > maxCardCount {
			maxCard = card
			maxCardCount = count
		}
	}
  counts[maxCard] += wildCount

	haveFive := false
	haveFour := false
	haveThree := false
	havePair := false
	pairCount := 0

	for _, count := range counts {
		if count == 5 {
			haveFive = true
		} else if count == 4 {
			haveFour = true
		} else if count == 3 {
			haveThree = true
		} else if count == 2 {
			havePair = true
			pairCount++
		}
	}

	if haveFive {
		return HandFiveKind
	} else if haveFour {
		return HandFourKind
	} else if haveThree && havePair {
		return HandFullhouse
	} else if haveThree && !havePair {
		return HandThreeKind
	} else if havePair && pairCount == 2 {
		return HandTwoPair
	} else if havePair && pairCount == 1 {
		return HandOnePair
	} else {
		return HandHighCard
	}
}

func compareHands(a, b []rune, isPartTwo bool) int {
	length := min(len(a), len(b))

	for i := 0; i < length; i++ {
		valA := getCardValue(a[i], isPartTwo)
		valB := getCardValue(b[i], isPartTwo)
		if valA == valB {
			continue
		} else {
			return valA - valB
		}
	}

	return 0
}

func solveA(plays []Play, isPartTwo bool) int {
	slices.SortFunc(plays, func(a, b Play) int {
		handA := getHandType(a.hand, isPartTwo)
		handB := getHandType(b.hand, isPartTwo)

		if handA != handB {
			return handA - handB
		} else {
			return compareHands(a.hand, b.hand, isPartTwo)
		}
	})

	acc := 0

	fmt.Println("Hand\tBid\tType\tRank")
	for i, p := range plays {
		fmt.Printf("%s\t%d\t%d\t%d\n", string(p.hand), p.bid, getHandType(p.hand, isPartTwo), i+1)
		acc += (i + 1) * p.bid
	}

	return acc
}

func main() {
	plays := parseInput("day7/in.txt")

	solnA := solveA(plays, false)
	fmt.Println("Solution A:", solnA)

	solnB := solveA(plays, true)
	fmt.Println("Solution B:", solnB)
}
