package main

import (
	"bytes"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"os"
	"time"
)

type Platform [][]byte
type Rock byte

type Coord struct{ x, y int }

const (
	RRound  Rock = 'O'
	RSquare Rock = '#'
	RTile   Rock = '.'
)

func parseInput(filename string) Platform {
	buff, _ := os.ReadFile(filename)
	buff = bytes.TrimSpace(buff)
	platform := bytes.Split(buff, []byte("\n"))
	return platform
}

func (p Platform) String() string {
	var buff bytes.Buffer
	for _, row := range p {
		buff.Write(row)
		buff.WriteByte('\n')
	}
	return buff.String()
}

func getHighestCoord(platform Platform, coord Coord) Coord {
	x := coord.x

	for y := coord.y; y > 0; y-- {
		destTile := Rock(platform[y-1][x])

		if destTile == RSquare || destTile == RRound {
			return Coord{x, y}
		}
	}

	return Coord{x, 0}
}

func getLowestCoord(platform Platform, coord Coord) Coord {
	x := coord.x
	for y := coord.y; y < len(platform)-1; y++ {
		destTile := Rock(platform[y+1][x])
		if destTile == RSquare || destTile == RRound {
			return Coord{x, y}
		}
	}
	return Coord{x, len(platform) - 1}
}

func getLeftmostCoord(platform Platform, coord Coord) Coord {
	y := coord.y
	for x := coord.x; x > 0; x-- {
		destTile := Rock(platform[y][x-1])
		if destTile == RSquare || destTile == RRound {
			return Coord{x, y}
		}
	}
	return Coord{0, y}
}

func getRightmostCoord(platform Platform, coord Coord) Coord {
	y := coord.y
	for x := coord.x; x < len(platform[0])-1; x++ {
		destTile := Rock(platform[y][x+1])
		if destTile == RSquare || destTile == RRound {
			return Coord{x, y}
		}
	}
	return Coord{len(platform[0]) - 1, y}
}

func moveRocksUp(platform Platform) Platform {
	// Start from the top-left
	// If we find a round rock, move it as high up as possible
	// To find the how high up, start scaning upwards from the curr position
	// Stop once we reach another rock or the highest row

	for x := range platform[0] {
		for y := range platform {
			if Rock(platform[y][x]) == RRound {
				currPos := Coord{x, y}
				nextPos := getHighestCoord(platform, currPos)

				platform[currPos.y][currPos.x] = byte(RTile)
				platform[nextPos.y][nextPos.x] = byte(RRound)
			}
		}
	}

	return platform
}

func moveRocksDown(platform Platform) Platform {
	// Start from the bottom-left
	// If we find a round rock, move it as low down as possible
	// To find the how low down, start scaning downwards from the curr position
	// Stop once we reach another rock or the lowest row
	for x := range platform[0] {
		for y := len(platform) - 1; y >= 0; y-- {
			if Rock(platform[y][x]) == RRound {
				currPos := Coord{x, y}
				nextPos := getLowestCoord(platform, currPos)
				platform[currPos.y][currPos.x] = byte(RTile)
				platform[nextPos.y][nextPos.x] = byte(RRound)
			}
		}
	}
	return platform
}

func moveRocksLeft(platform Platform) Platform {
	// Start from the bottom-left
	// If we find a round rock, move it as left as possible
	// To find the how left, start scaning leftwards from the curr position
	// Stop once we reach another rock or the leftmost col
	for y := range platform {
		for x := range platform[y] {
			if Rock(platform[y][x]) == RRound {
				currPos := Coord{x, y}
				nextPos := getLeftmostCoord(platform, currPos)
				platform[currPos.y][currPos.x] = byte(RTile)
				platform[nextPos.y][nextPos.x] = byte(RRound)
			}
		}
	}
	return platform
}

func moveRocksRight(platform Platform) Platform {
	// Start from the bottom-right
	// If we find a round rock, move it as right as possible
	// To find the how right, start scaning rightwards from the curr position
	// Stop once we reach another rock or the rightmost col
	for y := range platform {
		for x := len(platform[y]) - 1; x >= 0; x-- {
			if Rock(platform[y][x]) == RRound {
				currPos := Coord{x, y}
				nextPos := getRightmostCoord(platform, currPos)
				platform[currPos.y][currPos.x] = byte(RTile)
				platform[nextPos.y][nextPos.x] = byte(RRound)
			}
		}
	}
	return platform
}

// Multi-threaded version of moveRocksUp
// Since each column is independent of the other, we can parallelize the inner loop
// This is worse than the single threaded-version
// func moveRocksUpParallel(platform Platform) Platform {
// 	var wg sync.WaitGroup

// 	for x := range platform[0] {
//     wg.Add(1)
// 		go func(x int) {
//       defer wg.Done()
// 			for y := range platform {
// 				if Rock(platform[y][x]) == RRound {
// 					currPos := Coord{x, y}
// 					nextPos := getHighestCoord(platform, currPos)

//           // No need for mutex since each goroutine is working on a different column
// 					platform[currPos.y][currPos.x] = byte(RTile)
// 					platform[nextPos.y][nextPos.x] = byte(RRound)
// 				}
// 			}
// 		}(x)
// 	}

//   wg.Wait()

// 	return platform
// }

func countScore(platform Platform) (score int) {
	for y := range platform {
		for x := range platform[y] {
			tile := Rock(platform[y][x])

			if tile == RRound {
				score += len(platform) - y
			}
		}
	}
	return score
}

func solveA(platform Platform) int {
	defer timer("solveA")()
	return countScore(moveRocksUp(platform))
}

func solveB(platform Platform) int {
	defer timer("solveB")()

	cache := make(map[string]int)

	// Pretty Progress
	LIMIT := 1_000_000_000
	bar := progressbar.NewOptions(
		LIMIT,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionThrottle(1*time.Second),
		progressbar.OptionOnCompletion(func() { fmt.Println() }),
		progressbar.OptionShowElapsedTimeOnFinish(),
	)

	for i := 0; i < LIMIT; {
		if prevI, ok := cache[platform.String()]; !ok {
			cache[platform.String()] = i
		} else {
      delta := i - prevI
			fmt.Printf("Cycle detected from %d -> %d: %d\n", prevI, i, delta)

      // We can skip the remaining cycles since the pattern repeats
      repeats := delta * int((LIMIT - i) / delta)
      bar.Add(repeats)
      i += repeats
		}
		// fmt.Printf("%s\n%s\n", "ORIG:", platform)

		platform = moveRocksUp(platform)
		// fmt.Printf("%s\n%s\n", "UP:", platform)

		platform = moveRocksLeft(platform)
		// fmt.Printf("%s\n%s\n", "LEFT:", platform)

		platform = moveRocksDown(platform)
		// fmt.Printf("%s\n%s\n", "DOWN:", platform)

		platform = moveRocksRight(platform)
		// fmt.Printf("After %d cycles:\n%s\n", i+1, platform)
		bar.Add(1)
    i++
	}
	return countScore(platform)
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Println(name, "took:", time.Since(start))
	}
}

func main() {
	platform := parseInput("day14/a.txt")
	solnA := solveA(platform)
	fmt.Println("A:", solnA)

	fmt.Println()

	platform = parseInput("day14/a.txt")
	solnB := solveB(platform)
	fmt.Println("B:", solnB)
}
