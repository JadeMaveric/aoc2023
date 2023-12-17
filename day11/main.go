package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
)

type Grid [][]rune

type Coord struct{ x, y int }

func parseInput(filename string) Grid {
	buff, _ := os.ReadFile(filename)
	buff = bytes.TrimSpace(buff)

	lines := bytes.Split(buff, []byte("\n"))
	grid := make(Grid, len(lines))

	for i, line := range lines {
		grid[i] = bytes.Runes(line)
	}

	return grid
}

func isEmptyRow(row []rune) bool {
	for _, char := range row {
		if char != '.' {
			return false
		}
	}
	return true
}

func isEmptyCol(grid Grid, col int) bool {
	for _, row := range grid {
		char := row[col]
		if char != '.' {
			return false
		}
	}
	return true
}

func expandGrid(grid Grid) Grid {
	// Create the new grid - expand horizontally first
	ngH := make(Grid, len(grid))

	hExpansion := 0

	for x := range grid[0] {
		for y := range grid {
			ngH[y] = append(ngH[y], grid[y][x])
		}

		if isEmptyCol(grid, x) {
			hExpansion++
			for y := range grid {
				ngH[y] = append(ngH[y], grid[y][x])
			}
		}
	}

	// Expand vertically
	ngV := make(Grid, 0)
	vExpansion := 0
	for _, row := range ngH {
		ngV = append(ngV, row)
		if isEmptyRow(row) {
			ngV = append(ngV, row)
			vExpansion++
		}
	}

	fmt.Println("h:", vExpansion, "v:", hExpansion)

	return ngV
}

func findGalaxies(grid Grid) (galaxies []Coord) {
	galaxies = make([]Coord, 0)

	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == '#' {
				galaxies = append(galaxies, Coord{x, y})
			}
		}
	}

	return
}

func solveA(grid Grid) (solnA float64) {
	grid = expandGrid(grid)
	galaxies := findGalaxies(grid)

	fmt.Println("Processing", len(galaxies), "galaxies")

	for _, g1 := range galaxies {
		for _, g2 := range galaxies {
      dist := math.Abs(float64(g1.x-g2.x)) + math.Abs(float64(g1.y-g2.y))
      // fmt.Println(g1, g2, dist)
			solnA += dist
		}
	}

	return solnA / 2
}

func getEmptyRowCount(grid Grid, start, end int) (count int) {
	for y := start; y < end; y++ {
		if isEmptyRow(grid[y]) {
			count++
		}
	}
	return
}

func getEmptyColCount(grid Grid, start, end int) (count int) {
	for x := start; x < end; x++ {
		if isEmptyCol(grid, x) {
			count++
		}
	}
	return
}

func solveB(grid Grid) (solnB float64) {
  // B is similar to A, just that the expansion is 1_000_000 instead of 2
  // Since we're using taxi-cab distance, we can directly add the added distance
  // no need to actually expand the grid
	galaxies := findGalaxies(grid)

	fmt.Println("Processing", len(galaxies), "galaxies")

	for _, g1 := range galaxies {
		for _, g2 := range galaxies {
      dist := math.Abs(float64(g1.x-g2.x)) + math.Abs(float64(g1.y-g2.y))
			rowCount := (1_000_000-1) * getEmptyRowCount(grid, min(g1.y, g2.y), max(g1.y, g2.y))
			colCount := (1_000_000-1) * getEmptyColCount(grid, min(g1.x, g2.x), max(g1.x, g2.x))

			dist += float64(rowCount + colCount)
      // fmt.Println(g1, g2, dist, fmt.Sprintf("+(%d, %d)", colCount, rowCount))
      solnB += dist
		}
	}

	return solnB / 2
}

func main() {
	grid := parseInput("day11/in.txt")

	if len(grid) < 80 {
		for y := range grid {
			for x := range grid[y] {
				fmt.Printf("%c", grid[y][x])
			}
			fmt.Println()
		}
	}

  // for i := 0; i < 2; i++ {
  //   fmt.Printf("Col %d is %v\n", i, isEmptyCol(grid, i))
  // }

	solnA := solveA(grid)
	fmt.Println("A:", int(solnA))

	solnB := solveB(grid)
	fmt.Println("B:", int(solnB))
}
