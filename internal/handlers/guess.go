package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/isdiemer/crossword-backend/internal/model"
	"github.com/isdiemer/crossword-backend/internal/storage"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ValidateRequest struct {
	Grid datatypes.JSON `json:"grid"`
}

type ValidateResponse struct {
	Correct bool `json:"correct"`
}

func SubmitGuessHandler(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var input struct {
		PuzzleID int            `json:"puzzleId"`
		Grid     datatypes.JSON `json:"grid"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	var guess *model.Guess

	guess, err := storage.GetGuessByUserAndPuzzle(userID, uint(input.PuzzleID))

	if err == nil {
		guess.Grid = input.Grid
		if err := storage.DB.Save(&guess).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update guess"})
			return
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		guess = &model.Guess{
			UserID:   userID,
			PuzzleID: uint(input.PuzzleID),
			Grid:     input.Grid,
		}
		if err := storage.DB.Create(&guess).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save guess"})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error checking existing guess"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "guess saved"})
}

func ValidateGuessHandler(c *gin.Context) {
	idStr := c.Param("id")
	puzzleID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid puzzle id"})
		return
	}

	// Parse request body
	var req ValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	// Fetch puzzle from DB
	puzzle, err := storage.GetPuzzleByID(int(puzzleID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "puzzle not found"})
		return
	}

	// Compare grids
	var guessGrid, solutionGrid [][]string
	json.Unmarshal(req.Grid, &guessGrid)
	json.Unmarshal(puzzle.Solution, &solutionGrid)

	correct := reflect.DeepEqual(guessGrid, solutionGrid)

	// Add logging
	logMsg := struct {
		PuzzleID     int        `json:"puzzle_id"`
		GuessGrid    [][]string `json:"guess_grid"`
		SolutionGrid [][]string `json:"solution_grid"`
		Correct      bool       `json:"correct"`
	}{
		PuzzleID:     puzzleID,
		GuessGrid:    guessGrid,
		SolutionGrid: solutionGrid,
		Correct:      correct,
	}
	logData, _ := json.MarshalIndent(logMsg, "", "  ")
	println("[ValidateGuess]", string(logData))

	c.JSON(http.StatusOK, ValidateResponse{Correct: correct})
}

func GetGuessHandler(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	puzzleIDStr := c.Param("puzzleId")
	puzzleID, err := strconv.ParseUint(puzzleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid puzzle ID"})
		return
	}

	guess, err := storage.GetGuessByUserAndPuzzle(userID, uint(puzzleID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no saved progress found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching guess"})
		return
	}

	c.JSON(http.StatusOK, guess)
}
