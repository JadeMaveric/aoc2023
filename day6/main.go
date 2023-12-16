package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// Part A: Math
// T -> time limit
// s -> speed :: duration we held down the button for
// d -> distance
//   : d = speed * time_for_travel
//   : d = s     * (T - s)
//
// L -> record time
// To beat the race, I need to solve
//   : d > L
//   : s * (T - S) > L
//
// So, find the zeros of the equation
//   : d = s * (T - s) - D
//   : y = x * (T - x) - D
//   : 0 = xT - x^2 - D
//   : 0 = -1*x^2 + T*x - D
//
//   : x = (-b + sqrt(b^2 - 4ac)) / 2ab
//   : x = (-T + sqrt(T^2 - 4D)) / -2T
//   : x = (T +/- sqrt(T^2 - 4D)) / 2T

type Race struct {
	time     int
	distance int
}

func parseIntArr(line string) []int {
	fields := strings.Fields(line)

	nums := make([]int, len(fields))
	for i, f := range fields {
		nums[i], _ = strconv.Atoi(f)
	}

	return nums
}

func timer(name string) func() {
	start := time.Now()

	return func() {
		end := time.Now()
		s := end.Sub(start)
		fmt.Println(name, "took", s)
	}
}

func parseInputA(filename string) []Race {
	buff, _ := os.ReadFile(filename)

	lines := bytes.Split(buff, []byte("\n"))

	timeLine := bytes.Split(lines[0], []byte(":"))[1]
	distLine := bytes.Split(lines[1], []byte(":"))[1]

	times := parseIntArr(string(timeLine))
	distances := parseIntArr(string(distLine))

	if len(times) != len(distances) {
		panic(fmt.Sprintf("Times and distances have difference lengths: %d, %d", len(times), len(distances)))
	}

	races := make([]Race, len(times))
	for i := range times {
		races[i] = Race{time: times[i], distance: distances[i]}
	}

	return races
}

func parseInputB(filename string) Race {
	buff, _ := os.ReadFile(filename)

	lines := strings.Split(string(buff), "\n")

	timeLine := strings.Split(lines[0], ":")[1]
	timeLine = strings.ReplaceAll(timeLine, " ", "")
	time, _ := strconv.Atoi(timeLine)

	distLine := strings.Split(lines[1], ":")[1]
	distLine = strings.ReplaceAll(distLine, " ", "")
	dist, _ := strconv.Atoi(distLine)

	return Race{time: time, distance: dist}
}

func solveA(races []Race) int {
	defer timer("solveA")()
	acc := 1

	for _, race := range races {
		T := float64(race.time)
		D := float64(race.distance)

		numA := T
		numB := math.Sqrt(float64(T*T - 4*D))
		denum := float64(2)

		zeroA := (numA - numB) / denum
		zeroB := (numA + numB) / denum
		// fmt.Printf("(%.2f, %.2f)\n", zeroA, zeroB)

		ansA := math.Floor(zeroA) + 1
		ansB := math.Ceil(zeroB) - 1
		// fmt.Printf("(%.2f, %.2f)\n", ansA, ansB)

		waysToWin := int(ansB - ansA + 1)
		// fmt.Println("-->", waysToWin)
		acc *= waysToWin
	}

	return acc
}

func main() {
	defer timer("main")()

	races := parseInputA("day6/in.txt")
	solnA := solveA(races)
	fmt.Println("Solution A:", solnA)

	race := parseInputB("day6/in.txt")
	solnB := solveA([]Race{race})
	fmt.Println("Solution A:", solnB)
}
