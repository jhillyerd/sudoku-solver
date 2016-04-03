package main

import (
	"bufio"
	"fmt"
	"os"
)

// DIM is the dimension of the board
const DIM = 9

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Puzzle filename required")
		os.Exit(1)
	}
	board, err := readGame(os.Args[1])
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

// readGame reads a board from a text file, ignoring non-numeric characters
func readGame(fname string) (*Game, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	b := NewGame()
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
			// ASCII values 48..57 represent 0..9
			if 48 <= c && c <= 57 {
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
func validateSolution(b Game) {
	for row := 0; row < DIM; row++ {
		for col := 0; col < DIM; col++ {
			// Hold on to the move for this cell
			expect := b.board[row][col]
			// Reset move and check that the expected move is in the candidate list
			b.board[row][col] = 0
			candidates := b.CellCandidates(row, col)
			if !candidates[expect] {
				fmt.Printf("Invalid value %v at row %v, col %v\n", expect, row+1, col+1)
			}
		}
	}
}
