package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Step struct {
	label string
	op    byte
	focal int
}

type Lens struct {
	label string
	focal int
}

type Box []Lens

func (b Box) Add(l Lens) Box {
  for i, l1 := range b {
    if l1.label == l.label {
      b[i] = l
      return b
    }
  }
	return append(b, l)
}

func (b Box) Remove(label string) Box {
	b1 := make(Box, len(b))

	count := 0

	for _, l := range b {
		if l.label != label {
			b1[count] = l
			count++
		}
	}

	return b1[:count]
}

func parseInput(filename string) []string {
	buff, _ := os.ReadFile(filename)
	buff = bytes.TrimSpace(buff)
	buff = bytes.ReplaceAll(buff, []byte("\n"), []byte(""))

	codes := bytes.Split(buff, []byte(","))
	steps := make([]string, len(codes))
	for i, code := range codes {
		steps[i] = string(code)
	}

	return steps
}

func parseStep(step string) Step {
	isAdd, isRemove := false, false

	for _, c := range step {
		if c == '-' {
			isRemove = true
			break
		}

		if c == '=' {
			isAdd = true
			break
		}
	}

	if isAdd {
		fields := strings.Split(step, "=")
		label := fields[0]
		focal, _ := strconv.Atoi(fields[1])
		return Step{label, '=', focal}
	}

	if isRemove {
		label := step[:len(step)-1]
		return Step{label, '-', 0}
	}

	panic(fmt.Sprintf("Invalid step: %s\n", step))
}

func HASH(s string) (hash int) {
	for _, c := range s {
		ascii := int(c)
		hash += ascii
		hash *= 17
		hash = hash % 256
	}
	return hash
}

func getScore(boxes []Box) (score int) {
  for i, box := range boxes {
    for j, lens := range box {
      score += (i+1) * (j+1) * lens.focal
    }
  }
  return
}

func solveA(steps []string) (soln int) {
	for _, step := range steps {
		soln += HASH(step)
	}
	return
}

func solveB(steps []string) (soln int) {
  box := [256]Box{}
  for i := range box {
    box[i] = make([]Lens, 0)
  }

  for _, step := range steps {
    s := parseStep(step)
    bIdx := HASH(s.label)
    b := box[bIdx]

    if s.op == '=' {
      b = b.Add(Lens{s.label, s.focal})
    } else {
      b = b.Remove(s.label)
    }
    // fmt.Println(b)
    box[bIdx] = b
  }

	return getScore(box[:])
}

func main() {
	steps := parseInput("day15/in.txt")

	// Sanity test
	// hash := HASH("HASH")
	// fmt.Println(hash) -- 52

	solnA := solveA(steps)
	fmt.Println("A:", solnA)

  solnB := solveB(steps)
  fmt.Println("B:", solnB)
}
