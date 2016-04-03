package main

import (
	"bufio"
	"fmt"
	"os"
)

const DIM = 9

// Board represents a sudoku board
type Board struct {
	cell         [][]int
	rowRemaining []int
	colRemaining []int
	remaining    int
	backtracks   int
}

// NewBoard creates an empty sudoku board
func NewBoard() *Board {
	b := &Board{}
	// Build empty (zero) cell matrix
	b.cell = make([][]int, DIM)
	for i := range b.cell {
		b.cell[i] = make([]int, DIM)
	}
	// Init remaining counters
	b.rowRemaining = make([]int, DIM)
	for i := range b.rowRemaining {
		b.rowRemaining[i] = DIM
	}
	b.colRemaining = make([]int, DIM)
	for i := range b.colRemaining {
		b.colRemaining[i] = DIM
	}
	b.remaining = DIM * DIM
	return b
}

// String formats the board for human consumption
func (b *Board) String() string {
	var result = "    1 2 3 4 5 6 7 8 9\n"
	for i, row := range b.cell {
		result += fmt.Sprintf("%v: %v\n", i+1, row)
	}
	result += fmt.Sprintf("Remaining: %v, Backtracks: %v", b.remaining, b.backtracks)
	return result
}

// ValidSolution is true if remaining == 0
func (b *Board) ValidSolution() bool {
	return b.remaining == 0
}

// MakeMove adds a number to the board, row and col indices are 0 based
func (b *Board) MakeMove(row, col, val int) {
	if b.cell[row][col] == 0 && val != 0 {
		b.remaining--
		b.rowRemaining[row]--
		b.colRemaining[col]--
	}
	b.cell[row][col] = val
}

// UnmakeMove removes a number from the board, row and col indices are 0 based
func (b *Board) UnmakeMove(row, col int) {
	if b.cell[row][col] != 0 {
		b.remaining++
		b.rowRemaining[row]++
		b.colRemaining[col]++
		b.cell[row][col] = 0
	}
	b.backtracks++
}

// NextEmptyCell tells our solver which cell to work on next
func (b *Board) NextEmptyCell() (row, col int) {
	row = -1
	col = -1
	rmin := DIM + 1
	cmin := DIM + 1
	// Look for most constrained row
	for i, rem := range b.rowRemaining {
		if 0 < rem && rem < rmin {
			row, rmin = i, rem
		}
	}
	// Look for most constrained empty column in this row
	for i, rem := range b.colRemaining {
		if b.cell[row][i] == 0 && 0 < rem && rem < cmin {
			col, cmin = i, rem
		}
	}
	// Return row, col
	return
}

// CellCandidates returns a list of legal moves for specified cell
func (b *Board) CellCandidates(row, col int) []bool {
	if row < 0 || DIM < row {
		panic(fmt.Sprintf("Invalid row passed: %v", row))
	}
	if col < 0 || DIM < col {
		panic(fmt.Sprintf("Invalid col passed: %v", col))
	}
	// Will we use a 1-based slice for readability
	candidates := make([]bool, DIM+1)
	// Everything is valid initially
	for i := 1; i <= DIM; i++ {
		candidates[i] = true
	}
	// Check row
	for i := 0; i < DIM; i++ {
		candidates[b.cell[row][i]] = false
	}
	// Check column
	for i := 0; i < DIM; i++ {
		candidates[b.cell[i][col]] = false
	}
	// Check section
	rowStart := row / 3 * 3
	rowEnd := rowStart + DIM/3
	colStart := col / 3 * 3
	colEnd := colStart + DIM/3
	for ri := rowStart; ri < rowEnd; ri++ {
		for ci := colStart; ci < colEnd; ci++ {
			candidates[b.cell[ri][ci]] = false
		}
	}
	return candidates
}

// recursiveSolver tries to solve the board using a recursive backtracking
// algorithm
func recursiveSolver(b *Board) (solved bool) {
	if b.ValidSolution() {
		return true
	}

	row, col := b.NextEmptyCell()
	candidates := b.CellCandidates(row, col)

	// Try each candidate
	for val, avail := range candidates {
		if avail {
			b.MakeMove(row, col, val)
			solved = recursiveSolver(b)
			if solved {
				break
			}
			// Move was incorrect
			b.UnmakeMove(row, col)
		}
	}

	return solved
}

// readBoard reads a board from a text file, ignoring non-numeric characters
func readBoard(fname string) (*Board, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	b := NewBoard()
	for row := 0; row < DIM; row++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("EOF while reading row %v", row+1)
		}
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		col := 0
		for _, c := range line {
			if 47 < c && c < 58 {
				// c is numeric
				if c > 0 {
					b.MakeMove(row, col, int(c-48))
				}
				col++
			}
		}

	}

	return b, nil
}

// validateSolution cross checks each cell of the board.  Not part of the
// solver, but used to validate the solvers correctness.
func validateSolution(b Board) {
	for row := 0; row < DIM; row++ {
		for col := 0; col < DIM; col++ {
			expect := b.cell[row][col]
			b.cell[row][col] = 0
			candidates := b.CellCandidates(row, col)
			if !candidates[expect] {
				fmt.Printf("Invalid value %v at row %v, col %v\n", expect, row+1, col+1)
			}
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Puzzle filename required")
		os.Exit(1)
	}
	board, err := readBoard(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Starting configuration:")
	fmt.Println(board)

	solved := recursiveSolver(board)

	fmt.Printf("\nSolved? %v\n\n", solved)

	fmt.Println("Ending configuration:")
	fmt.Println(board)

	validateSolution(*board)
}
