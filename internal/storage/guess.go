package storage

import "github.com/isdiemer/crossword-backend/internal/model"

func GetGuessByUserAndPuzzle(userID uint, puzzleID uint) (*model.Guess, error) {
	var guess model.Guess
	err := DB.Where("user_id = ? AND puzzle_id = ?", userID, puzzleID).First(&guess).Error
	if err != nil {
		return nil, err
	}
	return &guess, nil
}
