package model

import "time"

// Puzzle represents a crossword puzzle. The Grid field stores the solved
// puzzle state. It is serialized as JSON in the database.
type Puzzle struct {
	ID        uint              `gorm:"primaryKey" json:"id"`
	Title     string            `json:"title"`
	Grid      string            `json:"grid"`
	Clues     map[string]string `gorm:"-" json:"clues"`
	CreatedAt time.Time         `json:"created"`
}
