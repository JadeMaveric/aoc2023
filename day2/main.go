package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type GameSet struct {
	Red   int
	Green int
	Blue  int
}

type Game struct {
	ID   int
	Sets []GameSet
}

func parseInput(filename string) ([]Game, error) {
	bytes, e := os.ReadFile(filename)
	if e != nil {
		return nil, e
	}

	chars := string(bytes)
	chars = strings.TrimSuffix(chars, "\n")

	content := strings.Split(chars, "\n")

	games := []Game{}

	for _, line := range content {
		currGame := Game{}
		gameAndBallSplit := strings.Split(line, ":")

		gameId, err := strconv.ParseInt(strings.TrimPrefix(gameAndBallSplit[0], "Game "), 10, 64)
		if err != nil {
			return nil, e
		}

		currGame.ID = int(gameId)

		gameSetsSplit := strings.Split(gameAndBallSplit[1], ";")
		for _, gameSets := range gameSetsSplit {
			gameSetsSplit := strings.Split(gameSets, ",")
			gameSet := GameSet{}

			for _, gameSetStr := range gameSetsSplit {
				gameSetStr := strings.TrimSpace(gameSetStr)

				ballAndColorSplit := strings.Split(gameSetStr, " ")
				color := ballAndColorSplit[1]
				count, err := strconv.ParseInt(ballAndColorSplit[0], 10, 64)
				if err != nil {
					return nil, e
				}

				if color == "green" {
					gameSet.Green = int(count)
				} else if color == "blue" {
					gameSet.Blue = int(count)
				} else if color == "red" {
					gameSet.Red = int(count)
				}

				// fmt.Println(ballAndColorSplit)
			}
			currGame.Sets = append(currGame.Sets, gameSet)
		}
		// fmt.Println("--")

		games = append(games, currGame)
	}

	return games, nil
}

func isGameValid(game Game, rule GameSet) bool {
	// 1. A game is invalid if the total num of balls in a given
	//    set exceeds the total num of balls in play
	// 2. A game is invalid if the num of balls of a color in a given
	//    set exceeds the total num of balls of that color in play

	for _, set := range game.Sets {
		// Check for each color
		if set.Blue > rule.Blue {
			return false
		}

		if set.Red > rule.Red {
			return false
		}

		if set.Green > rule.Green {
			return false
		}

		// Check total
		setTotal := set.Red + set.Green + set.Blue
		ruleTotal := rule.Red + rule.Green + rule.Blue
		if setTotal > ruleTotal {
			return false
		}
	}

	// If none of the conditions have broken, then the game is valid
	return true
}

func solveA(games []Game) int {
	// Which games would be possible if the bag
	// has been loaded with 12 Red, 13 Green, 14 Blue (total 39 cubes)
	acc := 0
	cond := GameSet{
		Red:   12,
		Green: 13,
		Blue:  14,
	}

	for _, game := range games {
		if isGameValid(game, cond) {
			fmt.Printf("Game %d: Valid\n", game.ID)
			acc += game.ID
		} else {
      fmt.Printf("Game %d: Invalid\n", game.ID)
    }
	}

	return acc
}

func calcMinSet(game Game) GameSet {
  minSet := GameSet{}

  for _, set := range game.Sets {
    if set.Blue > minSet.Blue {
      minSet.Blue = set.Blue
    }

    if set.Red > minSet.Red {
      minSet.Red = set.Red
    }

    if set.Green > minSet.Green {
      minSet.Green = set.Green
    }
  }

  return minSet
}

func calcSetPower(set GameSet) int {
  return set.Blue * set.Green * set.Red
}

func solveB(games []Game) int {
  // - Find the sum of the power of the min set of each game
  // - The min set of a game is the min num of cubes of each color required
  //   to play that game
  // - the power of a set is blueCount * redCount * greenCount

  acc := 0

  for _, game := range games {
    minSet := calcMinSet(game)
    power := calcSetPower(minSet)

    acc += power
  }

  return acc
}

func main() {
	games, _ := parseInput("b.txt")

  // score := solveA(games)
  // fmt.Println("Score:", score)

  power := solveB(games)
  fmt.Println("Power:", power)
}
