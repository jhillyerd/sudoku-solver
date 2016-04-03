package main

import (
	"fmt"
)

// Game represents a sudoku board
type Game struct {
	// board represents the game board, access as board[row][col]
	board      [][]int
	remaining  int
	backtracks int
}

// NewGame creates an empty sudoku board
func NewGame() *Game {
	g := &Game{}
	// Build empty (zero) board matrix
	g.board = make([][]int, DIM)
	for i := range g.board {
		g.board[i] = make([]int, DIM)
	}
	g.remaining = DIM * DIM
	return g
}

// String formats the board for human consumption
func (g *Game) String() string {
	var result = "    1 2 3 4 5 6 7 8 9\n"
	for i, row := range g.board {
		result += fmt.Sprintf("%v: %v\n", i+1, row)
	}
	result += fmt.Sprintf("Remaining: %v, Backtracks: %v", g.remaining, g.backtracks)
	return result
}

// ValidSolution is true if remaining == 0
func (g *Game) ValidSolution() bool {
	return g.remaining == 0
}

// MakeMove adds a number to the board, row and col indices are 0 based
func (g *Game) MakeMove(row, col, val int) {
	if g.board[row][col] == 0 && val != 0 {
		g.remaining--
	}
	g.board[row][col] = val
}

// UnmakeMove removes a number from the board, row and col indices are 0 based
func (g *Game) UnmakeMove(row, col int) {
	if g.board[row][col] != 0 {
		g.remaining++
		g.board[row][col] = 0
	}
	g.backtracks++
}

// NextEmptyCell tells our solver which cell to work on next
func (g *Game) NextEmptyCell() (row, col int) {
	min := DIM + 1
	for ri, cols := range g.board {
		for ci, val := range cols {
			if val == 0 {
				cur := 0
				candidates := g.CellCandidates(ri, ci)
				// Count candidates
				for i := 1; i <= DIM; i++ {
					if candidates[i] {
						cur++
					}
				}
				if cur < min {
					row, col = ri, ci
					min = cur
				}
			}
		}
	}
	// Return row, col
	return
}

// CellCandidates returns a list of legal moves for specified cell
func (g *Game) CellCandidates(row, col int) []bool {
	if row < 0 || DIM < row {
		panic(fmt.Sprintf("Invalid row passed: %v", row))
	}
	if col < 0 || DIM < col {
		panic(fmt.Sprintf("Invalid col passed: %v", col))
	}
	// Will we use a 1-based slice for readability, 0 will always be false
	candidates := make([]bool, DIM+1)
	// Set everything to valid (except 0)
	for i := 1; i <= DIM; i++ {
		candidates[i] = true
	}
	// Check row
	for i := 0; i < DIM; i++ {
		candidates[g.board[row][i]] = false
	}
	// Check column
	for i := 0; i < DIM; i++ {
		candidates[g.board[i][col]] = false
	}
	// Check section
	rowStart := row / 3 * 3
	rowEnd := rowStart + DIM/3
	colStart := col / 3 * 3
	colEnd := colStart + DIM/3
	for ri := rowStart; ri < rowEnd; ri++ {
		for ci := colStart; ci < colEnd; ci++ {
			candidates[g.board[ri][ci]] = false
		}
	}
	return candidates
}

// recursiveSolver tries to solve the board using a recursive backtracking
// algorithm
func recursiveSolver(g *Game) (solved bool) {
	if g.ValidSolution() {
		return true
	}

	row, col := g.NextEmptyCell()
	candidates := g.CellCandidates(row, col)

	// Try each candidate
	for val, avail := range candidates {
		if avail {
			g.MakeMove(row, col, val)
			solved = recursiveSolver(g)
			if solved {
				break
			}
			// Move was incorrect
			g.UnmakeMove(row, col)
		}
	}

	return solved
}
