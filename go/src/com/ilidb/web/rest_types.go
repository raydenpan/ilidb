package web

//BookVote User vote for a book
type BookVote struct {
	BookID string
	Rating int
}

//FacebookAccessTokenResponse Facebook access token response
type FacebookAccessTokenResponse struct {
	Access_token string
	Token_type   string
	Expires_in   int64
}

//FacebookUser Facebook user
type FacebookUser struct {
	ID   string
	Name string
}
