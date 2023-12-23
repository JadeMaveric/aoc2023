package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

type Tile byte

const (
	TileSpace       Tile = '.'
	TileMirrorRight Tile = '/'
	TileMirrorLeft  Tile = '\\'
	TileSplitX      Tile = '-'
	TileSplitY      Tile = '|'
)

type Cell struct {
	tile           Tile
	pN, pS, pE, pW bool
}
type Grid [][](*Cell)

type Dir struct{ x, y int }

type Photon struct {
	x, y int
	dir  Dir
}

func (photon Photon) incrementPos() Photon {
	photon.x += photon.dir.x
	photon.y -= photon.dir.y // y-axis in inverted on the grid
	return photon
}

func (grid *Grid) String() string {
	var buff string

	for _, line := range *grid {
		for _, cell := range line {
			buff += fmt.Sprint(string(cell.tile))
		}
		buff += fmt.Sprintln()
	}

	return buff
}

func (grid *Grid) EnergyString() string {
	var buff string

	for _, line := range *grid {
		for _, cell := range line {
			if cell.pN || cell.pS || cell.pE || cell.pW {
				buff += "#"
			} else {
				buff += "."
			}
		}
		buff += fmt.Sprintln()
	}

	return buff
}

func (grid *Grid) EnergyCountString() string {
	var buff string

	for _, line := range *grid {
		for _, cell := range line {
			if cell.pN || cell.pS || cell.pE || cell.pW {
				buff += fmt.Sprint(countTrue(cell.pN, cell.pS, cell.pE, cell.pW))
			} else {
				buff += "."
			}
		}
		buff += fmt.Sprintln()
	}

	return buff
}

func (grid *Grid) EnergyCount() (count int) {
	for _, line := range *grid {
		for _, cell := range line {
			if cell.pN || cell.pS || cell.pE || cell.pW {
				count++
			}
		}
	}
	return
}

func (grid *Grid) ResetEnergy() {
	for _, line := range *grid {
		for _, cell := range line {
			cell.pN = false
			cell.pS = false
			cell.pE = false
			cell.pW = false
		}
	}
}

func countTrue(b ...bool) (count int) {
	for _, v := range b {
		if v {
			count++
		}
	}
	return
}

func parseInput(filename string) Grid {
	buff, _ := os.ReadFile(filename)
	buff = bytes.TrimSpace(buff)

	lines := bytes.Split(buff, []byte("\n"))

	grid := make(Grid, len(lines))
	for y, line := range lines {
		line = bytes.TrimSpace(line)
		grid[y] = make([](*Cell), len(line))
		for x, char := range line {
			cell := &Cell{tile: Tile(char)}
			grid[y][x] = cell
		}
	}

	return grid
}

func tick(photon Photon, tile Tile) []Photon {
	// Calc direction
	if tile == TileSpace {
		return []Photon{photon.incrementPos()}
	} else if tile == TileMirrorRight {
		t := photon.dir.x
		photon.dir.x = photon.dir.y
		photon.dir.y = t
		return []Photon{photon.incrementPos()}
	} else if tile == TileMirrorLeft {
		t := photon.dir.x
		photon.dir.x = -photon.dir.y
		photon.dir.y = -t
		return []Photon{photon.incrementPos()}
	} else if tile == TileSplitY {
		if photon.dir.x == 0 {
			return []Photon{photon.incrementPos()}
		} else {
			photonA := Photon{x: photon.x, y: photon.y, dir: Dir{x: 0, y: -1}}.incrementPos()
			photonB := Photon{x: photon.x, y: photon.y, dir: Dir{x: 0, y: +1}}.incrementPos()
			return []Photon{photonA, photonB}
		}
	} else if tile == TileSplitX {
		if photon.dir.y == 0 {
			return []Photon{photon.incrementPos()}
		} else {
			photonA := Photon{x: photon.x, y: photon.y, dir: Dir{x: -1, y: 0}}.incrementPos()
			photonB := Photon{x: photon.x, y: photon.y, dir: Dir{x: +1, y: 0}}.incrementPos()
			return []Photon{photonA, photonB}
		}
	}

	return []Photon{}
}

func validPos(grid *Grid, x, y int) bool {
	return x >= 0 && y >= 0 && y < len(*grid) && x < len((*grid)[y])
}

func (grid *Grid) traceRay(start Photon) (photons []Photon) {
	// TODO: Figure out exit condition
	// Exit condition:
	// 1. photon is out of bounds
	// 2. OR photon is trapped in a loop

	photons = make([]Photon, 1)
	photons[0] = start

	for len(photons) > 0 {
		// Pop first photon
		photon := photons[0]
		photons = photons[1:]

		if !validPos(grid, photon.x, photon.y) {
			continue
		}

		cell := (*grid)[photon.y][photon.x]

		// Check if we've processed this tile
		// If we have, we're in a loop
		if photon.dir.x == 1 {
			if cell.pE {
				continue
			} else {
				cell.pE = true
			}
		} else if photon.dir.x == -1 {
			if cell.pW {
				continue
			} else {
				cell.pW = true
			}
		} else if photon.dir.y == 1 {
			if cell.pN {
				continue
			} else {
				cell.pN = true
			}
		} else if photon.dir.y == -1 {
			if cell.pS {
				continue
			} else {
				cell.pS = true
			}
		}

		newPhotons := tick(photon, cell.tile)
		if len(newPhotons) > 0 {
			photons = append(photons, newPhotons...)
		}
	}

	return
}

func solveA(grid *Grid) (solnA int) {
	if len(*grid) < 15 {
		fmt.Println(grid.String())
	}

	startPhoton := Photon{x: 0, y: 0, dir: Dir{x: 1, y: 0}}
	grid.traceRay(startPhoton)

	if len(*grid) < 15 {
		fmt.Println(grid.EnergyString())
		fmt.Println(grid.EnergyCountString())
	}

	for _, char := range grid.EnergyString() {
		if char == '#' {
			solnA++
		}
	}

	return solnA
}

func getStartPhoton(x, y, w, h int) *Photon {
	if x == -1 && y > -1 && y < h {
		return &Photon{x: 0, y: y, dir: Dir{x: 1, y: 0}}
	}

	if x == w && y > -1 && y < h {
		return &Photon{x: w - 1, y: y, dir: Dir{x: -1, y: 0}}
	}

	if y == -1 && x > -1 && x < w {
		return &Photon{x: x, y: 0, dir: Dir{x: 0, y: 1}}
	}

	if y == h && x > -1 && x < w {
		return &Photon{x: x, y: h - 1, dir: Dir{x: 0, y: -1}}
	}

	return nil
}

func solveB(grid *Grid) (solnB int) {
	defer timer("solveB")()
	w, h := len((*grid)[0]), len(*grid)

	var maxPhoton *Photon

	for x := -1; x <= w; x++ {
		for y := -1; y <= h; y++ {
			// Get start photon
			startPhoton := getStartPhoton(x, y, w, h)
			if startPhoton == nil {
				continue
			}

			// Calculate engery count
			grid.ResetEnergy()
			grid.traceRay(*startPhoton)
			energyCount := grid.EnergyCount()

			if energyCount > solnB {
				solnB = energyCount
				maxPhoton = startPhoton
			}
		}
	}

	fmt.Println("maxPhoton:", *maxPhoton)

	return solnB
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Println(name, "took", time.Since(start))
	}
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	grid := parseInput("day16/in.txt")

	solnA := solveA(&grid)
	fmt.Println("A:", solnA)

	solnB := solveB(&grid)
	fmt.Println("B:", solnB)
}
