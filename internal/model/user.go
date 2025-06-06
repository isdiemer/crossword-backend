package model

type User struct {
	ID       int    `json:"ID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	PW       string `json:"-"`
	Created  string `json:"created"`
}
