package storage

import "github.com/isdiemer/crossword-backend/internal/model"

func GetPuzzleByID(id uint) (*model.Puzzle, error) {
	var p model.Puzzle
	result := DB.First(&p, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &p, nil
}
