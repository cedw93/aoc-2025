package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type IdRange struct {
	start           int
	end             int
	invalid         []int
	invalidMultiple []int
}

var ranges []*IdRange

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		for _, seq := range strings.Split(scanner.Text(), ",") {
			parts := strings.Split(seq, "-")
			ranges = append(ranges, &IdRange{
				start:           aToIIgnoreError(parts[0]),
				end:             aToIIgnoreError(parts[1]),
				invalid:         []int{},
				invalidMultiple: []int{},
			})
		}
	}
}

func (i *IdRange) String() string {
	return fmt.Sprintf("%d to %d has %d invalid: %v", i.start, i.end, len(i.invalid), i.invalid)
}

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func calcInvalid() (int, int) {
	partOne := 0
	partTwo := 0
	for _, r := range ranges {
		for candidate := r.start; candidate <= r.end; candidate++ {
			candidateAsString := strconv.Itoa(candidate)
			if checkRepeatedInvalid(candidateAsString) {
				r.invalidMultiple = append(r.invalidMultiple, candidate)
				partTwo += candidate
			}
			if len(candidateAsString)%2 == 0 {
				mid := len(candidateAsString) / 2
				firstHalf := candidateAsString[:mid]
				secondHalf := candidateAsString[mid:]
				if firstHalf == secondHalf {
					r.invalid = append(r.invalid, candidate)
					partOne += candidate
				}
			}
		}
	}
	return partOne, partTwo
}

func checkRepeatedInvalid(id string) bool {
	idLength := len(id)

	// Try every possible block length from 1 up to n/2. We only need to check up to n/2 as there can't be repeats longer than half the string length.
	for blockLen := 1; blockLen <= idLength/2; blockLen++ {
		// The string length must be a multiple of blockLen or there is no duplicate pattern.
		if idLength%blockLen != 0 {
			continue
		}

		pattern := id[0:blockLen]
		allMatch := true

		// Check that every subsequent block matches the first one.
		// e.g. if pattern is "12" and blockLen is 2, check id[2:4], id[4:6], etc that they are all "12"
		for i := blockLen; i < idLength; i += blockLen {
			if id[i:i+blockLen] != pattern {
				allMatch = false
				break
			}
		}

		if allMatch {
			return true
		}
	}

	return false
}

func main() {
	partOne, partTwo := calcInvalid()
	fmt.Printf("PartOne invalid sum: %d\n", partOne)
	fmt.Printf("PartTwo invalid sum: %d\n", partTwo)
}
