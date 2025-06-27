package model

import (
	"time"

	"gorm.io/datatypes"
)

type Puzzle struct {
	ID       int            `gorm:"primaryKey" json:"id"`
	Title    string         `json:"title"`
	Grid     datatypes.JSON `json:"grid"`
	Clues    datatypes.JSON `json:"clues"`
	Solution datatypes.JSON `json:"solution"`
	Created  time.Time      `json:"created"`
	AuthorID uint           `json:"authorID"`
}

type Guess struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `json:"userID"`
	PuzzleID  uint           `json:"puzzleID"`
	Grid      datatypes.JSON `json:"grid"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

type CompletedPuzzle struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `json:"userID"`
	PuzzleID    uint      `json:"puzzleID"`
	CompletedAt time.Time `json:"completedAt"`
}
