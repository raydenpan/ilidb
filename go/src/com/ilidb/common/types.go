package common

import (
	"time"
)

//User Ilidb user
type User struct {
	FacebookID  string
	Name        string
	Created     time.Time
	LoginTokens []LoginToken
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

//FacebookUser Facebook user
type FacebookUser struct {
	ID   string
	Name string
}

//FacebookUserToken Facebook user login token
type FacebookUserToken struct {
	Value string
}
