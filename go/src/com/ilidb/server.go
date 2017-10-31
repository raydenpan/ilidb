package main

import (
	"com/ilidb/common"
	"com/ilidb/db"
	"com/ilidb/user/auth"
	"com/ilidb/web"
	"crypto/tls"
	"encoding/json"
	"strings"
	"time"

	"fmt"
	"net/http"
	"os"
)

func userAuthenticateFacebookHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	fmt.Printf("Authentication against Facebook initiated...\n")
	web.PrintRequest(req)

	// Get code parameter
	// TODO return error if missing code
	var tFacebookUserLoginCode = req.FormValue("code")
	fmt.Printf("User was redirected with FB login code:" + tFacebookUserLoginCode + "\n")

	var tLoginResult db.LoginResult
	tLoginResult = auth.HandleFacebookLogin(tFacebookUserLoginCode)

	// Set login cookies on response
	auth.SetLoginCookies(w, tLoginResult)

	// Redirect to index page (setting login cookies)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func userVoteBookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

	decoder := json.NewDecoder(r.Body)
	var tBookVote web.BookVote
	err := decoder.Decode(&tBookVote)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	fmt.Printf("Voting for book, BookID: " + tBookVote.BookID + " Rating:" + tBookVote.Rating + " UserID:" + tBookVote.UserID + "\n")
	tDBBookVote := db.BookVote{BookID: tBookVote.BookID, Rating: tBookVote.Rating, Timestamp: time.Now()}
	tSuccess := db.UpsertBookVote(tBookVote.UserID, tDBBookVote)
	if tSuccess {
		response := http.Response{StatusCode: http.StatusOK}
		response.Write(w)
	} else {
		response := http.Response{StatusCode: http.StatusInternalServerError}
		response.Write(w)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// TODO set content header on all?
	w.Header().Set("Content-Type", "text/html")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	tIndexPageString := web.GenerateIndexPage()
	fmt.Fprintf(w, tIndexPageString)
}

func booksCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

	// TODO fix me
	tBookCategory := strings.Split(r.URL.Path[1:], "/")[2]
	fmt.Printf("Fetching popular books for category: " + tBookCategory + "\n")

	tBookCategoryPageString := web.GenerateBookCategoryPage(tBookCategory)
	fmt.Fprintf(w, tBookCategoryPageString)
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

	fmt.Printf("Fetching popular books\n")

	tBooksPageString := web.GenerateBooksPage()
	fmt.Fprintf(w, tBooksPageString)
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

	// TODO fix me
	tBookID := strings.Split(r.URL.Path[1:], "/")[1]
	fmt.Printf("Fetching BookId: " + tBookID + "\n")
	// remove a in a12312321
	tBookID = tBookID[1:len(tBookID)]

	bookPage := web.GenerateBookPage(tBookID)
	fmt.Fprintf(w, bookPage)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	tPageString := common.LoadHTMLFileAsString("stats.html")
	fmt.Fprintf(w, tPageString)
}

func httpRedirectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	target := "https://" + r.Host + r.URL.Path
	fmt.Printf("Redirect to: " + target + "\n")
	http.Redirect(w, r, target, http.StatusMovedPermanently)
}

func main() {
	//Print server starting message
	fmt.Printf("Starting Ilidb.com server...\n")

	// Redirect HTTP
	go http.ListenAndServe(":8080", http.HandlerFunc(httpRedirectHandler))

	// Create file server for static html files
	fs := http.FileServer(http.Dir("../../"))
	http.Handle("/css/", fs)
	http.Handle("/img/", fs)
	http.Handle("/js/", fs)

	//Main page, index.html
	http.HandleFunc("/", indexHandler)

	// Authenticate Facebook
	// No trailing slash since FB will append URL params when redirecting
	http.HandleFunc("/user/authenticate/facebook", userAuthenticateFacebookHandler)

	// Vote
	http.HandleFunc("/user/vote/book/", userVoteBookHandler)

	// Book
	http.HandleFunc("/books/category/", booksCategoryHandler)
	http.HandleFunc("/books/", booksHandler)
	http.HandleFunc("/book/", bookHandler)

	// Stats
	http.HandleFunc("/stats/", statsHandler)

	myTLSConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			//	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			//	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			//	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			//	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			//	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			//	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		}}

	addrSecure := "0.0.0.0:8000"
	fmt.Println("Using HTTPS")
	fmt.Println("Listening on => " + addrSecure)
	myTLSWebServer := &http.Server{Addr: addrSecure, TLSConfig: myTLSConfig, Handler: nil}
	err := myTLSWebServer.ListenAndServeTLS("../../../certs/all_bundle.crt", "../../../certs/private.key")
	if err != nil {
		fmt.Println("Server cannot serve TLS on https port:8000 ", err)
		os.Exit(1)
	}
}
