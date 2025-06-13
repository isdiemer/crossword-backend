package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/isdiemer/crossword-backend/internal/model"
	"github.com/isdiemer/crossword-backend/internal/storage"
	"gorm.io/datatypes"
)

type createPuzzleRequest struct {
	Title    string          `json:"title"`
	Grid     json.RawMessage `json:"grid"`
	Clues    json.RawMessage `json:"clues"`
	Solution json.RawMessage `json:"solution"`
}

func CreatePuzzleHandler(c *gin.Context) {

	session, err := GetSessionFromContext(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req createPuzzleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	puzzle := model.Puzzle{
		Title:    req.Title,
		Grid:     datatypes.JSON(req.Grid),
		Clues:    datatypes.JSON(req.Clues),
		Solution: datatypes.JSON(req.Solution),
		Created:  time.Now(),
		AuthorID: session.UserID,
	}

	if err := storage.CreatePuzzle(&puzzle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save puzzle"})
		return
	}

	c.JSON(http.StatusOK, puzzle)
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
