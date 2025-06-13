package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/isdiemer/crossword-backend/internal/service"
)

// PuzzleValidationInput defines the expected input JSON for validation.
type PuzzleValidationInput struct {
	Grid [][]string `json:"grid" binding:"required"`
}

// ValidatePuzzle handles POST /puzzles/:id/validate requests.
func ValidatePuzzle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid puzzle id"})
		return
	}

	var input PuzzleValidationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	ok, err := service.CheckPuzzleSolution(uint(id), input.Grid)
	if err != nil {
		if err == service.ErrPuzzleNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"correct": ok})
}
