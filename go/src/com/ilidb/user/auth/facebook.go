package auth

import (
	"com/ilidb/web"
	"encoding/json"
	"fmt"
	"net/http"
)

//FetchFacebookUserToken Fetch a Facebook user token from Facebook Graph API
func FetchFacebookUserToken(aFacebookLoginCode string) string {
	fmt.Printf("###Trying to fetch AccessToken from Facebook with user Facebook login code :\ncode:" + aFacebookLoginCode + "\n")
	// Fetch access token from FB Graph API
	resp, err := http.Get("https://graph.facebook.com/v2.10/oauth/access_token?client_id=180292159051019&client_secret=cf374c2c25f6f06b8a1d64fa78517861&code=" + aFacebookLoginCode + "&redirect_uri=https://www.ilidb.com/authenticate/facebook")
	if err != nil {
		panic(err)
	}
	web.PrintResponse(resp)
	var tFacebookAccessTokenResponse web.FacebookAccessTokenResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tFacebookAccessTokenResponse)
	//TODO proper error handling
	if err != nil {
		fmt.Printf("###Call to Facebook /access_token failed...\n")
		panic(err)
	}
	fmt.Printf("###Successfully call to Facebook /access_token:\naccess_token:" + tFacebookAccessTokenResponse.Access_token + "\n")
	return tFacebookAccessTokenResponse.Access_token
}

//AuthenticateFacebookUserToken Authenticate a Facebook user token against Facebook Graph API
func AuthenticateFacebookUserToken(aFacebookUserAccessToken string) web.FacebookUser {
	fmt.Printf("###Trying to authenticate against Facebook /me:\nFacebook user token:" + aFacebookUserAccessToken + "\n")
	// Check access token against FB using /me
	resp, err := http.Get("https://graph.facebook.com/v2.10/me?access_token=" + aFacebookUserAccessToken)
	if err != nil {
		fmt.Printf("###Call to Facebook /me failed...\n")
		panic(err)
	}
	web.PrintResponse(resp)

	var tFacebookUser web.FacebookUser
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tFacebookUser)
	//TODO proper error handling
	if err != nil {
		panic("###Call to Facebook /me failed...")
	}
	fmt.Printf("###Successful call to Facebook /me:\nid:" + tFacebookUser.ID + "\nname:" + tFacebookUser.Name + "\n")
	return tFacebookUser
}
