package main

import (
	"bytes"
	"fmt"
	"os"
)

type Row []byte
type Pattern [][]byte
type Input [][][]byte

func parseInput(filename string) Input {
	buff, _ := os.ReadFile(filename)
	buff = bytes.TrimSpace(buff)

	sections := bytes.Split(buff, []byte("\n\n"))

	input := make([][][]byte, len(sections))
	for i, section := range sections {
		lines := bytes.Split(section, []byte("\n"))
		input[i] = lines
	}

	return input
}

func isValidRowPos(row Row, pos, smudgeLimit int) (valid bool, smudgeCount int) {
	// Can't be on the edge
	if pos == 0 || pos == len(row) {
		return false, 0
	}

	for i := 0; i < pos; i++ {
		leftIdx := pos - i - 1
		rightIdx := pos + i

		if leftIdx < 0 || rightIdx >= len(row) {
			break
		}

		if row[leftIdx] != row[rightIdx] {
			smudgeCount++
			if smudgeCount > smudgeLimit {
				return false, 0
			}
		}
	}

	return true, smudgeCount
}

type CandidatePos struct {
	pos         int
	smudgeCount int
}

func getMirrorRowPositions(row Row, candidatePos []CandidatePos, smudgeLimit int) (pos []CandidatePos) {
	pos = make([]CandidatePos, 0)

	if candidatePos == nil {
		for i := range row {
			valid, smudgeCount := isValidRowPos(row, i, smudgeLimit)
			if valid {
				pos = append(pos, CandidatePos{i, smudgeCount})
			}
		}
	} else {
		for _, p := range candidatePos {
			valid, smudgeCount := isValidRowPos(row, p.pos, smudgeLimit)
			newP := CandidatePos{p.pos, p.smudgeCount + smudgeCount}

			if valid && newP.smudgeCount <= smudgeLimit {
				pos = append(pos, newP)
			}
		}
	}

	return
}

func transpose(pattern Pattern) Pattern {
	rows := len(pattern)
	cols := len(pattern[0])

	transposed := make(Pattern, cols)
	for y := range transposed {
		transposed[y] = make(Row, rows)

		for x := range transposed[y] {
			transposed[y][x] = pattern[x][y]
		}
	}

	return transposed
}

func getScore(input Input, smudgeLimit int) (score int) {
	xMatches := 0
	yMatches := 0

	for _, pattern := range input {
		// Check for X-axis symmetry
		var xCandidatePos []CandidatePos
		for _, line := range pattern {
			xCandidatePos = getMirrorRowPositions(line, xCandidatePos, smudgeLimit)
			// fmt.Println(string(line), xCandidatePos)
			if len(xCandidatePos) == 0 {
				break
			}
		}
		for _, pos := range xCandidatePos {
			if pos.smudgeCount == smudgeLimit {
				xMatches += pos.pos
			}
		}

		// fmt.Println("")

		// Check for Y-axis symmetry
		var yCandidatePos []CandidatePos
		for _, line := range transpose(pattern) {
			yCandidatePos = getMirrorRowPositions(line, yCandidatePos, smudgeLimit)
			// fmt.Println(string(line), yCandidatePos)
			if len(yCandidatePos) == 0 {
				break
			}
		}
		for _, pos := range yCandidatePos {
			if pos.smudgeCount == smudgeLimit {
				yMatches += pos.pos
			}
		}

		// totalMatches := len(xCandidatePos) + len(yCandidatePos)
		// fmt.Printf("Summary:\t%d\t%d\t%d\n", totalMatches, xCandidatePos, yCandidatePos)
	}

	return xMatches + 100*yMatches
}

func solveA(input Input) (solnA int) {
	return getScore(input, 0)
}

func solveB(input Input) (solnB int) {
	return getScore(input, 1)
}

func main() {
	input := parseInput("day13/in.txt")

	solnA := solveA(input)
	fmt.Println("A:", solnA)

	solnB := solveB(input)
	fmt.Println("B:", solnB)
}
