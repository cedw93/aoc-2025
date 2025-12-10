package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type (
	button struct {
		lights []int
	}
	state struct {
		lights  uint16
		presses int
		parent  *state
	}
	// We bitmask the lights, my input only had at most 13 indicators so uint16 is sufficient
	diagram struct {
		target        uint16
		numIndicators int
		buttons       []button
		raw           string
		fewestPresses int
		targetString  string
		visited       map[uint16]int
	}
)

var (
	diagrams = []diagram{}
)

const (
	On  = '#' // treat this a 1 in the mask
	Off = '.'
)

func (d diagram) String() string {
	return fmt.Sprintf("Diagram targeting %016b with buttons %v", d.target, d.buttons)
}

func (b button) updateLights(current uint16, numIndicators int) uint16 {
	newLights := current
	for _, light := range b.lights {
		newLights = toggleBit(newLights, light, numIndicators)
	}
	return newLights
}

func solve(d *diagram) {
	d.visited = make(map[uint16]int)
	dfs(d, 0, 0)
}

func dfs(d *diagram, lights uint16, presses int) {
	if presses >= d.fewestPresses {
		return
	}

	// keyed by lights, for example 0001001010 so if we've seen this state in fewer presses, skip
	if prevPresses, seen := d.visited[lights]; seen && prevPresses <= presses {
		return
	}

	d.visited[lights] = presses

	if lights == d.target {
		if presses < d.fewestPresses {
			d.fewestPresses = presses
		}
		return
	}

	for _, b := range d.buttons {
		newLights := b.updateLights(lights, d.numIndicators)
		dfs(d, newLights, presses+1)
	}
}

func setBit(mask uint16, i int, numIndicators int) uint16 {
	return mask | (1 << (numIndicators - 1 - i))
}

func toggleBit(mask uint16, i int, numIndicators int) uint16 {
	return mask ^ (1 << (numIndicators - 1 - i))
}

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func stringToMask(raw string) uint16 {
	mask := uint16(0)
	raw = strings.Trim(raw, "[]")
	runes := []rune(raw)
	for i, r := range runes {
		if r == On {
			mask = setBit(mask, i, len(runes))
		}
	}
	return mask
}

func buttonFromText(raw string) button {
	lights := []int{}
	raw = strings.Trim(raw, "()")
	for _, light := range strings.Split(raw, ",") {
		lights = append(lights, aToIIgnoreError(light))
	}
	return button{lights: lights}
}

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		targetString := parts[0]
		buttons := []button{}
		for _, buttonText := range parts[1 : len(parts)-1] {
			buttons = append(buttons, buttonFromText(buttonText))
		}
		diagrams = append(diagrams, diagram{
			target:        stringToMask(targetString),
			buttons:       buttons,
			raw:           line,
			fewestPresses: math.MaxInt,
			targetString:  targetString,
			// account for brackets
			numIndicators: len(targetString) - 2,
		})
	}
}

func main() {
	partOne := 0
	for _, d := range diagrams {
		solve(&d)
		// fmt.Printf("%s can be solved in %d presses\n", d.targetString, d.fewestPresses)
		partOne += d.fewestPresses
	}

	println("Part One:", partOne)
	println("Part Two, No idea how to do!:")
}
