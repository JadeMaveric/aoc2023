package main

import (
	"bytes"
	"fmt"
	"os"
	"unicode"
)

const COLOR_RED = "\033[0;31m"
const COLOR_GOLD = "\033[0;33m"
const COLOR_NONE = "\033[0m"

type Grid [][]rune

type Box struct {
	x1, y1, x2, y2 int
}

type Point struct {
	x, y int
}

func parseInput(filename string) (Grid, error) {
	buff, e := os.ReadFile(filename)
	if e != nil {
		return nil, e
	}

	buff = bytes.TrimSpace(buff)

	content := bytes.Split(buff, []byte("\n"))

	grid := make(Grid, len(content))

	for i := range content {
		grid[i] = bytes.Runes(content[i])
	}

	return grid, e
}

func isNextToSymbol(grid Grid, box Box) bool {
	// Check vertical axes
	for _, col := range []int{box.x1 - 1, box.x2 + 1} {
		for row := box.y1 - 1; row <= box.y2+1; row++ {
			if row < 0 || row > (len(grid)-1) {
				continue
			}

			if col < 0 || col > (len(grid[0])-1) {
				continue
			}

			if grid[row][col] != '.' && !unicode.IsDigit(grid[row][col]) {
				return true
				// fmt.Printf("1. %s (%d, %d) - {%s} %s\n", COLOR_GOLD, row, col, string(grid[row][col]), COLOR_NONE)
			} else {
				// fmt.Printf("1. (%d, %d) - {%s}\n", row, col, string(grid[row][col]))
			}
		}
	}

	// Check horizontal axes
	for _, row := range []int{box.y1 - 1, box.y2 + 1} {
		for col := box.x1 - 1; col <= box.x2+1; col++ {
			if row < 0 || row > (len(grid)-1) {
				continue
			}

			if col < 0 || col > (len(grid[0])-1) {
				continue
			}

			if grid[row][col] != '.' && !unicode.IsDigit(grid[row][col]) {
				return true
				// fmt.Printf("2. %s (%d, %d) - {%s} %s\n", COLOR_GOLD, row, col, string(grid[row][col]), COLOR_NONE)
			} else {
				// fmt.Printf("2. (%d, %d) - {%s}\n", row, col, string(grid[row][col]))
			}
		}
	}

	return false
}

func solveA(grid Grid) int {
	acc := 0

	// Iterate through every cell
	for row := 0; row < len(grid); row++ {
		// Reset state for each row
		foundNum := false
		currNum := 0
		nextToSymbol := false

		for col := 0; col < len(grid[row]); col++ {
			cell := grid[row][col]
			isDigit := unicode.IsDigit(cell)

			// 1. Still searching for a digit
			if !foundNum && !isDigit {
				continue
			}

			// 2. Currently in a digit
			if isDigit {
				foundNum = true
				currNum = 10*currNum + int(cell-'0')

				box := Box{x1: col, y1: row, x2: col, y2: row}
				nextToSymbol = nextToSymbol || isNextToSymbol(grid, box)
			}

			// 3. Digit got over
			isLastChar := col == len(grid[row])-1
			if foundNum && (!isDigit || isLastChar) {
				if nextToSymbol {
					acc += currNum
				}

				// if nextToSymbol {
				// 	fmt.Printf("%d, YES\n", currNum)
				// } else {
				// 	fmt.Printf(("%d, NO\n"), currNum)
				// }

				// Reset state for next digit
				foundNum = false
				currNum = 0
				nextToSymbol = false
			}
		}
	}

	return acc
}

func removeDuplicates(arr []Point) []Point {
  seen := make(map[Point]bool)

  for _, point := range arr {
    if _, ok := seen[point]; !ok {
      seen[point] = true
    }
  }

  result := make([]Point, 0)

  for point := range seen {
    result = append(result, point)
  }

  return result
}

func findNearbyGears(grid Grid, point Point) []Point {
  // fmt.Printf("%c (%d, %d)\n", grid[point.y][point.x], point.y, point.x)
  gears := make([]Point, 0)

	// Check vertical axes
	for _, dx := range []int{-1, 0, 1} {
		for _, dy := range []int{-1, 0, 1} {
      if dx == 0 && dy == 0 {
        continue
      }

      y := point.y + dy
      x := point.x + dx

			if y < 0 || y > (len(grid)-1) {
				continue
			}

			if x < 0 || x > (len(grid[0])-1) {
				continue
			}

			if grid[y][x] == '*' {
        gears = append(gears, Point{x: x, y: y})
				// fmt.Printf("%s+(%d, %d)\t=(%d, %d)\t- {%s} %s\n", COLOR_GOLD, dy, dx, y, x, string(grid[y][x]), COLOR_NONE)
			} else {
				// fmt.Printf("+(%d, %d)\t=(%d, %d)\t- {%s}\n", dy, dx, y, x, string(grid[y][x]))
			}
		}
	}

	return gears
}

func solveB(grid Grid) int {
  engineParts := make(map[Point][]int)

	// Iterate through every cell
	for y := 0; y < len(grid); y++ {
		// Reset state for each row
		foundNum := false
		currNum := 0
    nearbyGears := make([]Point, 0)

		for x := 0; x < len(grid[y]); x++ {
			cell := grid[y][x]
			isDigit := unicode.IsDigit(cell)

			// 1. Still searching for a digit
			if !foundNum && !isDigit {
				continue
			}

			// 2. Currently in a digit
			if isDigit {
				foundNum = true
				currNum = 10*currNum + int(cell-'0')

        nearbyGears = append(nearbyGears, findNearbyGears(grid, Point{x: x, y: y})...)
			}

			// 3. Digit got over
			isLastChar := x == len(grid[y])-1
			if foundNum && (!isDigit || isLastChar) {
        nearbyGears = removeDuplicates(nearbyGears)

        for _, gear := range nearbyGears {
          engineParts[gear] = append(engineParts[gear], currNum)
        }

				// if len(nearbyGears) > 0 {
				// 	fmt.Printf("%d, YES\n", currNum)
				// } else {
				// 	fmt.Printf(("%d, NO\n"), currNum)
				// }

				// Reset state for next digit
				foundNum = false
				currNum = 0
        nearbyGears = make([]Point, 0)
			}
		}
	}

  // fmt.Println("engineParts:", engineParts)

  acc := 0

  for _, parts := range engineParts {
    if len(parts) == 2 {
      gearRatio := parts[0] * parts[1]

      // fmt.Printf("gearRatio for [%d, %d] = %d\n", parts[0], parts[1], gearRatio)

      acc += gearRatio
    }
  }

  return acc
}

func main() {
	grid, _ := parseInput("in.txt")
	// solution := solveA(grid)
	// fmt.Printf("Solution A: %d\n", solution)

  solution := solveB(grid)
  fmt.Printf("Solution B: %d\n", solution)
  // for row := 0; row < len(grid); row++ {
  //   for col := 0; col < len(grid[row]); col++ {
  //     if grid[row][col] == '*' {
  //       fmt.Printf("%s%c ", COLOR_GOLD, grid[row][col])
  //     } else if len(findNearbyGears(grid, Point{x: col, y: row})) > 0  {
  //       fmt.Printf("%s%c ", COLOR_RED, grid[row][col])
  //     } else {
  //       fmt.Printf("%s%c ", COLOR_NONE, grid[row][col])
  //     }
  //     // fmt.Println()
  //   }
  //   fmt.Println()
  // }
}
