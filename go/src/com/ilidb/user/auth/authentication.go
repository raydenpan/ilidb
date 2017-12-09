package auth

import (
	"com/ilidb/db"
	"com/ilidb/web"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func createRandomToken(n int) string {
	// Must seed or it is fucked up giving same values
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

//CreateILIDBLoginToken Create user ILIDB login token
func CreateILIDBLoginToken() db.LoginToken {
	// Create user ILIDB access token
	ilidbTokenString := createRandomToken(63)
	tLoginToken := db.LoginToken{Value: ilidbTokenString, Created: time.Now()}
	//fmt.Printf("Created new IlidbAccessToken:" + tLoginToken.Value + "\n")
	return tLoginToken
}

//AuthenticateUserSession authenticate user session for request
func AuthenticateUserSession(aUserIDCookieValue string, aSessionCookieValue string) (string, error) {
	return db.FetchUserSession(aUserIDCookieValue, aSessionCookieValue)
}

//HandleFacebookLogin Handle a Facebook user login
func HandleFacebookLogin(aFacebookUserLoginCode string) (db.LoginResult, error) {
	tFacebookUserToken, err := FetchFacebookUserToken(aFacebookUserLoginCode)
	if nil != err {
		return db.LoginResult{}, err
	}

	tFacebookUser, err := AuthenticateFacebookUserToken(tFacebookUserToken)
	if nil != err {
		return db.LoginResult{}, err
	}

	tLoginToken := CreateILIDBLoginToken()

	tLoginResult := CreateUserSession(tFacebookUser, tLoginToken)
	return tLoginResult, nil
}

// CreateUserSession Create a user session
func CreateUserSession(aFacebookUser web.FacebookUser, aLoginToken db.LoginToken) db.LoginResult {
	tLoginResult := db.LoginResult{Result: false, Token: "", ID: "", Name: ""}
	// Check if user already exists
	if !db.AddLoginToken(aFacebookUser.ID, aLoginToken) {
		fmt.Printf("Could not find any existing user for FacebookID:" + aFacebookUser.ID + " Name:" + aFacebookUser.Name + "\n")
		if !db.CreateUser(aFacebookUser.ID, aFacebookUser.Name, aLoginToken) {
			fmt.Printf("Could not create new user for FacebookID:" + aFacebookUser.ID + " Name:" + aFacebookUser.Name + "\n")
			return tLoginResult
		}
	}
	tLoginResult = db.LoginResult{Result: true, Token: aLoginToken.Value, ID: aFacebookUser.ID, Name: aFacebookUser.Name}
	return tLoginResult
}

//SetLoginCookies Set login cookies for successful login
func SetLoginCookies(w http.ResponseWriter, aLoginResult db.LoginResult) {
	tTimeNow := time.Now()
	tTimeTenYears := tTimeNow.AddDate(10, 0, 0)
	tMaxAge := 315360000

	var tCookieToken http.Cookie
	tCookieToken.Expires = tTimeTenYears
	tCookieToken.HttpOnly = false
	tCookieToken.MaxAge = tMaxAge
	tCookieToken.Name = "loginToken"
	tCookieToken.Path = "/"
	tCookieToken.Value = aLoginResult.Token
	http.SetCookie(w, &tCookieToken)

	var tCookieName http.Cookie
	tCookieName.Expires = tTimeTenYears
	tCookieName.HttpOnly = false
	tCookieName.MaxAge = tMaxAge
	tCookieName.Name = "name"
	tCookieName.Path = "/"
	tCookieName.Value = aLoginResult.Name
	http.SetCookie(w, &tCookieName)

	var tCookieID http.Cookie
	tCookieID.Expires = tTimeTenYears
	tCookieID.HttpOnly = false
	tCookieID.MaxAge = tMaxAge
	tCookieID.Name = "id"
	tCookieID.Path = "/"
	tCookieID.Value = aLoginResult.ID
	http.SetCookie(w, &tCookieID)
}
