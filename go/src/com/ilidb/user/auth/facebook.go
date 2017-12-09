package auth

import (
	"com/ilidb/web"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

//FetchFacebookUserToken Fetch a Facebook user token from Facebook Graph API
func FetchFacebookUserToken(aFacebookLoginCode string) (string, error) {
	fmt.Printf("###Trying to fetch AccessToken from Facebook with user Facebook login code :\ncode:" + aFacebookLoginCode + "\n")
	// Fetch access token from FB Graph API
	resp, err := http.Get("https://graph.facebook.com/v2.10/oauth/access_token?client_id=180292159051019&client_secret=9262e3fc316aa9c2807a7b456def16c6&code=" + aFacebookLoginCode + "&redirect_uri=https://www.ilidb.com/user/authenticate/facebook")
	if err != nil {
		// TODO log error
		return "", err
	}
	web.PrintResponse(resp)
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	var tFacebookAccessTokenResponse web.FacebookAccessTokenResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tFacebookAccessTokenResponse)
	//TODO proper error handling
	if err != nil {
		fmt.Printf("###Call to Facebook /access_token failed...\n")
		panic(err)
	}
	fmt.Printf("###Successfully call to Facebook /access_token...\naccess_token:" + tFacebookAccessTokenResponse.Access_token + "\n")
	return tFacebookAccessTokenResponse.Access_token, nil
}

//AuthenticateFacebookUserToken Authenticate a Facebook user token against Facebook Graph API
func AuthenticateFacebookUserToken(aFacebookUserAccessToken string) (web.FacebookUser, error) {
	fmt.Printf("###Trying to authenticate against Facebook /me:\nFacebook user token:" + aFacebookUserAccessToken + "\n")
	// Check access token against FB using /me
	resp, err := http.Get("https://graph.facebook.com/v2.10/me?access_token=" + aFacebookUserAccessToken)
	if err != nil {
		fmt.Printf("###Call to Facebook /me failed...\n")
		//TODO proper error logging
		return web.FacebookUser{}, err
	}
	web.PrintResponse(resp)
	if resp.StatusCode != http.StatusOK {
		//TODO proper error logging
		return web.FacebookUser{}, errors.New(resp.Status)
	}
	var tFacebookUser web.FacebookUser
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tFacebookUser)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Call to Facebook /me failed for FacebookAccessToken:"+aFacebookUserAccessToken)
		//TODO proper error logging
		return web.FacebookUser{}, err
	}
	fmt.Fprintln(os.Stdout, "###Successful call to Facebook /me:\nid:" + tFacebookUser.ID + "\nname:" + tFacebookUser.Name + "\n")
	return tFacebookUser, nil
}
