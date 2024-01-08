package main

import (
	"bytes"
	"fmt"
	"os"
)

func parseInput(filename string) (grid Grid) {
	buff, _ := os.ReadFile(filename)
	buff = bytes.Trim(buff, "\n")
	lines := bytes.Split(buff, []byte("\n"))

	grid = make(Grid, len(lines))

	for y := range lines {
		grid[y] = make([](*Block), len(lines[y]))

		for x := range lines[y] {
			grid[y][x] = &Block{
				heatLoss: int(lines[y][x] - '0'),
				posX:     x,
				posY:     y,
			}
		}
	}

	return grid
}

func main() {
	grid := parseInput("day17/a.txt")
	fmt.Println(grid.String() + "\n")

	start := grid[0][0]
	for i := range start.lossGrid {
		start.lossGrid[i] = &Step{cameFrom: start, bestLoss: -1, xStep: 0, yStep: 0}
	}

	end := grid[len(grid)-1][len(grid[0])-1]

	// path := AStarPathFinderSimple{}.FindPath(&grid, start, end)
	path := AStarPathFinderUltra{}.FindPath(&grid, start, end)

	fmt.Println(grid.VisitedString(path))
}
