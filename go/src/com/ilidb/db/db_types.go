package db

import (
	"time"
)

//User Ilidb user
type User struct {
	FacebookID  string
	Name        string
	Created     time.Time
	LoginTokens []LoginToken
	BookVotes   []BookVote
}

//LoginToken Ilidb user login token
type LoginToken struct {
	Value   string
	Created time.Time
}

//LoginResult Login result
type LoginResult struct {
	Result bool
	Token  string
	ID     string
	Name   string
}

//BookVote User vote for a book
type BookVote struct {
	BookID    string
	Rating    string
	Timestamp time.Time
}

//Book A book entry
type Book struct {
	AuthorID         string
	AuthorName       string
	ID               string
	Title            string
	OriginalLanguage string
	ReleaseYear      string
	NbrOfPages       string
	TopReview        string
	Rating           string
	NbrOfRatings     int64
	ImgURL           string
	PageURL          string
	Category         string
}
