package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
)

type Card struct {
	winningNums []int
	presentNums []int
}

func parseInput(filename string) ([]Card, error) {
	buff, e := os.ReadFile(filename)
	if e != nil {
		return nil, e
	}

	buff = bytes.TrimSpace(buff)

	content := bytes.Split(buff, []byte("\n"))

	cards := make([]Card, len(content))

	for i, line := range content {
		line := bytes.Split(line, []byte(":"))
		lists := bytes.Split(line[1], []byte("|"))

		left := bytes.TrimSpace(lists[0])
		winningNums := make([]int, 0)
		for _, digits := range bytes.Split(left, []byte(" ")) {
			num, _ := strconv.Atoi(string(digits))
			winningNums = append(winningNums, num)
		}

		right := bytes.TrimSpace(lists[1])
		presentNums := make([]int, 0)
		for _, digits := range bytes.Split(right, []byte(" ")) {
			num, _ := strconv.Atoi(string(digits))
			if num > 0 {
				presentNums = append(presentNums, num)
			}
		}

		cards[i] = Card{winningNums: winningNums, presentNums: presentNums}
	}

	return cards, nil
}

func getWinningNums(card Card) []int {
	winningSet := make(map[int]bool)
	for _, winningNum := range card.winningNums {
		winningSet[winningNum] = true
	}

  winningNums := make([]int, 0)

	for _, presentNum := range card.presentNums {
		if _, ok := winningSet[presentNum]; ok {
      winningNums = append(winningNums, presentNum)
		}
	}

  return winningNums
}

func getCardScore(matchCount int) int {
  return int(math.Pow(2, float64(matchCount-1)))
}

func solveA(cards []Card) int {
	var acc int = 0

	for _, card := range cards {
    winningNums := getWinningNums(card)
    score := getCardScore(len(winningNums))

		acc += score
	}

	return acc
}

func solveB(cards []Card) int {
	counts := make([]int, len(cards))

  for i := 0; i < len(cards); i++ {
    counts[i] = 1
  }

	for i, card := range cards {
    winningNums := getWinningNums(card)
    score := len(winningNums)

    for j := 0; j < score && j < len(cards); j++ {
      counts[i+j+1] += 1 * counts[i]
    }

    fmt.Println("Round ", i, ":", counts, "Score: ", score)
	}

	acc := 0
	for _, count := range counts {
		acc += count
	}
	return acc
}

func main() {
	cards, _ := parseInput("in.txt")

	solnA := solveA(cards)
	fmt.Println("Solution A: ", solnA)

  solnB := solveB(cards)
  fmt.Println("Solution B: ", solnB)
}
