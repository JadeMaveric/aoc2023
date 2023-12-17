package main

import (
	"bytes"
	"fmt"
	"os"
)

type Pipe int
type Grid [][]rune

const (
	Invalid = iota
	Ground  = iota
	Start   = iota
	PipeNS  = iota
	PipeWE  = iota
	PipeNE  = iota
	PipeNW  = iota
	PipeSE  = iota
	PipeSW  = iota
)

type Coord struct{ x, y int }

func charToPipe(char rune) Pipe {
	switch char {
	case '.':
		return Ground
	case '|':
		return PipeNS
	case '-':
		return PipeWE
	case 'L':
		return PipeNE
	case 'J':
		return PipeNW
	case 'F':
		return PipeSE
	case '7':
		return PipeSW
	case 'S':
		return Start
	default:
		return Invalid
	}
}

func pipeToChar(pipe Pipe) rune {
	switch pipe {
	case PipeNS:
		return '|'
	case PipeWE:
		return '-'
	case PipeNE:
		return 'L'
	case PipeNW:
		return 'J'
	case PipeSE:
		return 'F'
	case PipeSW:
		return '7'
	case Start:
		return 'S'
	case Ground:
		return '.'
	default:
		return '?'
	}
}

func pipeToDisplayChar(pipe Pipe) rune {
	switch pipe {
	case PipeNS:
		return '║'
	case PipeWE:
		return '═'
	case PipeNE:
		return '╚'
	case PipeNW:
		return '╝'
	case PipeSE:
		return '╔'
	case PipeSW:
		return '╗'
	case Start:
		return 'S'
	case Ground:
		return '.'
	default:
		return '?'
	}
}

func getNeighbours(curr Coord, pipe Pipe) []Coord {
	x, y := curr.x, curr.y
	// We're scanning top-down so the y-axis in inverted
	north := Coord{x, y - 1}
	south := Coord{x, y + 1}
	east := Coord{x + 1, y}
	west := Coord{x - 1, y}

	switch pipe {
	case PipeNS:
		return []Coord{north, south}
	case PipeWE:
		return []Coord{west, east}
	case PipeNE:
		return []Coord{north, east}
	case PipeNW:
		return []Coord{north, west}
	case PipeSE:
		return []Coord{south, east}
	case PipeSW:
		return []Coord{south, west}
	default:
		return []Coord{}
	}
}

func parseInput(filename string) Grid {
	buff, _ := os.ReadFile(filename)
	content := bytes.Split(bytes.TrimSpace(buff), []byte("\n"))

	lines := make(Grid, len(content))
	for i, c := range content {
		lines[i] = bytes.Runes(c)
	}

	return lines
}

func getStartCoord(grid Grid) Coord {
	for y, row := range grid {
		for x, char := range row {
			if charToPipe(char) == Start {
				return Coord{x, y}
			}
		}
	}

	panic("No start found")
}

func isValidCoord(coord Coord, grid Grid) bool {
	x, y := coord.x, coord.y

	if y < 0 || y >= len(grid) {
		return false
	}
	if x < 0 || x >= len(grid[0]) {
		return false
	}

	return true
}

func getNextCoord(curr Coord, nexts []Coord, grid Grid) Coord {
	for _, n := range nexts {
		if n != curr && isValidCoord(n, grid) {
			return n
		}
	}

	return Coord{-1, -1}
}

func getStartingNeighbors(start Coord, grid Grid) []Coord {
	x, y := start.x, start.y
	potentialNeighbors := []Coord{
		{x, y - 1},
		{x, y + 1},
		{x + 1, y},
		{x - 1, y},
	}

	neighbors := make([]Coord, 0)

	for _, n := range potentialNeighbors {
		if !isValidCoord(n, grid) {
			continue
		}
		if p := charToPipe(grid[n.y][n.x]); p > Invalid {
			neighborNeighbors := getNeighbours(n, p)
			for _, nn := range neighborNeighbors {
				if nn.x == x && nn.y == y {
					// fmt.Printf("(%d, %d): %c\n", n.x, n.y, pipeToDisplayChar(p))
					neighbors = append(neighbors, n)
				}
			}
		}
	}

	return neighbors
}

func getLoop(start Coord, next Coord, grid Grid) []Coord {
	loop := make([]Coord, 1)
	loop[0] = start

	curr := start

	// fmt.Println("Curr\tNext\tNeighbors")
	// for i := 0; i < 3; i++ {
	for {
		pipe := charToPipe(grid[next.y][next.x])
		allN := getNeighbours(next, pipe)
		// fmt.Printf("%v\t%v\t%v\n", curr, next, allN)

		temp := next
		next = getNextCoord(curr, allN, grid)
		curr = temp

		// Loop is broken
		if !isValidCoord(curr, grid) {
			return []Coord{}
		}

		// The loop is complete
		if curr == start {
			return loop
		}

		loop = append(loop, curr)
	}
}

func getMainLoop(grid Grid) []Coord {
	startCoord := getStartCoord(grid)

	startNeighbors := getStartingNeighbors(startCoord, grid)

	var mainLoop []Coord
	for _, n := range startNeighbors {
		l := getLoop(startCoord, n, grid)
		if len(mainLoop) < len(l) {
			mainLoop = l
		}
	}

	return mainLoop
}

func solveA(grid Grid) (solnA int) {
	mainLoop := getMainLoop(grid)
	solnA = len(mainLoop) / 2
	return
}

func getStartType(start, n1, n2 Coord) Pipe {
	x, y := start.x, start.y
	// We're scanning top-down so the y-axis in inverted
	north := Coord{x, y - 1}
	south := Coord{x, y + 1}
	east := Coord{x + 1, y}
	west := Coord{x - 1, y}

	isNorth := north == n1 || north == n2
	isSouth := south == n1 || south == n2
	isEast := east == n1 || east == n2
	isWest := west == n1 || west == n2

	if isNorth && isSouth {
		return PipeNS
	} else if isEast && isWest {
		return PipeWE
	} else if isNorth && isEast {
		return PipeNE
	} else if isNorth && isWest {
		return PipeNW
	} else if isSouth && isEast {
		return PipeSE
	} else if isSouth && isWest {
		return PipeSW
	} else {
		return Start
	}
}

func replaceStart(grid *Grid, loop []Coord) *Grid {
	// Figure out the type of start
	// The loop will always have start as the first coord
	startType := getStartType(loop[0], loop[1], loop[len(loop)-1])

	for y, row := range *grid {
		for x, char := range row {
			if charToPipe(char) == Start {
				(*grid)[y][x] = pipeToChar(startType)
			}
		}
	}
	return grid
}

func isCoordInLoop(coord Coord, loop map[Coord]bool, grid Grid) bool {
	// If the coord belongs to the loop, then it's not inside of the loop
	if _, ok := loop[coord]; ok {
		return false
	}

	// Count the number of intersections with the loop
	// Only points to the left, on the horizontal axis matter
	// Rules for counting:
	// 1. Ignore any '-'
	// 2. Any U-turns count as 2 intersections
	// 3. Any diagonal turns count as 1 intersection
	// 4. All '|' count as 1 intersection

	insertions := 0
	line := grid[coord.y]
	var prevPipe Pipe

	for i := coord.x + 1; i < len(line); i++ {
		pipe := charToPipe(line[i])
		// Only check pipe chars
		if pipe <= Invalid {
			continue
		}

		// We only check points on the loop
		if _, ok := loop[Coord{i, coord.y}]; !ok {
			continue
		}

		// Ignore horizontal pipes
		if pipe == PipeWE {
			continue
		}

		// All vertical pipes count as 1 insertion
		if pipe == PipeNS {
			prevPipe = pipe
			insertions++
			// Special checks for diagonal pipes
		} else {
			isUturn := (prevPipe == PipeNE && pipe == PipeNW) ||
				(prevPipe == PipeSE && pipe == PipeSW)

			isDiagonal := (prevPipe == PipeNE && pipe == PipeSW) ||
				(prevPipe == PipeSE && pipe == PipeNW)

			// U-turns count as 2 insertions
			if isUturn {
				insertions += 2
				// Diagonal turns count as 1 insertion
			} else if isDiagonal {
				insertions++
				// Wait for the next pipe before counting
			} else {
				prevPipe = pipe
			}
		}
	}

	return insertions%2 == 1
}

func solveB(grid Grid) (solnB int) {
	loop := getMainLoop(grid)
	replaceStart(&grid, loop)

	// Convert loop to a map for faster lookups
	loopMap := make(map[Coord]bool, len(loop))
	for _, c := range loop {
		loopMap[c] = true
	}

	for y := range grid {
		for x := range grid[y] {
			if isCoordInLoop(Coord{x, y}, loopMap, grid) {
				solnB++
			}
		}
	}

	return
}

func main() {
	grid := parseInput("day10/in.txt")

	// Preview for examples
	if len(grid[0]) < 80 {
		for _, line := range grid {
			for _, char := range line {
				pipe := charToPipe(char)
				displayChar := pipeToDisplayChar(pipe)
				fmt.Print(string(displayChar))
			}
			fmt.Println()
		}
		fmt.Println()
	}

	// solnA := solveA(grid)
	// fmt.Println("Solution A:", solnA)

	solnB := solveB(grid)
	fmt.Println("Solution B:", solnB)
}
