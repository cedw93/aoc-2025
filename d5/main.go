package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type freshRange struct {
	start int
	end   int
	raw   string
}

var freshRanges []freshRange
var ingredients []int

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			freshRanges = append(freshRanges, freshRange{
				start: aToIIgnoreError(parts[0]),
				end:   aToIIgnoreError(parts[1]),
				raw:   line,
			})
			continue
		} else {
			ingredients = append(ingredients, aToIIgnoreError(line))
		}
	}
	sort.Slice(freshRanges, func(i, j int) bool {
		return freshRanges[i].start < freshRanges[j].start
	})
}

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func (fr freshRange) String() string {
	return fmt.Sprintf("Fresh range from %d to %d (%s)", fr.start, fr.end, fr.raw)
}

func partOne() int {
	freshCount := 0
	for _, ingredient := range ingredients {
		for _, fr := range freshRanges {
			if ingredient >= fr.start && ingredient <= fr.end {
				// fmt.Println(ingredient, "is fresh")
				freshCount++
				// break to prevent double counting if an ingredient falls into multiple ranges
				break
			}
		}
	}
	return freshCount
}

// Can't brute force this one... need to merge ranges first
// sort the slice (already sorted in init)
// then iterate through, merging overlapping or contiguous ranges
// for example:
// [1-3], [2-4], [6-8], [7-10]
// becomes
// [1-4], [6-10]
// then we can just sum the lengths of the merged ranges using (end - start + 1). The +1 is because the ranges are inclusive
func partTwo() int {
	mergedRanges := []freshRange{}
	result := 0
	// Since slice is sorted this will be the smallest start value
	current := freshRanges[0]
	for i := 1; i < len(freshRanges); i++ {
		next := freshRanges[i]
		if next.start <= current.end+1 {
			// overlapping or contiguous ranges, merge them
			if next.end > current.end {
				current.end = next.end
			}
		} else {
			// no overlap, add current to merged and move to next one
			mergedRanges = append(mergedRanges, current)
			current = next
		}
	}
	// add the last range
	mergedRanges = append(mergedRanges, current)

	for _, mergedRanges := range mergedRanges {
		result += (mergedRanges.end - mergedRanges.start + 1)
	}
	return result
}

func main() {
	println("Part One:", partOne())
	println("Part Two:", partTwo())
}
