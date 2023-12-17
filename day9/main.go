package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func parseInput(filename string) [][]int {
	buff, _ := os.ReadFile(filename)
	buff = bytes.TrimSpace(buff)

	lines := bytes.Split(buff, []byte("\n"))
	sequences := make([][]int, len(lines))

	for i, line := range lines {
		fields := bytes.Fields(line)
		sequence := make([]int, len(fields))

		for j, f := range fields {
			sequence[j], _ = strconv.Atoi(string(f))
		}

		sequences[i] = sequence
	}

	return sequences
}

func isAllZeros(seq []int) bool {
	for _, num := range seq {
		if num != 0 {
			return false
		}
	}
	return true
}

func getPairwiseDiffs(seq []int) []int {
	diff := make([]int, len(seq)-1)

	for i := range diff {
		diff[i] = seq[i+1] - seq[i]
	}

	return diff
}

func predictNextNum(seq []int) (num int) {
	diffs := make([][]int, 0)
  diffs = append(diffs, seq)
	latestDiff := seq

	// Find Pairwise differences, util all are zero
	for !isAllZeros(latestDiff) {
		nextLayer := getPairwiseDiffs(latestDiff)
		diffs = append(diffs, nextLayer)
		latestDiff = nextLayer
	}

	// Add a `0` to the end of each layer
	for i, d := range diffs {
		diffs[i] = append(d, 0)
	}

	// Find sums between layers to predict the next num
	for i := len(diffs) - 2; i >= 0; i-- {
		lastIdx := len(diffs[i]) - 1
		diffs[i][lastIdx] = diffs[i][lastIdx-1] + diffs[i+1][lastIdx-1]
	}

	num = diffs[0][len(diffs[0])-1]
	return
}

func solveA(seqs [][]int) (acc int) {
	for _, seq := range seqs {
		num := predictNextNum(seq)
		fmt.Println(seq, "->", num)
		acc += num
	}

	return
}

func predictPrevNum(seq []int) (num int) {
	diffs := make([][]int, 0)
  diffs = append(diffs, seq)
	latestDiff := seq

	// Find Pairwise differences, util all are zero
	for !isAllZeros(latestDiff) {
		nextLayer := getPairwiseDiffs(latestDiff)
		diffs = append(diffs, nextLayer)
		latestDiff = nextLayer
	}

	// Add a `0` to the START of each layer
	for i, d := range diffs {
		diffs[i] = append([]int{0}, d...)
	}

	// Find sums between layers to predict the next num
	for i := len(diffs) - 2; i >= 0; i-- {
		diffs[i][0] = diffs[i][1] - diffs[i+1][0]
	}

	num = diffs[0][0]
	return
}

func solveB(seqs [][]int) (acc int) {
	for _, seq := range seqs {
		num := predictPrevNum(seq)
		fmt.Println(seq, "->", num)
		acc += num
	}

	return
}

func main() {
	seqs := parseInput("day9/in.txt")

	solnA := solveA(seqs)
	fmt.Println("Solution A:", solnA)

	solnB := solveB(seqs)
	fmt.Println("Solution B:", solnB)
}
