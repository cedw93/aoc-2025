package main

import (
	"bufio"
	"os"
	"strings"
)

type (
	// a graph of node names -> neighboring node names
	// example ["abc"] -> ["def, "xyz"]
	graph map[string][]string
	// For part 2, tracking state of path not just the path itself
	// This cares about which nodes have been seen up to this node
	pathState struct {
		id      string
		usedFFT bool
		usedDAC bool
	}
)

var (
	currentGraph graph = make(graph)
)

const (
	StartNode            = "you"
	EndNode              = "out"
	ServerRack           = "svr"
	DigiToAnoConverter   = "dac"
	FastFourierTransform = "fft"
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		nodeId := parts[0]
		currentGraph[nodeId] = strings.Fields(parts[1])
	}
}

// Very dumb way, bit like a DFS but summing everything
// Might loop forever if theres no end in input but assuming AoC inputs are always
// valid and possible (for this problem at least)
func partOne(currentNodeId string, cache map[string]int) int {
	paths := 0
	// Input is not really large enough to cache but doing it anyway
	if val, ok := cache[currentNodeId]; ok {
		return val
	}

	for _, next := range currentGraph[currentNodeId] {
		if next == EndNode {
			return 1
		}
		paths += partOne(next, cache)
	}
	cache[currentNodeId] = paths
	return paths
}

// This now needs to track which required nodes have been visited, could reuse part1 for both but
// duplicating for clarity/simplicity
func partTwo(currentNodeId string, visited map[pathState]int, usedFft bool, usedDac bool) int {
	result := 0
	// If everything matches, we've seen this exact state before.
	// Just getting to currentNodeId is not enough, we need to know if we've used the required nodes as well
	// currentNode="x" usedFft=true usedDac=false is different to currentNode="x" usedFft=true usedDac=true
	currentState := pathState{id: currentNodeId, usedFFT: usedFft, usedDAC: usedDac}
	if value, ok := visited[currentState]; ok {
		return value
	}

	// We are at the end, check if we've used both required nodes otherwise it's not a valid path
	if currentNodeId == EndNode {
		if usedFft && usedDac {
			return 1
		}
		return 0
	}

	for _, next := range currentGraph[currentNodeId] {
		if currentNodeId == FastFourierTransform {
			usedFft = true
		}
		if currentNodeId == DigiToAnoConverter {
			usedDac = true
		}
		result += partTwo(next, visited, usedFft, usedDac)
	}
	visited[currentState] = result

	return result
}

func main() {
	println("Part One:", partOne(StartNode, make(map[string]int)))
	println("Part Two:", partTwo(ServerRack, make(map[pathState]int), false, false))
}
