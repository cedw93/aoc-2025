package main

import (
	"bufio"
	"os"
)

var grid [][]cell

type cell struct {
	row           int
	col           int
	val           rune
	adjacentRolls int
	accessible    bool
}

type directionMap struct {
	rowOffset int
	colOffset int
	label     string
}

var directions = [...]directionMap{
	{0, 1, "right"},
	{1, 0, "down"},
	{0, -1, "left"},
	{-1, 0, "up"},
	{-1, -1, "upLeft"},
	{-1, 1, "upRight"},
	{1, -1, "downLeft"},
	{1, 1, "downRight"},
}

const (
	RollOfPaper = '@'
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	rowCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]cell, 0)
		col := 0
		for _, r := range []rune(line) {
			row = append(row, cell{
				row:           rowCount,
				col:           col,
				val:           r,
				adjacentRolls: 0,
				accessible:    false,
			})
			col++
		}
		rowCount++
		grid = append(grid, row)
	}
}

func withinGrindBoundary(row, col int, g [][]cell) bool {
	return row > -1 && col > -1 && row < len(g) && col < len(g[row])
}

func deepCopyGrid(original [][]cell) [][]cell {
	copied := make([][]cell, len(original))
	for i := range original {
		copied[i] = make([]cell, len(original[i]))
		copy(copied[i], original[i])
	}
	return copied
}

func calcAdjacent(c cell, g [][]cell) int {
	adjacent := 0
	for _, direction := range directions {
		offsetRow := c.row + direction.rowOffset
		offsetColumn := c.col + direction.colOffset
		if !withinGrindBoundary(offsetRow, offsetColumn, g) {
			continue
		}
		if g[offsetRow][offsetColumn].val == RollOfPaper {
			adjacent++
		}
	}
	return adjacent
}

func partOne() int {
	accessibleCount := 0
	for _, r := range grid {
		for _, c := range r {
			cell := grid[c.row][c.col]
			adjacent := calcAdjacent(cell, grid)
			if cell.val == RollOfPaper && adjacent < 4 {
				accessibleCount++
			}
			// count accessible rolls of paper, i.e. those with less than 4 adjacent rolls and the cell itself a roll of papers

		}
	}
	return accessibleCount
}

func removalCandidates(g [][]cell) []cell {
	result := []cell{}
	for _, r := range g {
		for _, c := range r {
			cell := g[c.row][c.col]
			adjacent := calcAdjacent(cell, g)
			if cell.val == RollOfPaper && adjacent < 4 {
				result = append(result, cell)
			}
		}
	}
	return result
}

// Bit wasteful as it doesn't reuse the result from part one but it is what it is for now
func partTwo() int {
	currentGrid := deepCopyGrid(grid)
	candidates := removalCandidates(currentGrid)
	removed := 0

	for len(candidates) > 0 {
		removed += len(candidates)
		for _, cell := range candidates {
			currentGrid[cell.row][cell.col].val = 'x'

		}
		currentGrid = deepCopyGrid(currentGrid)
		candidates = removalCandidates(currentGrid)
	}
	return removed
}

func main() {
	println("Part One:", partOne())
	println("Part Two:", partTwo())
}
