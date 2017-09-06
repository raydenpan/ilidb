package auth

import (
	"com/ilidb/common"
	"com/ilidb/db"
	"fmt"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func createIlidbToken(n int) string {
	// Must seed or it is fucked up giving same values
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
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

func createLoginToken() common.LoginToken {
	// Create user ILIDB access token
	ilidbTokenString := createIlidbToken(64)
	tLoginToken := common.LoginToken{Value: ilidbTokenString, Created: time.Now()}
	fmt.Printf("Created new IlidbAccessToken:" + tLoginToken.Value + "\n")
	return tLoginToken
}

//HandleFacebookLogin Handle a Facebook user login
func HandleFacebookLogin(tFacebookUserToken common.FacebookUserToken) common.LoginResult {
	var tFacebookUser common.FacebookUser
	tFacebookUser = AuthenticateFacebookUserToken(tFacebookUserToken)

	var tLoginToken common.LoginToken
	tLoginToken = createLoginToken()

	var tLoginResult common.LoginResult
	tLoginResult = db.CreateUserSession(tFacebookUser, tLoginToken)
	return tLoginResult
}
