package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const ilidbDB = "ilidb"

//BooksCollection Books collection
const BooksCollection = "books"

//UsersCollection Users collection
const UsersCollection = "users"

//SessionsCollection Session collection
const SessionsCollection = "sessions"

func getDatabaseCollection(aCollection string) *mgo.Collection {
	//TODO close all connections opened here
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not establish database connection: "+err.Error())
		os.Exit(1)
	}
	//defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	// Create (or fetch) ilidb DB
	var tCollection *mgo.Collection
	tCollection = session.DB(ilidbDB).C(aCollection)
	return tCollection
}

//AddLoginToken add a logintoken to an existing user
func AddLoginToken(aUserFacebookID string, aLoginToken LoginToken) bool {
	fmt.Printf("Adding new user login token for user:" + aUserFacebookID + "\n")
	tCollection := getDatabaseCollection(UsersCollection)
	change := mgo.Change{
		Update:    bson.M{"$push": bson.M{"logintokens": aLoginToken}},
		ReturnNew: true,
	}
	result := User{}
	_, err := tCollection.Find(bson.M{"facebookid": aUserFacebookID}).Apply(change, &result)
	if nil != err {
		fmt.Printf("Failed to add user login token for user:" + aUserFacebookID + "\n")
		return false
	}
	fmt.Printf("LoginToken was successfully added...\n")
	tokens, _ := json.Marshal(&result.LoginTokens)
	fmt.Printf("LoginTokens:" + string(tokens) + "\n")
	return true
}

//FetchUserSession fetch user session if it exists, otherwise return error
func FetchUserSession(aUserID string, aSessionToken string) (string, error) {
	fmt.Printf("Fetching user session data for UserID:" + aUserID + "\n")
	tCollection := getDatabaseCollection(UsersCollection)
	var tUser User
	err := tCollection.Find(bson.M{"facebookid": aUserID}).One(&tUser)
	if nil != err || nil == tUser.LoginTokens {
		return "", errors.New("could not find user session")
	}
	for i := 0; i < len(tUser.LoginTokens); i++ {
		if tUser.LoginTokens[i].Value == aSessionToken {
			return aUserID, nil
		}
	}
	return "", errors.New("")
}

//upsertBookVoteExistingUser add a book vote to an existing user
func upsertBookVoteExistingUser(aUserID string, aBookVote BookVote) bool {
	tVote, _ := json.Marshal(&aBookVote)
	fmt.Printf("Adding new book vote:" + string(tVote) + " for user:" + aUserID + "\n")
	tCollection := getDatabaseCollection(UsersCollection)
	change := mgo.Change{
		Update:    bson.M{"$push": bson.M{"bookvotes": aBookVote}},
		ReturnNew: true,
	}
	result := User{}
	_, err := tCollection.Find(bson.M{"facebookid": aUserID}).Apply(change, &result)
	if nil != err {
		panic(err)
	} else {
		fmt.Printf("Vote was successfully added...\n")
		tokens, _ := json.Marshal(&result.BookVotes)
		fmt.Printf("All user BookVotes:" + string(tokens) + "\n")
	}
	return true
}

//FetchUser fetch a user from users collection
// TODO change to use ilidb userid
func FetchUser(aUserFacebookID string) (User, error) {
	tCollection := getDatabaseCollection(UsersCollection)
	var tUser User
	err := tCollection.Find(bson.M{"facebookid": aUserFacebookID}).One(&tUser)
	return tUser, err
}

//CreateUser create a user in users collection
func CreateUser(aUserFacebookID string, aName string, aLoginToken LoginToken) bool {
	tCollection := getDatabaseCollection(UsersCollection)
	tUser := User{FacebookID: aUserFacebookID, Name: aName, Created: time.Now(), LoginTokens: []LoginToken{aLoginToken}}
	toPrint, _ := json.Marshal(&tUser)
	fmt.Printf("Creating user:\n" + string(toPrint) + "\n")
	err := tCollection.Insert(&tUser)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("User was successfully created...\n")
	}
	return true
}

//FetchBook fetch a book from books collection
func FetchBook(aBookID string) (Book, error) {
	tCollection := getDatabaseCollection(BooksCollection)
	var tBook Book
	err := tCollection.Find(bson.M{"id": aBookID}).One(&tBook)
	return tBook, err
}

//CreateBook insert a book in books collection
func CreateBook(aBook Book) bool {
	tCollection := getDatabaseCollection(BooksCollection)
	err := tCollection.Insert(&aBook)
	return nil == err
}

//UpsertBookVote User vote for a book
func UpsertBookVote(aUserID string, aBookVote BookVote) bool {
	// Check that book exists
	tBook, err := FetchBook(aBookVote.BookID)
	if nil != err || "" == tBook.ID {
		println("Could not find book with ID:" + aBookVote.BookID)
		return false
	}
	// Check that user exists
	tUser, err2 := FetchUser(aUserID)
	if nil != err2 || "" == tUser.FacebookID {
		println("Could not find user with facebookID:" + aUserID)
		return false
	}
	return upsertBookVoteExistingUser(aUserID, aBookVote)
}

//DeleteBook delete a book from books collection
func DeleteBook(aBookID string) bool {
	tCollection := getDatabaseCollection(BooksCollection)
	err := tCollection.Remove(bson.M{"id": aBookID})
	return nil == err
}

//FetchPopularBooksCategory fetch popular books of a category
func FetchPopularBooksCategory(aCategory string, aLimit int) []Book {
	tCollection := getDatabaseCollection(BooksCollection)
	tBooksIter := tCollection.Find(bson.M{"category": aCategory}).Sort("-nbrofratings").Limit(aLimit).Iter()
	tBooks := make([]Book, aLimit)
	var tBook Book
	i := 0
	for tBooksIter.Next(&tBook) {
		tBooks[i] = tBook
		i++
	}
	//If fewer results than limit
	if i < aLimit {
		tBooks = tBooks[0:i]
	}
	return tBooks
}

//FetchPopularBooks fetch popular books
func FetchPopularBooks(aLimit int) []Book {
	tCollection := getDatabaseCollection(BooksCollection)
	tBooksIter := tCollection.Find(bson.M{}).Sort("-nbrofratings").Limit(aLimit).Iter()
	tBooks := make([]Book, aLimit)
	var tBook Book
	i := 0
	for tBooksIter.Next(&tBook) {
		tBooks[i] = tBook
		i++
	}
	//If fewer results than limit
	if i < aLimit {
		tBooks = tBooks[0:i]
	}
	return tBooks
}

//deleteAllInCollection Delete all entries in a collection
func deleteAllInCollection(aCollection string) bool {
	tCollection := getDatabaseCollection(aCollection)
	_, err := tCollection.RemoveAll(nil)
	return nil == err
}

//DeleteAllBooks Delete all entries in books collection
func DeleteAllBooks() bool {
	return deleteAllInCollection(BooksCollection)
}

//FetchCollectionSize Fetch number of entries in a collection
func FetchCollectionSize(aCollection string) (int, error) {
	tCollection := getDatabaseCollection(aCollection)
	return tCollection.Count()
}

//DeleteUserByName For test purposes
func DeleteUserByName(aName string) bool {
	tCollection := getDatabaseCollection(UsersCollection)
	err := tCollection.Remove(bson.M{"name": aName})
	if err != nil {
		println("Could not remove user: " + aName)
		return false
	}
	return true
}
