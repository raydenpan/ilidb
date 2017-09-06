package auth

import (
	"com/ilidb/common"
	"encoding/json"
	"fmt"
	"net/http"
)

//AuthenticateFacebookUserToken Authenticate a Facebook user token against Facebook Graph API
func AuthenticateFacebookUserToken(aFacebookUserToken common.FacebookUserToken) common.FacebookUser {
	fmt.Printf("Trying to authenticate against Facebook /me:\nFacebook user token:" + aFacebookUserToken.Value + "\n")
	// Check access token against FB using /me
	resp, err := http.Get("https://graph.facebook.com/v2.5/me?access_token=" + aFacebookUserToken.Value)
	if err != nil {
		panic(err)
	}
	var tFacebookUser common.FacebookUser
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tFacebookUser)
	//TODO proper error handling
	if err != nil {
		panic("Call to Facebook /me failed...")
	}
	fmt.Printf("Successful call to Facebook /me:\nid:" + tFacebookUser.ID + "\nname:" + tFacebookUser.Name + "\n")
	return tFacebookUser
}
