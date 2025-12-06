package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	lines      []string
	dataLines  []string
	operations []string
)

// multiply multiplies together all numbers in a slice.
// Returns nil if the slice is empty.
func multiply(listOfNums []int) *int {
	if len(listOfNums) == 0 {
		return nil
	}
	result := listOfNums[0]
	for _, num := range listOfNums[1:] {
		result *= num
	}
	return &result
}

// calculateSum applies a sequence of operations column-wise and returns the summed result.
// For each list of numbers in list, the corresponding operator in
// `operators` determines how they are reduced:
//   - "*" multiplies all numbers together
//   - "+" sums all numbers together
func calculateSum(listOfNums [][]int, operators []string) int {
	results := []int{}
	for i, nums := range listOfNums {
		if i < len(operators) {
			operator := operators[i]
			if operator == "*" {
				product := multiply(nums)
				if product != nil {
					results = append(results, *product)
				}
			} else if operator == "+" {
				sum := 0
				for _, num := range nums {
					sum += num
				}
				results = append(results, sum)
			}
		}
	}

	total := 0
	for _, result := range results {
		total += result
	}
	return total
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	dataLines = lines[:len(lines)-1]
	operations = strings.Fields(lines[len(lines)-1])
}

// transforms horizontal rows of integers into vertical columns.
// for example, given input and input lines of ["123 328 51", "45 64 387", "6 98 215"],
// it would return [[123 45 6], [328 64 98], [51 387 215]]
// These can then be processed column-wise but it ignores the formatting of number in the input
func partOne() int {
	numbersHorizontal := [][]int{}
	for _, line := range dataLines {
		row := []int{}
		for _, x := range strings.Fields(line) {
			num, _ := strconv.Atoi(x)
			row = append(row, num)
		}
		numbersHorizontal = append(numbersHorizontal, row)
	}

	numColumns := len(numbersHorizontal[0])
	numbersVertical := [][]int{}

	for index := 0; index < numColumns; index++ {
		column := []int{}
		for _, row := range numbersHorizontal {
			column = append(column, row[index])
		}
		numbersVertical = append(numbersVertical, column)
	}

	return calculateSum(numbersVertical, operations)
}

// partTwo extracts numbers by reading the input grid column-by-column from right to left.
// Each column of characters is treated as a vertical string. These strings are reversed
// per row and stitched together to form integer values. Blank columns indicate separation.
//
// Example: given input lines:
//
//	"123 328  51 64 "
//	" 45 64  387 23 "
//	"  6 98  215 314"
//
// Reading right-to-left produces groups: [[4, 431, 623], [175, 581, 32], [8, 248, 369], [356, 24, 1]]
// this is because the final column for example, right to left by column is:
// 4 (4 from " 314") + 431 (4 from " 64 ") + 23 (3 from " 23 ") + 1 (from "123 ") and then the same for 623
func partTwo() int {
	// part 2 says it needs to processed right to left to reverse the operators
	reversedOperators := make([]string, len(operations))
	for i, j := 0, len(operations)-1; i < j; i, j = i+1, j-1 {
		reversedOperators[i], reversedOperators[j] = operations[j], operations[i]
	}

	numPositions := len(dataLines[0])

	numbers := []int{}
	numbersLeftToRight := [][]int{}

	for index := 0; index < numPositions; index++ {
		// Collect characters from each row at this column position (measured from the right)
		num := ""
		for _, row := range dataLines {
			reversed := reverseString(row)
			if index < len(reversed) {
				num += string(reversed[index])
			}
		}

		cleaned := strings.TrimSpace(num)
		if cleaned == "" {
			// Empty column → end of group
			numbersLeftToRight = append(numbersLeftToRight, numbers)
			numbers = []int{}
		} else if index == numPositions-1 {
			// Last column → finalize group
			val, _ := strconv.Atoi(cleaned)
			numbers = append(numbers, val)
			numbersLeftToRight = append(numbersLeftToRight, numbers)
		} else {
			val, _ := strconv.Atoi(cleaned)
			numbers = append(numbers, val)
		}
	}

	fmt.Println(numbersLeftToRight)

	return calculateSum(numbersLeftToRight, reversedOperators)

}

func main() {
	println("Part One:", partOne())
	println("Part Two:", partTwo())
}
