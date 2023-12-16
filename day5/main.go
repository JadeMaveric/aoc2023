package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

type Range struct {
	SrcStart  int
	DestStart int
	Length    int
}

type Mapping struct {
	SrcType  string
	DestType string
	Ranges   []Range
}

type Almanac struct {
	Seeds    []int
	Mappings []Mapping
}

func parseIntArr(line string) []int {
	fields := strings.Fields(line)

	nums := make([]int, len(fields))
	for i, f := range fields {
		nums[i], _ = strconv.Atoi(f)
	}

	return nums
}

func parseInput(filename string) Almanac {
	buff, _ := os.ReadFile(filename)
	content := string(buff)

	almanac := Almanac{}

	groups := strings.Split(content, "\n\n")

	almanac.Mappings = make([]Mapping, len(groups)-1)

	for _, group := range groups {
		group := strings.TrimSpace(group)
		// Map the `seeds`
		if strings.HasPrefix(group, "seeds: ") {
			line := strings.TrimPrefix(group, "seeds: ")
			almanac.Seeds = parseIntArr(line)

			// Map the mappings
		} else {
			lines := strings.Split(group, "\n")

			headerStr := strings.TrimSuffix(lines[0], " map:")
			header := strings.Split(headerStr, "-to-")
			srcType, destType := header[0], header[1]
			ranges := make([]Range, len(lines)-1)

			value := Mapping{SrcType: srcType, DestType: destType, Ranges: ranges}

			for i, line := range lines[1:] {
				nums := parseIntArr(line)
				r := Range{SrcStart: nums[1], DestStart: nums[0], Length: nums[2]}
				value.Ranges[i] = r
			}

			almanac.Mappings = append(almanac.Mappings, value)
		}
	}

	return almanac
}

func mapValue(val int, mapping Mapping) int {
	// Find the range that the value lies in
	var currRange *Range

	for _, r := range mapping.Ranges {
		if r.SrcStart <= val && val < (r.SrcStart+r.Length) {
			currRange = &r
			break
		}
	}

	if currRange == nil {
		return val
	} else {
		offset := val - currRange.SrcStart
		return currRange.DestStart + offset
	}
}

func mapSrcArray(
	vals []int,
	mappings []Mapping,
	srcType string,
) (mappedVals []int, destType string) {
	var mapping Mapping

	for _, val := range mappings {
		if val.SrcType == srcType {
			mapping = val
		}
	}

	destType = mapping.DestType
	mappedVals = make([]int, len(vals))

	for i, val := range vals {
		mappedVals[i] = mapValue(val, mapping)
	}

	// return mappedVals, destType
	return
}

func solveA(almanac Almanac) int {
	currVals := almanac.Seeds
	currType := "seed"
	targetType := "location"

	for currType != targetType {
		currVals, currType = mapSrcArray(currVals, almanac.Mappings, currType)
	}

	return slices.Min(currVals)
}

func buildPipeline(mappings []Mapping, srcType string, destType string) []Mapping {
	pipeline := make([]Mapping, 0)

	currType := srcType

	for currType != destType {
		for i, m := range mappings {
			if m.SrcType == currType {
				pipeline = append(pipeline, m)
				currType = m.DestType
				break
			}

			if i == len(mappings)-1 {
				panic(fmt.Sprintln("Could not find mapping for", currType))
			}
		}
	}

	return pipeline
}

func processPipeline(val int, pipeline []Mapping) int {
	mappedVal := val

	for _, mapping := range pipeline {
		mappedVal = mapValue(mappedVal, mapping)
	}

	return mappedVal
}

func getTuples(val []int) [][]int {
	if len(val)%2 != 0 {
		panic("Array length must be even!")
	}

	tuples := make([][]int, len(val)/2)

	for i := range tuples {
		tuples[i] = []int{val[i*2], val[i*2+1]}
	}

	return tuples
}

func solveB(almanac Almanac) int {
	pipeline := buildPipeline(almanac.Mappings, "seed", "location")

	// --- TEST --- Are the stages correct?
	// stages := make([]string, len(pipeline))
	// for i, p := range pipeline {
	//   stages[i] = p.DestType
	// }
	// fmt.Println("seed", stages)

	// --- TEST --- Is the pipeline working? Can we replicate solnA?
	// finalVals := make([]int, len(almanac.Seeds))
	// for i, val := range almanac.Seeds {
	//   finalVal := processPipeline(val, pipeline)
	//   finalVals[i] = finalVal
	// }
	// fmt.Println("location", finalVals)

	seeds := almanac.Seeds
	seedGroups := getTuples(seeds)

	minLocation := -1
	for i, group := range seedGroups {
		start, length := group[0], group[1]

		// Pretty Progress
		bar := progressbar.NewOptions(
			length,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionSetDescription(
				fmt.Sprintf(
					"[cyan][%d/%d][reset]", i, len(seedGroups),
				),
			),
			progressbar.OptionShowCount(),
      progressbar.OptionThrottle(1 * time.Second),
			progressbar.OptionOnCompletion(func() { fmt.Println("") }),
			progressbar.OptionShowElapsedTimeOnFinish(),
		)

		for val := start; val < (start + length); val++ {
			bar.Add(1)
			currLocation := processPipeline(val, pipeline)

			if minLocation == -1 {
				minLocation = currLocation
			} else {
				minLocation = min(currLocation, minLocation)
			}
		}
	}

	return minLocation
}

func main() {
	almanac := parseInput("day5/in.txt")

	// Test input parsing...
	_, _ = json.MarshalIndent(almanac, "", "  ")
	// fmt.Println(string(s))

	solnA := solveA(almanac)
	fmt.Println("Solution A:", solnA)

	solnB := solveB(almanac)
	fmt.Println("Solution B:", solnB)
}
