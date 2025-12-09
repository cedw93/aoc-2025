package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type (
	rect struct {
		width  int
		height int
		area   int
	}
	tile struct {
		row int
		col int
	}
)

var (
	tiles      = []tile{}
	rectangles = make(map[string]rect)
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		tiles = append(tiles, tile{
			row: aToIIgnoreError(parts[0]),
			col: aToIIgnoreError(parts[1]),
		})
	}
}

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

// not using math.x as that is for float64, not point converting for known int types
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func area(w, h int) int {
	return w * h
}

// constructPerimeter connects each tile in the list to the next tile with a straight line
// of intermediate tiles, forming a closed perimeter. Each tile (red) is connected to the
// tile before and after it by a straight line of green tiles. The list wraps, so the first
// tile is also connected to the last tile. Consecutive tiles are always on the same row or
// same column, and the function moves one step at a time horizontally or vertically.
//
// Example: Given tiles [(0,0), (0,2), (2,2)], it traces:
//   - (0,0) -> (0,1) -> (0,2)  [straight line right to next red tile]
//   - (0,2) -> (1,2) -> (2,2)  [straight line down to next red tile]
//   - (2,2) -> (2,1) -> (2,0) -> (1,0) -> (0,0)  [wrapping back to start]
//
// Result: all tiles visited along this closed path form the perimeter.
func constructPerimeter(tiles []tile) map[tile]bool {
	perimeterTiles := make(map[tile]bool)
	current := tiles[0]

	for _, t := range tiles[1:] {
		for current != t {
			perimeterTiles[current] = true
			if current.row != t.row {
				if current.row < t.row {
					current = tile{current.row + 1, current.col}
				} else {
					current = tile{current.row - 1, current.col}
				}
			} else {
				if current.col < t.col {
					current = tile{current.row, current.col + 1}
				} else {
					current = tile{current.row, current.col - 1}
				}
			}
		}
	}

	// Return to starting tile to complete the perimeter
	for current != tiles[0] {
		perimeterTiles[current] = true
		if current.row != tiles[0].row {
			if current.row < tiles[0].row {
				current = tile{current.row + 1, current.col}
			} else {
				current = tile{current.row - 1, current.col}
			}
		} else {
			if current.col < tiles[0].col {
				current = tile{current.row, current.col + 1}
			} else {
				current = tile{current.row, current.col - 1}
			}
		}
	}

	return perimeterTiles
}

// isValidRectangle checks if a rectangle formed by two corner tiles is valid by ensuring
// the perimeter doesn't cross through it. If the perimeter cuts through the rectangle, it's
// no longer an intact rectangle. The function detects crossing by finding two adjacent tiles
// (one on the edge, one just inside) that are both part of the perimeter path.
//
// Example of INVALID rectangle (perimeter crosses the left edge):
//
//	. P P P .     P = perimeter tile
//	. P R R .     R = rectangle corner/edge
//	. P P R .     . = empty
//	. . R R .
//
// At the left edge, both the edge tile and the interior tile (one step right) are perimeter
// tiles, meaning the perimeter path crosses into the rectangle, dividing it.
//
// Example of VALID rectangle (perimeter goes around, not through):
//
//	P P P P P
//	P R R R P
//	P R R R P
//	P R R R P
//	P P P P P
//
// The perimeter surrounds the rectangle without crossing through its interior.
// Therefore if a tile that connects all 4 corners is deemed to be inside the rectangle, then it is no longer a valid rectangle
// Basically for each pair of tiles along each edge of the rectangle, only 1 can be in the perimeter or its invalid
func isValidRectangle(a, b tile, perimeter map[tile]bool) bool {
	topLeftRow := min(a.row, b.row)
	topLeftCol := min(a.col, b.col)
	bottomRightRow := max(a.row, b.row)
	bottomRightCol := max(a.col, b.col)

	// Walk Left Side
	for i := 1; i < bottomRightRow-topLeftRow; i++ {
		first := tile{topLeftRow + i, topLeftCol}
		second := tile{topLeftRow + i, topLeftCol + 1}
		if perimeter[first] && perimeter[second] {
			return false
		}
	}

	// Walk Right Side
	for i := 1; i < bottomRightRow-topLeftRow; i++ {
		first := tile{topLeftRow + i, bottomRightCol}
		second := tile{topLeftRow + i, bottomRightCol - 1}
		if perimeter[first] && perimeter[second] {
			return false
		}
	}

	// Walk Top Side
	for i := 1; i < bottomRightCol-topLeftCol; i++ {
		first := tile{topLeftRow, topLeftCol + i}
		second := tile{topLeftRow + 1, topLeftCol + i}
		if perimeter[first] && perimeter[second] {
			return false
		}
	}

	// Walk Bottom Side
	for i := 1; i < bottomRightCol-topLeftCol; i++ {
		first := tile{bottomRightRow, topLeftCol + i}
		second := tile{bottomRightRow - 1, topLeftCol + i}
		if perimeter[first] && perimeter[second] {
			return false
		}
	}

	return true
}

func bothParts() (int, int) {
	largestArea := -1
	largestValidArea := -1
	perimeter := constructPerimeter(tiles)
	for i := 0; i < len(tiles)-1; i++ {
		tileOne := tiles[i]
		for j := i + 1; j < len(tiles); j++ {
			tileTwo := tiles[j]
			if _, exists := rectangles[fmt.Sprintf("%d,%d-%d,%d", tileOne.row, tileOne.col, tileTwo.row, tileTwo.col)]; !exists {
				// Hacky +1 here because coordinates are inclusive
				// e.g. (1,1) to (2,2) is width 2, height 2
				// not width 1, height 1
				// so we add one to each dimension
				// to get the correct area
				w := abs(tileOne.col-tileTwo.col) + 1
				h := abs(tileOne.row-tileTwo.row) + 1
				area := area(w, h)
				if area > largestArea {
					// fmt.Println("Largest area is between", tileOne, tileTwo)
					largestArea = area
				}

				if isValidRectangle(tileOne, tileTwo, perimeter) && area > largestValidArea {
					fmt.Println("Largest Valid Area", area)
					largestValidArea = area
				}
				rectangles[fmt.Sprintf("%d,%d-%d,%d", tileOne.row, tileOne.col, tileTwo.row, tileTwo.col)] = rect{
					width:  w,
					height: h,
					area:   area,
				}
			}
		}
	}
	return largestArea, largestValidArea
}

func main() {
	fmt.Println("!!! Warning !!!")
	fmt.Println("!!! Part 2 is slow on the real input !!!")
	fmt.Println("!!!!!!!!!!!!!!!!")
	partOne, partTwo := bothParts()
	println("Part One:", partOne)
	println("Part Two:", partTwo)
}
