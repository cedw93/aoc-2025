package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var rotations []int

const (
	dialMax   = 100
	dialStart = 50
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		runes := []rune(scanner.Text())
		dir := runes[0]
		runes = runes[1:]
		if dir == 'L' {
			runes = []rune("-" + string(runes))
		}

		rotations = append(rotations, aToIIgnoreError(string(runes)))
	}
}

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

// equivalent to Python's divmod but works for negative numbers as expected
func divideAndModulo(a, b int) (int, int) {
	quotient := a / b
	rem := a % b

	if rem != 0 && ((rem > 0) != (b > 0)) {
		quotient--
		rem += b
	}
	return quotient, rem
}

func main() {
	current := dialStart
	zeroCount := 0
	throughZero := 0
	for _, rotation := range rotations {

		if rotation < 0 {
			div, remainder := divideAndModulo(rotation, -dialMax)
			throughZero += div
			if current != 0 && current+remainder <= 0 {
				throughZero++
			}
		} else {
			div, remainder := divideAndModulo(rotation, dialMax)
			throughZero += div
			if current+remainder >= dialMax {
				throughZero++
			}
		}

		current = (current + rotation) % dialMax
		// Golang % can return negative so this is a hack to make it positive
		if current < 0 {
			current += dialMax
		}
		if current == 0 {
			zeroCount++
		}
	}
	fmt.Println("Times at Zero:", zeroCount)
	fmt.Println("Times through Zero:", throughZero)
}
