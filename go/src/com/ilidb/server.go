package main

import (
	"com/ilidb/auth"
	"com/ilidb/common"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func authenticateHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Authentication against Facebook initiated...\n")
	var tFacebookToken common.FacebookUserToken
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&tFacebookToken)
	if err != nil {
		panic("aaaa")
	}
	fmt.Printf("User POSTed FB access token:" + tFacebookToken.Value + "\n")

	var tLoginResult common.LoginResult
	tLoginResult = auth.HandleFacebookLogin(tFacebookToken)

	// Return user login info
	w.Header().Set("Content-Type", "application/json")
	result, _ := json.Marshal(tLoginResult)
	io.WriteString(w, string(result))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, common.LoadHTMLFileAsString("index.html"))
}

func main() {
	//Print server starting message
	fmt.Printf("Starting Ilidb.com server...\n")

	// Create file server for static html files
	fs := http.FileServer(http.Dir("../../"))

	http.Handle("/css/", fs)
	http.Handle("/img/", fs)
	http.Handle("/js/", fs)

	//Main page, index.html
	http.HandleFunc("/", indexHandler)

	// Authenticate
	http.HandleFunc("/authenticate", authenticateHandler)

	//Listen on port 9080
	http.ListenAndServe(":9080", nil)
}
