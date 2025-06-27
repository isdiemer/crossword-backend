package storage

import (
	"time"

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

func MarkPuzzleComplete(userID uint, puzzleID uint) error {
	completed := &model.CompletedPuzzle{
		UserID:      userID,
		PuzzleID:    puzzleID,
		CompletedAt: time.Now(),
	}
	return DB.Create(completed).Error
}

func GetCompletedPuzzles(userID uint) ([]uint, error) {
	var completed []model.CompletedPuzzle
	err := DB.Where("user_id = ?", userID).Find(&completed).Error
	if err != nil {
		return nil, err
	}

	puzzleIDs := make([]uint, len(completed))
	for i, c := range completed {
		puzzleIDs[i] = c.PuzzleID
	}
	return puzzleIDs, nil
}

func IsPuzzleCompleted(userID uint, puzzleID uint) (bool, error) {
	var count int64
	err := DB.Model(&model.CompletedPuzzle{}).
		Where("user_id = ? AND puzzle_id = ?", userID, puzzleID).
		Count(&count).Error
	return count > 0, err
}
