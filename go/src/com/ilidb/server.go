package main

import (
	"fmt"
	"net/http"
)

func main() {
	//Print server starting message
	fmt.Printf("Starting Ilidb.com server...\n")

	// Create file server for static html files
	fs := http.FileServer(http.Dir("../../../../html"))

	//Main page, index.html
	http.Handle("/", fs)

	//Listen on port 9080
	http.ListenAndServe(":9080", nil)
}
