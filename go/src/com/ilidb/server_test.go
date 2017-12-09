package main

import (
	"bytes"
	"com/ilidb/db"
	"com/ilidb/web"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestBooksPage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(booksHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		fmt.Println("Expected OK status but was:" + res.Status)
		t.Fail()
	}
	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if len(page) < 500 {
		fmt.Println("Expected books page to be more than 500 characters..." + string(page))
		t.Fail()
	}
}
func TestIndexPage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(indexHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		fmt.Println("Expected OK status but was:" + res.Status)
		t.Fail()
	}
	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if len(page) < 50 {
		fmt.Println("Expected index page to be more than 500 characters..." + string(page))
		t.Fail()
	}
}
func TestUserAuthenticationErrorPage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(userAuthenticateErrorHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		fmt.Println("Expected OK status but was:" + res.Status)
		t.Fail()
	}
	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if len(page) < 50 {
		fmt.Println("Expected user authentication error page to be more than 50 characters..." + string(page))
		t.Fail()
	}
}

func createUserAndBookInDB(t *testing.T) {
	tFaceBookID := "adas43hgffd"
	tUserName := "aName Nameaa"
	tLoginToken := db.LoginToken{Value: "sdasdsa", Created: time.Now()}
	success := db.CreateUser(tFaceBookID, tUserName, tLoginToken)
	if !success {
		t.Fail()
		println("Failed to create user")
	}
	tBook := db.Book{AuthorName: "Kalle Anka", AuthorID: "23482970", ID: "23234984", Title: "The Master and Margarita", OriginalLanguage: "Bulgaria", ReleaseYear: "1923", NbrOfPages: "478", TopReview: "This book is taking place in the eastern europe in the beginning of the 20th century", Rating: "7.7", NbrOfRatings: 323699, ImgURL: "/img/bb/23234984.jpg", PageURL: "/book/b23234984/", Category: "thriller"}
	success = db.CreateBook(tBook)
	if !success {
		t.Fail()
		println("Failed to create book")
	}
}

func deleteUserAndBookInDB(t *testing.T) {
	tBookID := "23234984"
	tUserName := "aName Nameaa"
	success := db.DeleteBook(tBookID)
	if !success {
		t.Fail()
		println("Failed to delete book")
	}
	success = db.DeleteUserByName(tUserName)
	if !success {
		t.Fail()
		println("Failed to delete user")
	}
}

type myjar struct {
	jar map[string][]*http.Cookie
}

func (p *myjar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	fmt.Printf("The URL is : %s\n", u.String())
	fmt.Printf("The cookie being set is : %s\n", cookies)
	p.jar[u.Host] = cookies
}

func (p *myjar) Cookies(u *url.URL) []*http.Cookie {
	fmt.Printf("The URL is : %s\n", u.String())
	fmt.Printf("Cookie being returned is : %s\n", p.jar[u.Host])
	return p.jar[u.Host]
}

func getClientWithAuthenticationHeaders(aURL string) *http.Client {
	client := &http.Client{}
	expiration := time.Now().Add(5 * time.Minute)
	cookie1 := &http.Cookie{Name: "id", Value: "adas43hgffd", Expires: expiration}
	cookie2 := &http.Cookie{Name: "loginToken", Value: "sdasdsa", Expires: expiration}
	jar := &myjar{}
	jar.jar = make(map[string][]*http.Cookie)
	tURL, err := url.Parse(aURL)
	if nil != err {
		log.Fatal(err)
	}
	tCookies := []*http.Cookie{cookie1, cookie2}
	jar.SetCookies(tURL, tCookies)
	client.Jar = jar
	return client
}

func TestUserSessionInvalidHeaderUserIDCookie(t *testing.T) {
	// Setup, create required test data in database
	createUserAndBookInDB(t)
	tBookID := "23234984"
	// Begin test
	ts := httptest.NewServer(http.HandlerFunc(userVoteBookHandler))
	defer ts.Close()

	// Create vote
	tBookVote := web.BookVote{BookID: tBookID, Rating: 5}
	tJSONBody, err := json.Marshal(&tBookVote)
	if err != nil {
		println("Failed to create post json body...")
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	tClient := &http.Client{}
	expiration := time.Now().Add(5 * time.Minute)
	cookie1 := &http.Cookie{Name: "id", Value: "adAs43hgffd", Expires: expiration}
	cookie2 := &http.Cookie{Name: "loginToken", Value: "sdasdsa", Expires: expiration}
	jar := &myjar{}
	jar.jar = make(map[string][]*http.Cookie)
	tURL, err := url.Parse(ts.URL)
	if nil != err {
		log.Fatal(err)
	}
	tCookies := []*http.Cookie{cookie1, cookie2}
	jar.SetCookies(tURL, tCookies)
	tClient.Jar = jar
	req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(tJSONBody))
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	res, err := tClient.Do(req)
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		fmt.Println("Expected Unauthorized status but was:" + res.Status)
		deleteUserAndBookInDB(t)
		t.Fail()
	}
	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	if len(page) != 0 {
		fmt.Println("Expected no body in response but was: " + string(page))
		deleteUserAndBookInDB(t)
		t.Fail()
	}

	// Teardown
	deleteUserAndBookInDB(t)
}

func TestUserSessionInvalidHeaderSessionCookie(t *testing.T) {
	// Setup, create required test data in database
	createUserAndBookInDB(t)
	tBookID := "23234984"
	// Begin test
	ts := httptest.NewServer(http.HandlerFunc(userVoteBookHandler))
	defer ts.Close()

	// Create vote
	tBookVote := web.BookVote{BookID: tBookID, Rating: 5}
	tJSONBody, err := json.Marshal(&tBookVote)
	if err != nil {
		println("Failed to create post json body...")
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	tClient := &http.Client{}
	expiration := time.Now().Add(5 * time.Minute)
	cookie1 := &http.Cookie{Name: "id", Value: "adas43hgffd", Expires: expiration}
	cookie2 := &http.Cookie{Name: "loginToken", Value: "sDasdsa", Expires: expiration}
	jar := &myjar{}
	jar.jar = make(map[string][]*http.Cookie)
	tURL, err := url.Parse(ts.URL)
	if nil != err {
		log.Fatal(err)
	}
	tCookies := []*http.Cookie{cookie1, cookie2}
	jar.SetCookies(tURL, tCookies)
	tClient.Jar = jar
	req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(tJSONBody))
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	res, err := tClient.Do(req)
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		fmt.Println("Expected Unauthorized status but was:" + res.Status)
		deleteUserAndBookInDB(t)
		t.Fail()
	}
	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	if len(page) != 0 {
		fmt.Println("Expected no body in response but was: " + string(page))
		deleteUserAndBookInDB(t)
		t.Fail()
	}

	// Teardown
	deleteUserAndBookInDB(t)
}

func TestUserSessionMissingHeaderUserIdCookie(t *testing.T) {
	// Setup, create required test data in database
	createUserAndBookInDB(t)
	tBookID := "23234984"
	// Begin test
	ts := httptest.NewServer(http.HandlerFunc(userVoteBookHandler))
	defer ts.Close()

	// Create vote
	tBookVote := web.BookVote{BookID: tBookID, Rating: 5}
	tJSONBody, err := json.Marshal(&tBookVote)
	if err != nil {
		println("Failed to create post json body...")
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	tClient := &http.Client{}
	expiration := time.Now().Add(5 * time.Minute)
	cookie2 := &http.Cookie{Name: "loginToken", Value: "sdasdsa", Expires: expiration}
	jar := &myjar{}
	jar.jar = make(map[string][]*http.Cookie)
	tURL, err := url.Parse(ts.URL)
	if nil != err {
		log.Fatal(err)
	}
	tCookies := []*http.Cookie{cookie2}
	jar.SetCookies(tURL, tCookies)
	tClient.Jar = jar
	req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(tJSONBody))
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	res, err := tClient.Do(req)
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		fmt.Println("Expected Unauthorized status but was:" + res.Status)
		deleteUserAndBookInDB(t)
		t.Fail()
	}
	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	if len(page) != 0 {
		fmt.Println("Expected no body in response but was: " + string(page))
		deleteUserAndBookInDB(t)
		t.Fail()
	}

	// Teardown
	deleteUserAndBookInDB(t)
}

func TestUserSessionMissingHeaderSessionCookie(t *testing.T) {
	// Setup, create required test data in database
	createUserAndBookInDB(t)
	tBookID := "23234984"
	// Begin test
	ts := httptest.NewServer(http.HandlerFunc(userVoteBookHandler))
	defer ts.Close()

	// Create vote
	tBookVote := web.BookVote{BookID: tBookID, Rating: 5}
	tJSONBody, err := json.Marshal(&tBookVote)
	if err != nil {
		println("Failed to create post json body...")
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	tClient := &http.Client{}
	expiration := time.Now().Add(5 * time.Minute)
	cookie1 := &http.Cookie{Name: "id", Value: "adas43hgffd", Expires: expiration}
	jar := &myjar{}
	jar.jar = make(map[string][]*http.Cookie)
	tURL, err := url.Parse(ts.URL)
	if nil != err {
		log.Fatal(err)
	}
	tCookies := []*http.Cookie{cookie1}
	jar.SetCookies(tURL, tCookies)
	tClient.Jar = jar
	req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(tJSONBody))
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	res, err := tClient.Do(req)
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		fmt.Println("Expected Unauthorized status but was:" + res.Status)
		deleteUserAndBookInDB(t)
		t.Fail()
	}
	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	if len(page) != 0 {
		fmt.Println("Expected no body in response but was: " + string(page))
		deleteUserAndBookInDB(t)
		t.Fail()
	}

	// Teardown
	deleteUserAndBookInDB(t)
}

func TestUserVoteBook(t *testing.T) {
	// Setup, create required test data in database
	createUserAndBookInDB(t)
	tBookID := "23234984"
	// Begin test
	ts := httptest.NewServer(http.HandlerFunc(userVoteBookHandler))
	defer ts.Close()

	// Create vote
	tBookVote := web.BookVote{BookID: tBookID, Rating: 5}
	tJSONBody, err := json.Marshal(&tBookVote)
	if err != nil {
		println("Failed to create post json body...")
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	tClient := getClientWithAuthenticationHeaders(ts.URL)
	req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(tJSONBody))
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	res, err := tClient.Do(req)
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusCreated {
		fmt.Println("Expected CREATED status but was:" + res.Status)
		deleteUserAndBookInDB(t)
		t.Fail()
	}
	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	if len(page) != 0 {
		fmt.Println("Expected no body in response but was: " + string(page))
		deleteUserAndBookInDB(t)
		t.Fail()
	}

	// Teardown
	deleteUserAndBookInDB(t)
}

func TestUserVoteBookInvalidRating(t *testing.T) {
	// Setup, create required test data in database
	createUserAndBookInDB(t)
	tBookID := "23234984"
	tBadRating := 11
	// Begin test
	ts := httptest.NewServer(http.HandlerFunc(userVoteBookHandler))
	defer ts.Close()

	// Create vote
	tBookVote := web.BookVote{BookID: tBookID, Rating: tBadRating}
	tJSONBody, err := json.Marshal(&tBookVote)
	if err != nil {
		println("Failed to create post json body...")
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	tClient := getClientWithAuthenticationHeaders(ts.URL)
	req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(tJSONBody))
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	res, err := tClient.Do(req)
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusBadRequest {
		fmt.Println("Expected BAD_REQUEST status but was:" + res.Status)
		deleteUserAndBookInDB(t)
		t.Fail()
	}
	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	if len(page) != 0 {
		fmt.Println("Expected no body in response but was: " + string(page))
		deleteUserAndBookInDB(t)
		t.Fail()
	}

	// Teardown
	deleteUserAndBookInDB(t)
}

func TestUserVoteBookInvalidBody(t *testing.T) {
	// Setup, create required test data in database
	createUserAndBookInDB(t)

	// Begin test
	ts := httptest.NewServer(http.HandlerFunc(userVoteBookHandler))
	defer ts.Close()

	// Create vote
	tBookVote := db.LoginToken{Value: "hello", Created: time.Now()}
	tJSONBody, err := json.Marshal(&tBookVote)
	if err != nil {
		println("Failed to create post json body...")
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}

	tClient := getClientWithAuthenticationHeaders(ts.URL)
	req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(tJSONBody))
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	res, err := tClient.Do(req)
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusBadRequest {
		fmt.Println("Expected BAD_REQUEST status but was:" + res.Status)
		deleteUserAndBookInDB(t)
		t.Fail()
	}
	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	if len(page) != 0 {
		fmt.Println("Expected no body in response but was: " + string(page))
		deleteUserAndBookInDB(t)
		t.Fail()
	}

	// Teardown
	deleteUserAndBookInDB(t)
}
func TestUserVoteBookInvalidBookID(t *testing.T) {
	// Setup, create required test data in database
	createUserAndBookInDB(t)
	tInvalidBookID := "9999999"
	tRating := 5
	// Begin test
	ts := httptest.NewServer(http.HandlerFunc(userVoteBookHandler))
	defer ts.Close()

	// Create vote
	tBookVote := web.BookVote{BookID: tInvalidBookID, Rating: tRating}
	tJSONBody, err := json.Marshal(&tBookVote)
	if err != nil {
		println("Failed to create post json body...")
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}

	tClient := getClientWithAuthenticationHeaders(ts.URL)
	req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(tJSONBody))
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	res, err := tClient.Do(req)
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusBadRequest {
		fmt.Println("Expected BAD_REQUEST status but was:" + res.Status)
		deleteUserAndBookInDB(t)
		t.Fail()
	}
	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		deleteUserAndBookInDB(t)
		log.Fatal(err)
	}
	if len(page) != 0 {
		fmt.Println("Expected no body in response but was: " + string(page))
		deleteUserAndBookInDB(t)
		t.Fail()
	}

	// Teardown
	deleteUserAndBookInDB(t)
}

func TestBookPage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(bookHandler))
	defer ts.Close()

	mostPopularBook := db.FetchPopularBooks(1)
	if len(mostPopularBook) != 1 {
		fmt.Println("Could not find the most popular book in DB...")
		t.Fail()
		return
	}
	mostPopularBookURL := ts.URL + "/book/b" + mostPopularBook[0].ID + "/"
	fmt.Println("Fetching page of most popular book:" + mostPopularBookURL)
	res, err := http.Get(mostPopularBookURL)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		fmt.Println("Expected OK status but was:" + res.Status)
		t.Fail()
	}
	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if len(page) < 500 {
		fmt.Println("Expected book page to be more than 500 characters..." + string(page))
		t.Fail()
	}
}

func TestBookCategoryPage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(booksCategoryHandler))
	defer ts.Close()

	mostPopularBook := db.FetchPopularBooks(1)
	if len(mostPopularBook) != 1 {
		fmt.Println("Could not find the most popular book in DB...")
		t.Fail()
		return
	}
	mostPopularBookURL := ts.URL + "/books/category/" + mostPopularBook[0].Category + "/"
	fmt.Println("Fetching page of most popular books category:" + mostPopularBookURL)
	res, err := http.Get(mostPopularBookURL)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		fmt.Println("Expected OK status but was:" + res.Status)
		t.Fail()
	}
	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if len(page) < 500 {
		fmt.Println("Expected books category page to be more than 500 characters..." + string(page))
		t.Fail()
	}
}

func TestServerHttpRedirectToHttps(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(httpRedirectHandler))
	defer ts.Close()
	fmt.Println("Test Server URL:" + ts.URL)
	tExpectedLocation := ts.URL + "/blalala/asdasads/324324/"
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("GET", tExpectedLocation, nil)
	if err != nil {
		fmt.Println("Could not create new http request")
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		if !strings.Contains(err.Error(), "http: server gave HTTP response to HTTPS client") {
			fmt.Println("Expected HTTPS server serving HTTP error but was: " + err.Error())
			t.Fail()
			return
		}
	}
	if res.StatusCode != http.StatusMovedPermanently {
		fmt.Println("Expected StatusMovedPermanently status but was:" + res.Status)
		t.Fail()
		return
	}
	tLocation, err := res.Location()
	if err != nil {
		fmt.Println("Could not get location from redirect response...")
		log.Fatal(err)
	}
	tExpectedLocation = strings.Replace(tExpectedLocation, "http://", "https://", 1)
	if tExpectedLocation != tLocation.String() {
		fmt.Println("Expected correct https redirect URL: " + tExpectedLocation + " but was: " + tLocation.String())
		t.Fail()
		return
	}
	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if !strings.Contains(string(page), "Moved Permanently") {
		fmt.Println("Expected \"Moved Permanently\" in response body but was:" + string(page))
		t.Fail()
	}
}
