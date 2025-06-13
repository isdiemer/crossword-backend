package storage

import (
	"github.com/isdiemer/crossword-backend/internal/model"
)

func CreatePuzzle(puzzle *model.Puzzle) error {
	return DB.Create(puzzle).Error
}

func GetPuzzlesByUserID(userID uint) ([]model.Puzzle, error) {
	var puzzles []model.Puzzle
	err := DB.Where("author_id = ?", userID).Find(&puzzles).Error
	return puzzles, err
}

func GetPuzzleByID(id int) (*model.Puzzle, error) {
	var puzzle model.Puzzle
	err := DB.First(&puzzle, id).Error
	if err != nil {
		return nil, err
	}
	return &puzzle, nil
}

func UpdatePuzzle(p *model.Puzzle) error {
	return DB.Save(p).Error
}
func DeletePuzzleByID(id uint) error {
	return DB.Delete(&model.Puzzle{}, id).Error
}
