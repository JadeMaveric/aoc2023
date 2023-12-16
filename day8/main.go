package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Map struct {
	directions string
	nodes      map[string][]string
}

func parseInput(filename string) (doc Map) {
	buff, _ := os.ReadFile(filename)
	lines := strings.Split(strings.TrimSpace(string(buff)), "\n")

	doc.directions = lines[0]
	doc.nodes = make(map[string][]string, len(lines)-2)

	re := regexp.MustCompile(`(\w{3}) = \((\w{3}), (\w{3})\)`)

	for _, line := range lines[2:] {
		match := re.FindStringSubmatch(line)
		node, left, right := match[1], match[2], match[3]
		doc.nodes[node] = []string{left, right}
	}

	return doc
}

func nextDir(directions string, currIdx int) (nextIdx int) {
	if currIdx == len(directions)-1 {
		return 0
	}
	return currIdx + 1
}

func solveA(doc Map) (steps int) {
	currNode := "AAA"
	currIdx := 0

	for currNode != "ZZZ" {
		steps++
		currDir := doc.directions[currIdx]

		if currDir == 'L' {
			currNode = doc.nodes[currNode][0]
		} else {
			currNode = doc.nodes[currNode][1]
		}

		currIdx = nextDir(doc.directions, currIdx)
	}

	return steps
}

func getStartingNodes(nodes map[string][]string) (startingNodes []string) {
	for node, _ := range nodes {
		if strings.HasSuffix(node, "A") {
			startingNodes = append(startingNodes, node)
		}
	}

	return
}

func GCD(a, b int) int {
	if b == 0 {
		return a
	}
	return GCD(b, a%b)
}

func LCMMultiple(nums []int) (lcm int) {
	lcm = nums[0]
	for _, num := range nums[1:] {
		lcm = (lcm * num) / GCD(lcm, num)
	}
	return
}

func solveB(doc Map) int {
	currNodes := getStartingNodes(doc.nodes)
	nodeSteps := make([]int, len(currNodes))

	for i, currNode := range currNodes {
		currIdx := 0
		nodeStep := 0

		for !strings.HasSuffix(currNode, "Z") {
			nodeStep++
			currDir := doc.directions[currIdx]

			if currDir == 'L' {
				currNode = doc.nodes[currNode][0]
			} else {
				currNode = doc.nodes[currNode][1]
			}

			currIdx = nextDir(doc.directions, currIdx)
		}

		nodeSteps[i] = nodeStep
	}

	fmt.Println(nodeSteps)

	return LCMMultiple(nodeSteps)
}

func main() {
	doc := parseInput("day8/c.txt")
	// fmt.Println(doc)

	// solnA := solveA(doc)
	// fmt.Println("Solution A:", solnA)

	solnB := solveB(doc)
	fmt.Println("Solution B:", solnB)
}
