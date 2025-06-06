package model

type Puzzle struct {
	ID      int               `json:"id"`
	Title   string            `json:"title"`
	Grid    [][]string        `json:"grid"`
	Clues   map[string]string `json:"Clues"`
	Created string            `json:"Created"`
}
