package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/draffensperger/golp"
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
		joltage       []int
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

func joltageFromText(raw string) []int {
	joltage := []int{}
	raw = strings.Trim(raw, "{}")
	for _, val := range strings.Split(raw, ",") {
		joltage = append(joltage, aToIIgnoreError(val))
	}
	return joltage
}

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		targetString := parts[0]
		buttons := []button{}
		var joltage []int

		for i := 1; i < len(parts); i++ {
			part := parts[i]
			if strings.HasPrefix(part, "(") {
				buttons = append(buttons, buttonFromText(part))
			} else if strings.HasPrefix(part, "{") {
				joltage = joltageFromText(part)
			}
		}

		diagrams = append(diagrams, diagram{
			target:        stringToMask(targetString),
			buttons:       buttons,
			raw:           line,
			fewestPresses: math.MaxInt,
			targetString:  targetString,
			// account for brackets
			numIndicators: len(targetString) - 2,
			joltage:       joltage,
		})
	}
}

// solveForJoltage finds the minimum number of button presses needed to achieve the desired joltage
// for each indicator using linear programming. Each button press toggles specific indicators and
// increases their joltage by 1.
//
// Example: If we want joltage [3,5,4] for indicators [0,1,2] and have buttons:
//   - Button A toggles [0,1]
//   - Button B toggles [1,2]
//
// Then pressing A 3 times and B 2 times gives joltage [3,5,2], but we need [3,5,4].
// The LP solver finds the optimal combination: A=3, B=4 gives [3,7,4], which is incorrect.
// The correct solution minimizes total presses while satisfying all joltage constraints exactly.
func solveForJoltage(d *diagram) int {
	if len(d.joltage) == 0 {
		return 0
	}

	numButtons := len(d.buttons)
	numJoltages := len(d.joltage)

	lp := golp.NewLP(0, numButtons)
	lp.SetVerboseLevel(golp.NEUTRAL)

	// Objective for solve minimize total button presses
	objectiveCoeffs := make([]float64, numButtons)
	for i := range numButtons {
		objectiveCoeffs[i] = 1.0
	}
	lp.SetObjFn(objectiveCoeffs)

	// Set variable bounds: each button can be pressed 0 to 1000 times (integer)
	for i := range numButtons {
		lp.SetInt(i, true)
		lp.SetBounds(i, 0.0, 1000.0)
	}

	for i := 0; i < numJoltages; i++ {
		var entries []golp.Entry
		for j, btn := range d.buttons {
			if slices.Contains(btn.lights, i) {
				entries = append(entries, golp.Entry{Col: j, Val: 1.0})
			}
		}
		targetValue := float64(d.joltage[i])
		if err := lp.AddConstraintSparse(entries, golp.EQ, targetValue); err != nil {
			panic(err)
		}
	}

	// Solve the problem using linear programming library
	status := lp.Solve()

	if status != golp.OPTIMAL {
		return 0
	}

	// Get solution and sum up total presses
	solution := lp.Variables()
	totalPresses := 0
	for _, val := range solution {
		totalPresses += int(val + 0.5) // round to nearest integer
	}

	return totalPresses
}

func main() {
	partOne := 0
	for _, d := range diagrams {
		solve(&d)
		// fmt.Printf("%s can be solved in %d presses\n", d.targetString, d.fewestPresses)
		partOne += d.fewestPresses
	}

	println("Part One:", partOne)

	partTwo := 0
	for _, d := range diagrams {
		presses := solveForJoltage(&d)
		partTwo += presses
	}
	println("Part Two:", partTwo)
}
