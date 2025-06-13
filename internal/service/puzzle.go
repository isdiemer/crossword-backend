package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/isdiemer/crossword-backend/internal/storage"
	"gorm.io/gorm"
)

// ErrPuzzleNotFound indicates the requested puzzle does not exist.
var ErrPuzzleNotFound = errors.New("puzzle not found")

// ValidatePuzzle checks if the provided grid is well formed.
// It currently performs simple validation ensuring that:
//   - the grid is not empty
//   - all rows have the same length
//   - each cell contains at most one character
func ValidatePuzzle(grid [][]string) error {
	if len(grid) == 0 {
		return errors.New("grid cannot be empty")
	}

	rowLen := len(grid[0])
	if rowLen == 0 {
		return errors.New("grid rows cannot be empty")
	}

	for i, row := range grid {
		if len(row) != rowLen {
			return fmt.Errorf("row %d has incorrect length", i)
		}
		for j, cell := range row {
			if len(cell) > 1 {
				return fmt.Errorf("cell (%d,%d) has more than one character", i, j)
			}
		}
	}
	return nil
}

// CheckPuzzleSolution fetches a puzzle by ID and compares the provided guess
// grid against the stored solution. It returns true when every cell matches.
func CheckPuzzleSolution(id uint, guess [][]string) (bool, error) {
	puzzle, err := storage.GetPuzzleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, ErrPuzzleNotFound
		}
		return false, err
	}

	var solution [][]string
	if err := json.Unmarshal([]byte(puzzle.Grid), &solution); err != nil {
		return false, fmt.Errorf("invalid stored grid")
	}

	if err := ValidatePuzzle(guess); err != nil {
		return false, err
	}

	if len(solution) != len(guess) {
		return false, fmt.Errorf("grid size mismatch")
	}
	for i := range solution {
		if len(solution[i]) != len(guess[i]) {
			return false, fmt.Errorf("grid size mismatch")
		}
		for j := range solution[i] {
			if solution[i][j] != guess[i][j] {
				return false, nil
			}
		}
	}
	return true, nil
}
