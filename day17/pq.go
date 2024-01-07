package main

import (
	"fmt"
	"strings"
)

// PRIORITY QUEUE

const (
  DirX = iota
  DirN = iota
  DirS = iota
  DirE = iota
  DirW = iota
)

type QueueStep struct {
	block        *Block
	dir          int
  lossSoFar    int
	xStep, yStep int
	index        int
}

type PriorityQueue []*QueueStep

func (q PriorityQueue) Len() int { return len(q) }

func (q PriorityQueue) Less(i, j int) bool {
	iCost := q[i].block.GetLoss(q[i].xStep, q[i].yStep) + q[i].block.GetHeuristicCost()
	jCost := q[j].block.GetLoss(q[j].xStep, q[j].yStep) + q[j].block.GetHeuristicCost()
	return iCost < jCost
}

func (q PriorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *PriorityQueue) Push(step interface{}) {
	n := len(*q)
	step.(*QueueStep).index = n
	*q = append(*q, step.(*QueueStep))
}

func (q *PriorityQueue) Pop() interface{} {
	old := *q
	n := len(old)
	step := old[n-1]
	old[n-1] = nil  // avoid memory leak
	step.index = -1 // for safety
	*q = old[0 : n-1]
	return step
}

func (q *PriorityQueue) String() string {
	blockBuff := make([]string, 0, len(*q))
	for _, step := range *q {
		block := step.block
		blockBuff = append(blockBuff, fmt.Sprintf("{(%d,%d), %d}", block.posX, block.posY, block.BestLoss().bestLoss))
	}

	return "[" + strings.Join(blockBuff, ", ") + "]"
}
