package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/isdiemer/crossword-backend/internal/model"
	"github.com/isdiemer/crossword-backend/internal/storage"
	"gorm.io/datatypes"
)

type CreatePuzzleRequest struct {
	Title    string         `json:"title"`
	Grid     datatypes.JSON `json:"grid"`
	Clues    datatypes.JSON `json:"clues"`
	Solution datatypes.JSON `json:"solution"`
}

func CreatePuzzleHandler(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req CreatePuzzleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	newPuzzle := model.Puzzle{
		Title:    req.Title,
		Grid:     req.Grid,
		Clues:    req.Clues,
		Solution: req.Solution,
		AuthorID: userID,
		Created:  time.Now(),
	}

	if err := storage.CreatePuzzle(&newPuzzle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save puzzle"})
		return
	}

	c.JSON(http.StatusCreated, newPuzzle)
}

func GetMyPuzzlesHandler(c *gin.Context) {
	session, _ := GetSessionFromContext(c)
	puzzles, _ := storage.GetPuzzlesByUserID(session.UserID)
	c.JSON(http.StatusOK, puzzles)
}
func GetPuzzleByIDHandler(c *gin.Context) {
	id := c.Param("id")
	var puzzle model.Puzzle
	err := storage.DB.First(&puzzle, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "puzzle not found"})
		return
	}
	c.JSON(http.StatusOK, puzzle)
}
func UpdatePuzzleHandler(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id, _ := strconv.Atoi(c.Param("id"))

	puzzle, err := storage.GetPuzzleByID(int(id))
	if err != nil || puzzle.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	var req CreatePuzzleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	puzzle.Title = req.Title
	puzzle.Grid = req.Grid
	puzzle.Clues = req.Clues
	puzzle.Solution = req.Solution

	if err := storage.UpdatePuzzle(puzzle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	c.JSON(http.StatusOK, puzzle)
}
func DeletePuzzleHandler(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id, _ := strconv.Atoi(c.Param("id"))

	puzzle, err := storage.GetPuzzleByID(int(id))
	if err != nil || puzzle.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	if err := storage.DeletePuzzleByID(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}

	c.Status(http.StatusNoContent)
}
