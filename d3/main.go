package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type bank struct {
	batteries []int
	raw       []rune
}

var banks []bank

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		raw := []rune(line)
		batteries := make([]int, len(raw))

		for i, ch := range raw {
			batteries[i] = int(ch - '0')
		}

		banks = append(banks, bank{
			raw:       raw,
			batteries: batteries,
		})
	}
}

func maxIndex(candidates []int, remainingLength int) int {
	currentMax := -1
	idx := 0
	for i, candidate := range candidates {
		// ensure enough elements remain after choosing this candidate
		// as we need to pick remainingLength elements in total
		// for example, if we have
		// [0,1,2,3,4] and need to pick 3 elements
		// at index 2 (value 2) we can only pick two more elements (3,4)
		// so we cannot pick index 2 as 2 + 3 > 5 (len of candidates)
		// therefore in that scenario we can only pick between indices 0 and 1
		// which results in idx = 1 being returned
		if i+remainingLength <= len(candidates) && candidate > currentMax {
			currentMax = candidate
			idx = i
		}
	}
	return idx
}

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func (b bank) maxVoltage(expectedLength int) int {
	results := []rune{}
	offset := 0
	for i := expectedLength; i > 0; i-- {
		maxIdx := maxIndex(b.batteries[offset:], i)
		results = append(results, b.raw[maxIdx+offset])
		// since maxId is relative to the slice when maxId returns 0 its actually b.batteries[offset + maxId]
		// +1 just progresses past the found max for the next search as the slice indexing is inclusive of the start index
		offset = offset + maxIdx + 1
	}
	return aToIIgnoreError(string(results))
}

func partOne() int {
	result := 0
	for _, b := range banks {
		result += b.maxVoltage(2)
	}

	return result
}

func partTwo() int {
	result := 0
	for _, b := range banks {
		result += b.maxVoltage(12)

	}

	return result
}

func main() {

	partOne := partOne()
	partTwo := partTwo()
	fmt.Printf("PartOne invalid sum: %d\n", partOne)
	fmt.Printf("PartTwo invalid sum: %d\n", partTwo)

}
