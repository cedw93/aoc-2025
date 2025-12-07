package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	grid         [][]cell
	possibleWays []int
)

type cell struct {
	row int
	col int
	val rune
}

const (
	beam     = '|'
	blank    = '.'
	Start    = 'S'
	splitter = '^'
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	rowIdx := 0
	for scanner.Scan() {
		row := []cell{}
		for i, rune := range []rune(scanner.Text()) {
			c := cell{
				row: rowIdx,
				col: i,
				val: rune,
			}
			row = append(row, c)
		}
		grid = append(grid, row)
		rowIdx++
	}
	possibleWays = make([]int, len(grid[0]))
}

// splitsAndRealities simulates how beams/paths propagate down the grid and
// counts split events and resulting realities.
//
// It treats each column as a counter of "ways" to reach that column. For
// each row (top-to-bottom):
//   - If a cell is `Start` ('S'), that column is initialized with one way.
//   - If a cell is a `splitter` ('^') and the column currently has >0 ways,
//     the splitter is counted (increments `splits`), all ways at that column
//     are transferred to the immediate left and right columns (added to
//     their counters), and the current column's ways are cleared. Other
//     cells are ignored for path propagation.
//
// After processing every row, the function sums the remaining ways across
// columns to compute the total number of distinct resulting "realities".
// It returns (splits, realities).
//
// Note: We never need to calculate or update the state of the grid, we only need how many times we split
// and how many times we are in a column at the end
func splitsAndRealities(input [][]cell) (int, int) {
	splits := 0
	for _, row := range input {
		for j, c := range row {
			if c.val == Start {
				// Starting point, exactly one way to get there
				possibleWays[j] = 1
			} else if c.val == splitter && possibleWays[j] > 0 {
				// only count a split if there's a way to get here
				// there are more splitters than ways to get to them, some are unreachable
				splits++
				possibleWays[j-1] += possibleWays[j]
				possibleWays[j+1] += possibleWays[j]
				// reset to zero to avoid double counting for exmaple, let's say
				// ways = [0, 1, 0] and we hit a splitter at index 1
				// we want to move that 1 way to index 0 and 2, then set index 1 to 0 as we don't move to a splitter
				// we split once, now we have ways = [1, 0, 1] and the rest can continue
				possibleWays[j] = 0
			}
		}
	}

	realities := 0
	// Sum up the remianing possible ways across all columns which should be the number of total realities
	for _, v := range possibleWays {
		realities += v
	}
	return splits, realities
}

func main() {
	splits, realities := splitsAndRealities(grid)
	fmt.Println("Part One:", splits)
	fmt.Println("Part Two:", realities)
}
