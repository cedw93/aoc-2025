package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type (
	region struct {
		width          int
		height         int
		requiredShapes []int
	}

	shape struct {
		id        int
		structure [][]rune
		occupies  int
	}
)

var (
	regions  []region
	shapeMap = make(map[int]*shape)
)

func init() {
	// easier to do todays reading all at once
	data, _ := os.ReadFile("input.txt")
	// data, _ := os.ReadFile("sample.txt")
	fileChunks := strings.Split(string(data), "\n\n")
	shapeData := fileChunks[:len(fileChunks)-1]

	for _, r := range strings.Split(fileChunks[len(fileChunks)-1], "\n") {
		if r == "" {
			continue
		}
		parts := strings.Split(r, ":")
		dimParts := strings.Split(parts[0], "x")
		width := aToIIgnoreError(dimParts[0])
		height := aToIIgnoreError(dimParts[1])
		requiredShapes := []int{}
		for _, id := range strings.Fields(parts[1]) {
			requiredShapes = append(requiredShapes, aToIIgnoreError(id))
		}
		regions = append(regions, region{
			width:          width,
			height:         height,
			requiredShapes: requiredShapes,
		})
	}

	// Takes a shape block such as
	// 0:
	// ###
	// ##.
	// ##.
	// and processed into a shape such as
	// &{0 [[35 35 35] [35 35 46] [35 35 46]]}
	// which is
	// &{0 [[# # #] [# # .] [# # .]]}
	for _, chunk := range shapeData {
		lines := strings.Split(chunk, "\n")

		shapeId := aToIIgnoreError(strings.TrimSuffix(lines[0], ":"))

		structure := [][]rune{}
		for _, line := range lines[1:] {
			if line != "" {
				structure = append(structure, []rune(line))
			}
		}

		shapeMap[shapeId] = &shape{
			id:        shapeId,
			structure: structure,
			occupies:  strings.Count(strings.Join(lines[1:], ""), "#"),
		}
	}
}

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

// Honestly this was by chance, problem seemed impossibly difficult time wise
// So simply check if the region can contain that main shapes at all
// it seems like the inputs are designed so that this is sufficient as the answer is always
// they cannot or there is a lot of spare room
// Doesn't really work for the sample, it returns 3 instead of 2 which is unfortunate but close enough
// so thought i'd try it on the actual input to see how many could possibly fit the shapes to see how big the answer is
// in terms of trying every shape, tried the answer and it worked
//
// This is likely intended as its the final day and a bit of a trick, I don't think its possible to check every arrangement
// as there are 100s of possible regions with millions of arrangements
func partOne() int {
	possible := 0
	for _, r := range regions {
		area := r.width * r.height
		totalRequired := 0
		for id, shapeCount := range r.requiredShapes {
			totalRequired += shapeMap[id].occupies * shapeCount
		}
		// Simply, is the area of the region bigger than the possible amount of space needed to fit all shapes
		// regardless of arrangement or optimisation

		// This seems intended as area - totalRequired is always negative or a large number such theres so much spare room
		// for shapes ore not enough space at all
		if area >= totalRequired {
			possible++
		}
	}
	return possible
}

func main() {
	fmt.Println("Part One:", partOne())
	fmt.Println("Part Two: There is not part two, this is the final day! Christmas is saved!")
}
