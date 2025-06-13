package service

import (
	"testing"

	"github.com/isdiemer/crossword-backend/internal/model"
	"github.com/isdiemer/crossword-backend/internal/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestValidatePuzzle_Valid(t *testing.T) {
	grid := [][]string{{"a", "b"}, {"c", "d"}}
	if err := ValidatePuzzle(grid); err != nil {
		t.Fatalf("expected valid grid, got error: %v", err)
	}
}

func TestValidatePuzzle_Empty(t *testing.T) {
	if err := ValidatePuzzle(nil); err == nil {
		t.Fatalf("expected error for empty grid")
	}
}

func TestValidatePuzzle_RowLength(t *testing.T) {
	grid := [][]string{{"a"}, {"b", "c"}}
	if err := ValidatePuzzle(grid); err == nil {
		t.Fatalf("expected row length error")
	}
}

func TestValidatePuzzle_CellLength(t *testing.T) {
	grid := [][]string{{"ab"}}
	if err := ValidatePuzzle(grid); err == nil {
		t.Fatalf("expected cell length error")
	}
}

func TestCheckPuzzleSolution(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	storage.DB = db
	db.AutoMigrate(&model.Puzzle{})

	puzzle := model.Puzzle{Title: "test", Grid: `[["a"]]`}
	db.Create(&puzzle)

	ok, err := CheckPuzzleSolution(puzzle.ID, [][]string{{"a"}})
	if err != nil || !ok {
		t.Fatalf("expected correct solution")
	}

	ok, err = CheckPuzzleSolution(puzzle.ID, [][]string{{"b"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("expected mismatch")
	}
}

func TestCheckPuzzleSolution_NotFound(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	storage.DB = db
	db.AutoMigrate(&model.Puzzle{})

	_, err := CheckPuzzleSolution(1, [][]string{{"a"}})
	if err != ErrPuzzleNotFound {
		t.Fatalf("expected ErrPuzzleNotFound, got %v", err)
	}
}
