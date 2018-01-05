package db

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TestDeleteAllBooks(t *testing.T) {
	if !DeleteAllBooks() {
		t.Fail()
		println("Could not delete all entries in books collection")
	}
}

func TestFetchUsersDB(t *testing.T) {
	tCollection := getDatabaseCollection(UsersCollection)
	if nil == tCollection {
		t.Fail()
	}
}

func TestFetchUserJohnLindskog(t *testing.T) {
	tCollection := getDatabaseCollection(UsersCollection)
	var result User
	err := tCollection.Find(bson.M{"name": "John Lindskog"}).One(&result)
	if err != nil || result.Name == "" {
		t.Fail()
		println("Could not find user with name: John Lindskog in collection users")
	}
	if result.FacebookID == "" {
		t.Fail()
		println("User John Lindskog does not have a FacebookID")
	}
}

func TestFetchInvalidUser(t *testing.T) {
	tFacebookUserID := "asdsadsa"
	_, err := FetchUser(tFacebookUserID)
	if nil == err {
		t.Fail()
		println("Should not be able to fetch non-existing user")
	}
}

func TestAddLoginToken(t *testing.T) {
	tFaceBookID := "adas43hgffd"
	tUserName := "aName Nameaa"
	tLoginToken := LoginToken{Value: "sdasdsa", Created: time.Now()}
	success := CreateUser(tFaceBookID, tUserName, tLoginToken)
	if !success {
		t.Fail()
		println("Failed to create user")
	}
	success = AddLoginToken(tFaceBookID, tLoginToken)
	if !success {
		t.Fail()
		println("Failed to add LoginToken")
	}
	if !DeleteUserByName(tUserName) {
		t.Fail()
		println("Failed to delete test user")
	}
}

func TestAddLoginTokenFail(t *testing.T) {
	tFaceBookID := "adas43hgffd"
	tLoginToken := LoginToken{Value: "sdasdsa", Created: time.Now()}
	success := AddLoginToken(tFaceBookID, tLoginToken)
	if success {
		t.Fail()
		println("Should not be able to add LoginToken")
	}
}

func TestFetchUserSession(t *testing.T) {
	tFaceBookID := "adas43hgffd"
	tUserName := "aName Nameaa"
	tLoginToken := LoginToken{Value: "sdasdsa", Created: time.Now()}
	success := CreateUser(tFaceBookID, tUserName, tLoginToken)
	if !success {
		t.Fail()
		println("Failed to create user")
	}
	success = AddLoginToken(tFaceBookID, tLoginToken)
	if !success {
		t.Fail()
		println("Failed to add LoginToken")
	}
	tUser, err := FetchUserSession(tFaceBookID, tLoginToken.Value)
	if nil != err {
		t.Fail()
		println("Failed to fetch user session for UserID: " + tUser.FacebookID)
	}
	if !DeleteUserByName(tUserName) {
		t.Fail()
		println("Failed to delete test user")
	}
}

func TestCreateFetchRemoveBook(t *testing.T) {
	tBook := Book{AuthorName: "Kalle Anka", AuthorID: "23482970", ID: "23234984", Title: "The Master and Margarita", OriginalLanguage: "Bulgaria", ReleaseYear: "1923", NbrOfPages: "478", TopReview: "This book is taking place in the eastern europe in the beginning of the 20th century", Rating: "7.7", NbrOfRatings: 323699, ImgURL: "/img/bb/23234984.jpg", PageURL: "/book/b23234984/", Category: "thriller"}
	success := CreateBook(tBook)
	if !success {
		t.Fail()
		println("Failed to create book")
	}
	tBook2, err := FetchBook(tBook.ID)
	if nil != err || tBook != tBook2 {
		t.Fail()
		println("Could not fetch book ID:" + tBook.ID)
	}
	success = DeleteBook(tBook.ID)
	if !success {
		t.Fail()
		println("Failed to delete book")
	}
	_, err = FetchBook(tBook.ID)
	if nil == err {
		t.Fail()
		println("Book was not deleted, bookID:" + tBook.ID)
	}
}

func TestUpsertBookVote(t *testing.T) {
	tFaceBookID := "adas43hgffd"
	tUserName := "aName Nameaa"
	tLoginToken := LoginToken{Value: "sdasdsa", Created: time.Now()}
	success := CreateUser(tFaceBookID, tUserName, tLoginToken)
	if !success {
		t.Fail()
		println("Failed to create user")
	}
	tBook := Book{AuthorName: "Kalle Anka", AuthorID: "23482970", ID: "23234984", Title: "The Master and Margarita", OriginalLanguage: "Bulgaria", ReleaseYear: "1923", NbrOfPages: "478", TopReview: "This book is taking place in the eastern europe in the beginning of the 20th century", Rating: "7.7", NbrOfRatings: 323699, ImgURL: "/img/bb/23234984.jpg", PageURL: "/book/b23234984/", Category: "thriller"}
	success = CreateBook(tBook)
	if !success {
		t.Fail()
		println("Failed to create book")
	}
	tBookVote := BookVote{BookID: tBook.ID, Rating: 5, Timestamp: time.Now()}
	success = UpsertBookVote(tFaceBookID, tBookVote)
	if !success {
		t.Fail()
		println("Failed to vote for book")
	}
	success = DeleteBook(tBook.ID)
	if !success {
		t.Fail()
		println("Failed to delete book")
	}
	success = DeleteUserByName(tUserName)
	if !success {
		t.Fail()
		println("Failed to delete user")
	}
}

func TestUpsertBookVoteInvalidUser(t *testing.T) {
	tFaceBookID := "adas43h6ffd"

	tBook := Book{AuthorName: "Kalle Anka", AuthorID: "23482970", ID: "23234984", Title: "The Master and Margarita", OriginalLanguage: "Bulgaria", ReleaseYear: "1923", NbrOfPages: "478", TopReview: "This book is taking place in the eastern europe in the beginning of the 20th century", Rating: "7.7", NbrOfRatings: 323699, ImgURL: "/img/bb/23234984.jpg", PageURL: "/book/b23234984/", Category: "thriller"}

	success := CreateBook(tBook)
	if !success {
		t.Fail()
		println("Failed to create book")
	}
	tBookVote := BookVote{BookID: tBook.ID, Rating: 5, Timestamp: time.Now()}
	success = UpsertBookVote(tFaceBookID, tBookVote)
	if success {
		t.Fail()
		println("Should not be able to vote for book without existing user")
	}
	success = DeleteBook(tBook.ID)
	if !success {
		t.Fail()
		println("Failed to delete book")
	}
}

func TestUpsertBookVoteInvalidBook(t *testing.T) {
	tFaceBookID := "adas43hgffd"
	tUserName := "aName Nameaa"
	tLoginToken := LoginToken{Value: "sdasdsa", Created: time.Now()}
	success := CreateUser(tFaceBookID, tUserName, tLoginToken)
	if !success {
		t.Fail()
		println("Failed to create user")
	}
	tBook := Book{AuthorName: "Kalle Anka", AuthorID: "23482970", ID: "23234984", Title: "The Master and Margarita", OriginalLanguage: "Bulgaria", ReleaseYear: "1923", NbrOfPages: "478", TopReview: "This book is taking place in the eastern europe in the beginning of the 20th century", Rating: "7.7", NbrOfRatings: 323699, ImgURL: "/img/bb/23234984.jpg", PageURL: "/book/b23234984/", Category: "thriller"}

	tBookVote := BookVote{BookID: tBook.ID, Rating: 4, Timestamp: time.Now()}
	success = UpsertBookVote(tFaceBookID, tBookVote)
	if success {
		t.Fail()
		println("Should not be able to vote for non-existing book")
	}
	success = DeleteUserByName(tUserName)
	if !success {
		t.Fail()
		println("Failed to delete user")
	}
}

func TestFetchPopularBooks(t *testing.T) {
	tSize, err := FetchCollectionSize(BooksCollection)
	if nil != err {
		t.Fail()
		println("Could not fetch size of books collection")
		return
	}

	if tSize != 0 {
		t.Fail()
		println("Collection books is not empty...")
		return
	}

	tBook1 := Book{AuthorID: "23482971", ID: "23234984", Title: "The Master and Margarita1", OriginalLanguage: "Bulgaria", ReleaseYear: "1923", NbrOfPages: "478", TopReview: "This book is taking place in the eastern europe in the beginning of the 20th century", Rating: "7.7", NbrOfRatings: 14, ImgURL: "/img/bb/23234984.jpg", PageURL: "/book/b23234984/", Category: "thriller"}
	tBook2 := Book{AuthorID: "23482971", ID: "23234985", Title: "The Master and Margarita2", OriginalLanguage: "Bulgaria", ReleaseYear: "1923", NbrOfPages: "478", TopReview: "This book is taking place in the eastern europe in the beginning of the 20th century", Rating: "7.7", NbrOfRatings: 788899, ImgURL: "/img/bb/23234984.jpg", PageURL: "/book/b23234984/", Category: "romance"}
	tBook3 := Book{AuthorID: "23482971", ID: "23234986", Title: "The Master and Margarita3", OriginalLanguage: "Bulgaria", ReleaseYear: "1923", NbrOfPages: "478", TopReview: "This book is taking place in the eastern europe in the beginning of the 20th century", Rating: "7.7", NbrOfRatings: 9123431, ImgURL: "/img/bb/23234984.jpg", PageURL: "/book/b23234984/", Category: "thriller"}
	tBook4 := Book{AuthorID: "23482971", ID: "23234987", Title: "The Master and Margarita4", OriginalLanguage: "Bulgaria", ReleaseYear: "1923", NbrOfPages: "478", TopReview: "This book is taking place in the eastern europe in the beginning of the 20th century", Rating: "7.7", NbrOfRatings: 12333, ImgURL: "/img/bb/23234984.jpg", PageURL: "/book/b23234984/", Category: "thriller"}

	success1 := CreateBook(tBook1)
	success2 := CreateBook(tBook2)
	success3 := CreateBook(tBook3)
	success4 := CreateBook(tBook4)
	if !success1 || !success2 || !success3 || !success4 {
		t.Fail()
		println("Failed to create books")
		return
	}

	tCategory := "thriller"
	tBooksArray := FetchPopularBooks(20)
	if len(tBooksArray) != 4 {
		t.Fail()
		println("Failed to fetch popular books, expected 4 books but found:%d", len(tBooksArray))
		return
	}
	tBooksArray = FetchPopularBooksCategory(tCategory, 20)
	if len(tBooksArray) != 3 {
		t.Fail()
		println("Failed to fetch popular books for category:"+tCategory+" expected 3 books but found:%d", len(tBooksArray))
		return
	}
	tFoundBook := tBooksArray[0]
	if tFoundBook.ID != tBook3.ID {
		t.Fail()
		println("Expected book with most ratings: " + printBook(tBook3) + " but was: " + printBook(tFoundBook))
		return
	}
	tFoundBook = tBooksArray[1]
	if tFoundBook.ID != tBook4.ID {
		t.Fail()
		println("Expected book with second most ratings: " + printBook(tBook4) + " but was: " + printBook(tFoundBook))
		return
	}
	tFoundBook = tBooksArray[2]
	if tFoundBook.ID != tBook1.ID {
		t.Fail()
		println("Expected book with third most ratings: " + printBook(tBook1) + " but was: " + printBook(tFoundBook))
		return
	}
	success1 = DeleteBook(tBook1.ID)
	success2 = DeleteBook(tBook2.ID)
	success3 = DeleteBook(tBook3.ID)
	success4 = DeleteBook(tBook4.ID)
	if !success1 || !success2 || !success3 || !success4 {
		t.Fail()
		println("Failed to delete books")
	}
}

func printBook(book Book) string {
	return "ID:" + book.ID + ", NbrOfRatings:" + fmt.Sprint(book.NbrOfRatings) + ", Category:" + book.Category
}

func TestCreateBooks(t *testing.T) {
	tBook1 := Book{AuthorName: "Mikhail Bulgakov", AuthorID: "23482971", ID: "23234967", Title: "The Master and Margarita1", OriginalLanguage: "Bulgaria", ReleaseYear: "1923", NbrOfPages: "478", TopReview: "This book is taking place in the eastern europe in the beginning of the 20th century", Rating: "7.7", NbrOfRatings: 14, ImgURL: "/img/bb/23234984.jpg", PageURL: "/book/b23234967/", Category: "thriller"}
	tBook2 := Book{AuthorName: "Mikhail Bulgakov", AuthorID: "23482971", ID: "23234968", Title: "The Master and Margarita2", OriginalLanguage: "Bulgaria", ReleaseYear: "1923", NbrOfPages: "478", TopReview: "This book is taking place in the eastern europe in the beginning of the 20th century", Rating: "7.7", NbrOfRatings: 788899, ImgURL: "/img/bb/23234984.jpg", PageURL: "/book/b23234968/", Category: "thriller"}
	tBook3 := Book{AuthorName: "Mikhail Bulgakov", AuthorID: "23482971", ID: "23234969", Title: "The Master and Margarita3", OriginalLanguage: "Bulgaria", ReleaseYear: "1923", NbrOfPages: "478", TopReview: "This book is taking place in the eastern europe in the beginning of the 20th century", Rating: "7.7", NbrOfRatings: 9123431, ImgURL: "/img/bb/23234984.jpg", PageURL: "/book/b23234969/", Category: "scifi"}
	tBook4 := Book{AuthorName: "Mikhail Bulgakov", AuthorID: "23482971", ID: "23234970", Title: "The Master and Margarita4", OriginalLanguage: "Bulgaria", ReleaseYear: "1923", NbrOfPages: "478", TopReview: "This book is taking place in the eastern europe in the beginning of the 20th century", Rating: "7.7", NbrOfRatings: 12333, ImgURL: "/img/bb/23234984.jpg", PageURL: "/book/b23234970/", Category: "fantasy"}

	success1 := CreateBook(tBook1)
	success2 := CreateBook(tBook2)
	success3 := CreateBook(tBook3)
	success4 := CreateBook(tBook4)
	if !success1 || !success2 || !success3 || !success4 {
		t.Fail()
		println("Failed to create books")
		return
	}
}

func TestSearchBookTitle(t *testing.T) {
	tTitle := "The"
	tResult := SearchBookTitle(tTitle)
	if len(tResult) == 0 {
		t.Fail()
		println("Found zero books matching:" + tTitle)
		return
	}
	t.Fail()
	println("Found " + strconv.Itoa(len(tResult)) + "books matching:" + tTitle)
}

func TestCreateIndex(t *testing.T) {
	tCollection := getDatabaseCollection(BooksCollection)
	index := mgo.Index{
		Key: []string{"$text:title"},
	}

	err := tCollection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func TestDropIndex(t *testing.T) {
	tCollection := getDatabaseCollection(BooksCollection)
	err := tCollection.DropIndexName("title_text")
	if err != nil {
		panic(err)
	}
}
