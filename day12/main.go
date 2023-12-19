package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Spring rune

const (
	SpringUnknown Spring = '?'
	SpringDamaged Spring = '#'
	SpringActive  Spring = '.'
)

type Row struct {
	springs []Spring
	ecc     []int
}

type Input []Row

type Stack [][]Spring

func (s Stack) Push(v []Spring) Stack {
	return append(s, v)
}

func (s Stack) Pop() (Stack, []Spring) {
	l := len(s)
	return s[:l-1], s[l-1]
}

func parseInput(filename string) (input Input) {
	buff, _ := os.ReadFile(filename)
	buff = bytes.TrimSpace(buff)

	lines := bytes.Split(buff, []byte("\n"))
	input = make(Input, len(lines))

	for i, line := range lines {
		fields := bytes.Fields(line)
		springs := fields[0]
		ecc := []byte{}
		if len(fields) > 1 {
			ecc = fields[1]
		}

		input[i].springs = make([]Spring, len(springs))
		for j, spring := range springs {
			input[i].springs[j] = Spring(spring)
		}

		codes := bytes.Split(ecc, []byte(","))
		input[i].ecc = make([]int, len(codes))
		for j, code := range codes {
			input[i].ecc[j], _ = strconv.Atoi(string(code))
		}
	}

	return
}

func generateEcc(springs []Spring, partial bool) (ecc []int) {
	ecc = make([]int, 0)

	count := 0
	for _, s := range springs {
		if s == SpringDamaged {
			count++
		}

		if s == SpringActive && count > 0 {
			ecc = append(ecc, count)
			count = 0
		}

		if s == SpringUnknown {
			if partial {
				break
			} else {
				return
			}
		}
	}

	if count > 0 {
		ecc = append(ecc, count)
	}

	return
}

// Is `b` a prefix of `a`?
func isPrefix(a, b []int) bool {
	if len(a) < len(b) {
		return false
	}

	for i := range b {
		// Special case to early exit out of `##??? - 1`
		if i == len(b)-1 && a[i] > b[i] {
			return true
			// Positive match: `##??`
			// Negative match: `##.`

		} else if a[i] != b[i] {
			return false
		}
	}

	return true
}

// Is `a` the same as `b`?
func isMatch(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range b {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func numOfValidPos(row Row, idx, total int) (count int) {
	stack := make(Stack, 0).Push(row.springs)
	processed := 0
	// validPos := make([][]Spring, 0)

	unknownCount := 0
	for _, s := range row.springs {
		if s == SpringUnknown {
			unknownCount++
		}
	}
	fmt.Printf("[%d/%d]: bitcount=%d ", idx, total, unknownCount)

	for len(stack) > 0 {
		s, curr := stack.Pop()
		stack = s // idk why it needs to be this way

		ecc := generateEcc(curr, true)
		if !isPrefix(row.ecc, ecc) {
			continue
		} else {
			processed++
		}

		// Fill the next iteration

		// Find the position of the `?` to be replaced
		p := -1
		for i, s := range curr {
			if s == SpringUnknown {
				p = i
				break
			}
		}

		if p == -1 {
			ecc := generateEcc(curr, false)
			if isMatch(row.ecc, ecc) {
				// validPos = append(validPos, curr)
				count++
			} else {
				continue
			}
		} else {
			activeCandiate := make([]Spring, len(row.springs))
			copy(activeCandiate, curr)
			activeCandiate[p] = SpringActive
			stack = stack.Push(activeCandiate)

			brokenCandiate := make([]Spring, len(row.springs))
			copy(brokenCandiate, curr)
			brokenCandiate[p] = SpringDamaged
			stack = stack.Push(brokenCandiate)
		}
	}

	fmt.Println("Proccessed:", processed, " -- count:", count)
	// for _, pos := range validPos {
	//   for _, s := range pos {
	//     fmt.Printf("%c", rune(s))
	//   }
	//   fmt.Printf("\t%v\n", generateEcc(pos))
	// }
	// fmt.Println()

	return count
}

func solveA(in Input) (solnA int) {
	defer timer("solveA")()
	for i, row := range in {
		solnA += numOfValidPos(row, i, len(in))
	}
	return solnA
}

func unfoldInput(in Input, factor int) (unfoldedIn Input) {
	unfoldedIn = make(Input, len(in))

	for i, row := range in {
		uRow := unfoldedIn[i]
		uRow.springs = make([]Spring, 0)
		uRow.ecc = make([]int, 0)

		for j := 0; j < factor; j++ {
			uRow.springs = append(uRow.springs, row.springs...)
			uRow.ecc = append(uRow.ecc, row.ecc...)

			if j != factor-1 {
				uRow.springs = append(uRow.springs, SpringUnknown)
			}
		}

		unfoldedIn[i] = uRow
	}

	return unfoldedIn
}

func isValidSpan(row []Spring, size int) bool {
	if len(row) < (size + 1) {
		return false
	}

	for i := 0; i < size; i++ {
		if row[i] == SpringActive {
			return false
		}
	}

	if row[size] == SpringDamaged {
		return false
	} else {
		return true
	}
}

func countValidPos(row []Spring, ecc []int, cache map[string]int) (count, cacheHit, processed int) {
	hash := string(row) + fmt.Sprint(ecc)
	if c, ok := cache[hash]; ok {
		count += c
    cacheHit++
		return
	} else {
		processed++
	}

	if len(row) == 0 {
		if len(ecc) == 0 {
			count = 1
			return
		} else {
			count = 0
			return
		}
	}

	if row[0] == SpringActive || row[0] == SpringUnknown {
		nCount, nCacheHit, nProcessed := countValidPos(row[1:], ecc, cache)
		count += nCount
		cacheHit += nCacheHit
		processed += nProcessed
	}

	if len(ecc) > 0 && (row[0] == SpringDamaged || row[0] == SpringUnknown) {
		spanSize := ecc[0]
		if isValidSpan(row, spanSize) {
			nCount, nCacheHit, nProcessed := countValidPos(row[spanSize+1:], ecc[1:], cache)
			count += nCount
			cacheHit += nCacheHit
			processed += nProcessed
		}
	}

	// fmt.Println(hash+":", count)

	cache[hash] = count
	return
}

func solveB(in Input) (solnB int) {
	defer timer("solveB")()
	unfoldedIn := unfoldInput(in, 5)

	cache := make(map[string]int)
	cacheHit, processed := 0, 0

	for _, row := range unfoldedIn {
		// Add a `.` to make boundary checking easier
		row.springs = append(row.springs, SpringActive)

		count, nCacheHit, nProcessed := countValidPos(row.springs, row.ecc, cache)
		cacheHit += nCacheHit
		processed += nProcessed

		hash := string(row.springs) + fmt.Sprint(row.ecc)
		fmt.Println(hash+":", count)

		solnB += count
	}

	fmt.Println("Processed:", processed, " -- Cache Hit:", cacheHit)

	return solnB
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func main() {
	in := parseInput("day12/in.txt")

	// if len(in[0].springs) < 80 {
	//   for y := range in {
	//     for x := range in[y].springs {
	//       fmt.Printf("%c", in[y].springs[x])
	//     }
	//     fmt.Print(" ")

	//     for x := range in[y].ecc {
	//       fmt.Printf("%d,", in[y].ecc[x])
	//     }
	//     fmt.Println()
	//   }
	//   fmt.Println()
	// }

	// solnA := solveA(in)
	// fmt.Println("A:", solnA)

	// This takes too long!!
	solnB := solveB(in)
	fmt.Println("B:", solnB)
}
