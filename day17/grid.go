package main

import (
	"bytes"
	"math"
)

const RANGE = 3
const SIZE = 2*RANGE + 1

type Block struct {
	heatLoss   int
	lossGrid   [SIZE * SIZE]*Step
	posX, posY int
}

type Step struct {
	cameFrom     *Block
	bestLoss     int
	xStep, yStep int
}

func (bPtr *Block) GetHeuristicCost() int {
	// Inverse manhattan distance since
	// target is at the bottom-right corner
	return -1 * (bPtr.posX + bPtr.posY)
}

func (bPtr *Block) IsSeen() bool {
	for _, s := range bPtr.lossGrid {
		if s != nil {
			return true
		}
	}

	return false
}

func (bPtr *Block) SetLoss(xStep, yStep int, cameFrom *Block, loss int) {
	idx := (yStep+RANGE)*SIZE + (xStep + RANGE)
	step := &Step{
		cameFrom: cameFrom,
		bestLoss: loss,
		xStep:    xStep,
		yStep:    yStep,
	}
	bPtr.lossGrid[idx] = step
}

func (bPtr *Block) GetLoss(xStep, yStep int) int {
	idx := (yStep+3)*7 + (xStep + 3)
	if bPtr.lossGrid[idx] == nil {
		return math.MaxInt
	}
	return bPtr.lossGrid[idx].bestLoss
}

func (bPtr *Block) BestLoss() *Step {
	minLoss := math.MaxInt
	var bestStep *Step = nil
	for _, s := range bPtr.lossGrid {
		if s == nil {
			continue
		}
		if s.bestLoss < minLoss {
			minLoss = s.bestLoss
			bestStep = s
		}
	}
	return bestStep
}

type Grid [][](*Block)

func (gPtr *Grid) String() string {
	g := *gPtr
	strGrid := make([][]byte, len(g))

	for y := range g {
		strGrid[y] = make([]byte, len(g[y]))

		for x := range g[y] {
			strGrid[y][x] = byte(g[y][x].heatLoss + '0')
		}
	}

	return string(bytes.Join(strGrid, []byte("\n")))
}

func getDirection(curr, next *Block) int {
	if curr.posY == next.posY {
		if curr.posX == next.posX {
			return DirX
		} else if curr.posX > next.posX {
			return DirW
		} else {
			return DirE
		}
	} else {
		if curr.posY > next.posY {
			return DirN
		} else {
			return DirS
		}
	}
}

func getDirectionChar(x1, y1, x2, y2 int) byte {
	// same row
	if y1 == y2 {
		if x1 == x2 {
			return '='
		} else if x1 > x2 {
			return '<'
		} else {
			return '>'
		}
	} else {
		if y1 > y2 {
			return '^'
		} else {
			return 'v'
		}
	}
}

func (gPtr *Grid) VisitedString(path []*Block) string {
	g := *gPtr
	strGrid := make([][]byte, len(g))

	for y := range g {
		strGrid[y] = make([]byte, len(g[y]))

		for x := range g[y] {
			block := g[y][x]
			if block.IsSeen() {
				strGrid[y][x] = byte(block.heatLoss + '0')
			} else {
				strGrid[y][x] = '.'
			}
		}
	}

	for _, block := range path {
		minLoss := math.MaxInt
		var prevBlock *Block = nil

		for _, s := range block.lossGrid {
			if s == nil {
				continue
			}
			if s.bestLoss < minLoss {
				minLoss = s.bestLoss
				prevBlock = s.cameFrom
			}
		}

		if prevBlock == nil {
			continue
		}

		dirChar := getDirectionChar(prevBlock.posX, prevBlock.posY, block.posX, block.posY)
		strGrid[prevBlock.posY][prevBlock.posX] = dirChar

		strGrid[block.posY][block.posX] = byte('X')
	}

	return string(bytes.Join(strGrid, []byte("\n")))
}

func (gPtr *Grid) Neighbors(curr *Block) []*Block {
	g := *gPtr
	neighbors := make([]*Block, 0)
	// up
	if curr.posY > 0 {
		neighbors = append(neighbors, g[curr.posY-1][curr.posX])
	}
	// down
	if curr.posY < len(g)-1 {
		neighbors = append(neighbors, g[curr.posY+1][curr.posX])
	}
	// left
	if curr.posX > 0 {
		neighbors = append(neighbors, g[curr.posY][curr.posX-1])
	}
	// right
	if curr.posX < len(g[0])-1 {
		neighbors = append(neighbors, g[curr.posY][curr.posX+1])
	}
	return neighbors
}
