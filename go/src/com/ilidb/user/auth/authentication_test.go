package auth

import (
	"com/ilidb/db"
	"com/ilidb/web"
	"testing"
	"time"
)

func TestAuthenticateUserSession(t *testing.T) {
	tName := "test Pekka"
	db.DeleteUserByName(tName)
	tFacebookUser := web.FacebookUser{ID: "5ad3s131231", Name: tName}
	tLoginToken := db.LoginToken{Value: "35345dfgfdgfdgdfg", Created: time.Now()}
	result := CreateUserSession(tFacebookUser, tLoginToken)
	if result.Result != true {
		t.Fail()
		println("Could not create user session for user: " + tName)
	}
	tUserID, err := AuthenticateUserSession(tFacebookUser.ID, tLoginToken.Value)
	if nil != err {
		println(err.Error())
		t.Fail()
	}
	if "" == tUserID {
		println("Did not find session for user...")
		t.Fail()
	}
	if !db.DeleteUserByName(tName) {
		t.Fail()
	}
}

func TestCreateUserSessionNewUser(t *testing.T) {
	tName := "test Pekka"
	db.DeleteUserByName(tName)
	tFacebookUser := web.FacebookUser{ID: "5ad3s131231", Name: tName}
	tLoginToken := db.LoginToken{Value: "35345dfgfdgfdgdfg", Created: time.Now()}
	result := CreateUserSession(tFacebookUser, tLoginToken)
	if result.Result != true {
		t.Fail()
		println("Could not create user session for user: " + tName)
	}
	if !db.DeleteUserByName(tName) {
		t.Fail()
	}
}

func TestCreateUserSessionExistingUser(t *testing.T) {
	tName := "test Pekka"
	db.DeleteUserByName(tName)
	tFacebookUser := web.FacebookUser{ID: "5ad3s131231", Name: tName}
	tLoginToken := db.LoginToken{Value: "35345dfgfdgfdgdfg", Created: time.Now()}
	result := CreateUserSession(tFacebookUser, tLoginToken)
	if result.Result != true {
		t.Fail()
		println("Could not create user session for new user: " + tName)
	}
	tLoginToken = db.LoginToken{Value: "asa678678trty", Created: time.Now()}
	result = CreateUserSession(tFacebookUser, tLoginToken)
	if result.Result != true {
		t.Fail()
		println("Could not create user session for existing user: " + tName)
	}
	if !db.DeleteUserByName(tName) {
		t.Fail()
	}
}

func TestCreateRandomToken(t *testing.T) {
	stringSet := getStringSet()
	for count := 0; count < 100; count++ {
		if !stringSet.Add(CreateILIDBLoginToken().Value) {
			t.Fail()
			println("Failed to create unique ILIDBLoginToken")
		}
		time.Sleep(1)
	}
}

func TestStringSet(t *testing.T) {
	stringSet := getStringSet()
	if !stringSet.Add("asd") {
		t.Fail()
	}
	if stringSet.Add("asd") {
		t.Fail()
	}
}

func getStringSet() stringSet {
	set := make(map[string]bool)
	stringSet := stringSet{set: set}
	return stringSet
}

type stringSet struct {
	set map[string]bool
}

func (set *stringSet) Add(s string) bool {
	_, found := set.set[s]
	set.set[s] = true
	return !found //False if it existed already
}
