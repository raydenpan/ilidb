package db

import (
	"com/ilidb/common"
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getDatabaseCollection(aCollection string) *mgo.Collection {
	//TODO close all connections opened here
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	//defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	// Create (or fetch) ilidb DB
	var tCollection *mgo.Collection
	tCollection = session.DB("swoosh").C(aCollection)
	return tCollection
}

// CreateUserSession Create a user session
func CreateUserSession(aFacebookUser common.FacebookUser, aLoginToken common.LoginToken) common.LoginResult {
	tCollection := getDatabaseCollection("users")
	// Check if user already exists
	var tUser common.User
	err := tCollection.Find(bson.M{"facebookid": aFacebookUser.ID}).One(&tUser)
	if err != nil {
		fmt.Printf("Could not find user in DB:" + aFacebookUser.ID + "\nCreating new user...\n")
		tUser = common.User{FacebookID: aFacebookUser.ID, Name: aFacebookUser.Name, Created: time.Now(), LoginTokens: []common.LoginToken{aLoginToken}}
		toPrint, _ := json.Marshal(&tUser)
		fmt.Printf("Creating user:\n" + string(toPrint) + "\n")
		err = tCollection.Insert(&tUser)
		if err != nil {
			panic(err)
		} else {
			fmt.Printf("User was successfully created...\n")
		}
	} else {
		fmt.Printf("Found user in DB:" + tUser.FacebookID + "\nAdding new user login token...\n")
		change := mgo.Change{
			Update:    bson.M{"$push": bson.M{"logintokens": aLoginToken}},
			ReturnNew: true,
		}
		result := common.User{}
		_, err = tCollection.Find(bson.M{"facebookid": aFacebookUser.ID}).Apply(change, &result)
		if err != nil {
			panic(err)
		} else {
			fmt.Printf("LoginToken was successfully added...\n")
			tokens, _ := json.Marshal(&result.LoginTokens)
			fmt.Printf("LoginTokens:" + string(tokens) + "\n")
		}
	}
	tLoginResult := common.LoginResult{Result: true, Token: aLoginToken.Value, ID: aFacebookUser.ID, Name: aFacebookUser.Name}
	return tLoginResult
}
