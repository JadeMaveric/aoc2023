package main

import (
	"container/heap"
	"fmt"
)

type PathFinder interface {
	FindPath(grid *Grid, start, end *Block) (path []*Block)
}

type MockPathFinder struct{}

func (pf MockPathFinder) FindPath(g *Grid, start, end *Block) (path []*Block) {
	grid := *g
	path = make([]*Block, 0)

	// Run algo
	mock := []*Block{grid[0][0], grid[0][1], grid[1][1], grid[2][1], grid[2][2], grid[1][2], grid[1][3], grid[2][3]}
	for _, p := range mock {
		path = append(path, p)
	}

	// Build path
	for i, p := range path {
		if i == 0 {
			p.SetLoss(p.posX, p.posY, p, 1)
		} else {
			p.SetLoss(p.posX, p.posY, path[i-1], 1)
		}
	}
	return path
}

func isWrongDir(curr, next int) bool {
	switch curr {
	case DirN:
		return next == DirS
	case DirS:
		return next == DirN
	case DirE:
		return next == DirW
	case DirW:
		return next == DirE
	default:
		return false
	}
}

func getNextStep(curr *QueueStep, nDir int) (int, int) {
	if curr.dir == nDir {
		xStep, yStep := curr.xStep, curr.yStep

		switch nDir {
		case DirN:
			return xStep, yStep - 1
		case DirS:
			return xStep, yStep + 1
		case DirE:
			return xStep + 1, yStep
		case DirW:
			return xStep - 1, yStep
		default:
			return xStep, yStep
		}
	}

	switch nDir {
	case DirN:
		return -0, -1
	case DirS:
		return +0, +1
	case DirE:
		return +1, +0
	case DirW:
		return -1, -0
	default:
		return +0, -0
	}
}

type AStarPathFinderSimple struct{}

func (pf AStarPathFinderSimple) FindPath(g *Grid, start, end *Block) (path []*Block) {
	grid := *g

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Initial conditions
	curr := QueueStep{block: start, dir: DirX, xStep: 0, yStep: 0}

	for curr.block != end {
		for _, neighbor := range grid.Neighbors(curr.block) {
			nDir := getDirection(curr.block, neighbor)
			if isWrongDir(curr.dir, nDir) {
				continue
			}

			xStep, yStep := getNextStep(&curr, nDir)
			if xStep > RANGE || xStep < -RANGE || yStep > RANGE || yStep < -RANGE {
				continue
			}
			// xStep, yStep := curr.xStep, curr.yStep
			nextLoss := curr.lossSoFar + neighbor.heatLoss

			if nextLoss < neighbor.GetLoss(xStep, yStep) {
				neighbor.SetLoss(xStep, yStep, curr.block, nextLoss)
				heap.Push(&pq, &QueueStep{block: neighbor, lossSoFar: nextLoss, dir: nDir, xStep: xStep, yStep: yStep})
			}
		}

		curr = *heap.Pop(&pq).(*QueueStep)
	}

	fmt.Println("Done traversing", end.BestLoss().bestLoss)

	// Build path
	path = make([]*Block, 0)
	currNode := end
	for currNode != start {
		// fmt.Printf("(%d,%d)[%d]\n", currNode.posX, currNode.posY, currNode.heatLoss)
		path = append([]*Block{currNode}, path...)
		currNode = currNode.BestLoss().cameFrom
	}

	return path
}


type AStarPathFinderUltra struct{}

func (pf AStarPathFinderUltra) FindPath(g *Grid, start, end *Block) (path []*Block) {
	grid := *g

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Initial conditions
	curr := QueueStep{block: start, dir: DirX, xStep: 0, yStep: 0}

	for curr.block != end {
    dirChange := (curr.xStep + curr.yStep) <= 1
		for _, neighbor := range grid.UltraNeighbors(curr.block, dirChange) {
			nDir := getDirection(curr.block, neighbor)
			if isWrongDir(curr.dir, nDir) {
				continue
			}

			xStep, yStep := getNextStep(&curr, nDir)
			if xStep > RANGE || xStep < -RANGE || yStep > RANGE || yStep < -RANGE {
				continue
			}
			// xStep, yStep := curr.xStep, curr.yStep
			nextLoss := curr.lossSoFar + grid.GetTotalLoss(curr.block, neighbor)

			if nextLoss < neighbor.GetLoss(xStep, yStep) {
				neighbor.SetLoss(xStep, yStep, curr.block, nextLoss)
				heap.Push(&pq, &QueueStep{block: neighbor, lossSoFar: nextLoss, dir: nDir, xStep: xStep, yStep: yStep})
			}
		}

		curr = *heap.Pop(&pq).(*QueueStep)
	}

	fmt.Println("Done traversing", end.BestLoss().bestLoss)

	// Build path
	path = make([]*Block, 0)
	currNode := end
	for currNode != start {
		// fmt.Printf("(%d,%d)[%d]\n", currNode.posX, currNode.posY, currNode.heatLoss)
		path = append([]*Block{currNode}, path...)
		currNode = currNode.BestLoss().cameFrom
	}

	return path
}
