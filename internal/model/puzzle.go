package model

import (
	"time"

	"gorm.io/datatypes"
)

type Puzzle struct {
	ID       int            `gorm:"primaryKey" json:"id"`
	Title    string         `json:"title"`
	Grid     datatypes.JSON `json:"grid"`
	Clues    datatypes.JSON `json:"Clues"`
	Solution datatypes.JSON `json:"-"`
	Created  time.Time      `json:"Created"`
	AuthorID uint           `json:"AuthorID"`
}
type Guess struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	UserID    uint           `json:"userID"`
	PuzzleID  uint           `json:"puzzleID"`
	Grid      datatypes.JSON `json:"grid"`
}
