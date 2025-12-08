package main

import (
	"bufio"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type (
	box struct {
		x int
		y int
		z int
	}
	pair struct {
		boxOne box
		boxTwo box
		dist   float64
	}
)

var (
	boxes []box
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		boxes = append(boxes, box{
			x: aToIIgnoreError(parts[0]),
			y: aToIIgnoreError(parts[1]),
			z: aToIIgnoreError(parts[2]),
		})
	}
}

// https://en.wikipedia.org/wiki/Euclidean_distance
func distance(a, b box) float64 {
	dx := a.x - b.x
	dy := a.y - b.y
	dz := a.z - b.z
	return math.Sqrt(float64(dx*dx + dy*dy + dz*dz))
}

// remove removes a target pointer from a slice of pointers.
// Example: if circuits = [&a, &b, &c] and target = &b, returns [&a, &c]
// This uses pointer equality (same memory address), not value equality.
func remove(circuits []*map[box]bool, target *map[box]bool) []*map[box]bool {
	for i, circuit := range circuits {
		if circuit == target {
			return append(circuits[:i], circuits[i+1:]...)
		}
	}
	return circuits
}

func partOne(batchSize int) int {
	pairs := []pair{}
	for i := 0; i < len(boxes); i++ {
		boxOne := boxes[i]
		for j := i + 1; j < len(boxes); j++ {
			boxTwo := boxes[j]
			pairs = append(pairs, pair{
				boxOne: boxOne,
				boxTwo: boxTwo,
				dist:   distance(boxOne, boxTwo),
			})
		}
	}
	// Sort pairs by distance as per the problem statement
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].dist < pairs[j].dist
	})
	// Using pointers so we can manipulate circuits more easily
	circuitMap := make(map[box]*map[box]bool)
	foundCircuits := []*map[box]bool{}

	for i := 0; i < batchSize; i++ {
		pair := pairs[i]
		circuitOne, existsOne := circuitMap[pair.boxOne]
		circuitTwo, existsTwo := circuitMap[pair.boxTwo]

		// None of the boxes exist in a circuit yet
		if !existsOne && !existsTwo {
			circuit := make(map[box]bool)
			circuit[pair.boxOne] = true
			circuit[pair.boxTwo] = true
			foundCircuits = append(foundCircuits, &circuit)
			circuitMap[pair.boxOne] = &circuit
			circuitMap[pair.boxTwo] = &circuit
			continue
		}
		// boxOne exists in a circuit but boxTwo does not
		// Add boxTwo to boxOne's circuit
		if existsOne && !existsTwo {
			circuit := circuitMap[pair.boxOne]
			(*circuit)[pair.boxTwo] = true
			circuitMap[pair.boxTwo] = circuit
			continue
		}

		if !existsOne && existsTwo {
			circuit := circuitMap[pair.boxTwo]
			(*circuit)[pair.boxOne] = true
			circuitMap[pair.boxOne] = circuit
			continue
		}

		// Both boxes exist in different circuits, need to merge
		if existsOne && existsTwo && circuitOne != circuitTwo {
			// Merge the two circuits
			circuit := make(map[box]bool)
			for b := range *circuitOne {
				circuit[b] = true
			}
			for b := range *circuitTwo {
				circuit[b] = true
			}

			for b := range circuit {
				circuitMap[b] = &circuit
			}
			// Remove the now merged circuit then add the merged one
			foundCircuits = remove(foundCircuits, circuitOne)
			foundCircuits = remove(foundCircuits, circuitTwo)
			foundCircuits = append(foundCircuits, &circuit)
		}
	}

	sort.Slice(foundCircuits, func(i, j int) bool {
		return len(*foundCircuits[i]) > len(*foundCircuits[j])
	})

	return len(*foundCircuits[0]) * len(*foundCircuits[1]) * len(*foundCircuits[2])
}

// Pretty wasteful, basically the same as part one but ignores batch size and goes until all boxes are connected in 1 circuit
// could reuse this for both answers but runs quick enough to not care
func partTwo() int {
	pairs := []pair{}
	for i := 0; i < len(boxes); i++ {
		boxOne := boxes[i]
		for j := i + 1; j < len(boxes); j++ {
			boxTwo := boxes[j]
			pairs = append(pairs, pair{
				boxOne: boxOne,
				boxTwo: boxTwo,
				dist:   distance(boxOne, boxTwo),
			})
		}
	}
	// Sort pairs by distance as per the problem statement
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].dist < pairs[j].dist
	})
	// Using pointers so we can manipulate circuits more easily
	circuitMap := make(map[box]*map[box]bool)
	foundCircuits := []*map[box]bool{}

	for i := 0; i < len(pairs); i++ {
		pair := pairs[i]
		// Get existing circuits for both boxes, if any, could end up with nil, false for example
		circuitOne, existsOne := circuitMap[pair.boxOne]
		circuitTwo, existsTwo := circuitMap[pair.boxTwo]

		// None of the boxes exist in a circuit yet
		// Adding together in a new circuit is quite easy
		// make a new circuit and flag them as being part of it (true)
		// simply append it to foundCircuits for tracking
		// then update the global map of which circuit each box belongs to
		if !existsOne && !existsTwo {
			circuit := make(map[box]bool)
			circuit[pair.boxOne] = true
			circuit[pair.boxTwo] = true
			foundCircuits = append(foundCircuits, &circuit)
			circuitMap[pair.boxOne] = &circuit
			circuitMap[pair.boxTwo] = &circuit

		} else if existsOne && !existsTwo {
			// boxOne exists in a circuit but boxTwo does not
			// Add boxTwo to boxOne's circuit
			// This is the same for the opposite case below
			// simply get the existing circuit pointer from the map for the box that has a circuit
			// set the flag to true for box 2 now its going to be in a circuit
			// then finally update the global map to point boxTwo to the same circuit as boxOne
			circuit := circuitMap[pair.boxOne]
			(*circuit)[pair.boxTwo] = true
			circuitMap[pair.boxTwo] = circuit

		} else if !existsOne && existsTwo {
			// same as above
			circuit := circuitMap[pair.boxTwo]
			(*circuit)[pair.boxOne] = true
			circuitMap[pair.boxOne] = circuit

		} else if existsOne && existsTwo && circuitOne != circuitTwo {
			// Both boxes exist in different circuits, need to merge
			// Merge the two circuits, this is easy too
			// make a new merged circuit, empty to start
			merged := make(map[box]bool)
			// Go through all boxes in first circuit and flag them as part of merged
			for b := range *circuitOne {
				merged[b] = true
			}
			// Go through all boxes in second circuit and flag them as part of merged
			for b := range *circuitTwo {
				merged[b] = true
			}
			// go through all the boxes in merged (both circuits now) and then the location in the global map to the merged circuit
			for b := range merged {
				circuitMap[b] = &merged
			}
			// Remove the now merged circuit then add the merged one
			foundCircuits = remove(foundCircuits, circuitOne)
			foundCircuits = remove(foundCircuits, circuitTwo)
			foundCircuits = append(foundCircuits, &merged)
		}

		// Check after every pair if we've connected all boxes
		if len(foundCircuits) == 1 && len(*foundCircuits[0]) == len(boxes) {
			return pair.boxOne.x * pair.boxTwo.x
		}
	}

	return -1
}

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func main() {
	//sample
	// batchSize := 10
	// real input
	batchSize := 1000

	println("Part One:", partOne(batchSize))
	println("Part Two:", partTwo())

}
