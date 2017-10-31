package web

import (
	"fmt"
	"net/http"
	"strings"
)

//PrintRequest generates ascii representation of a request
func PrintRequest(r *http.Request) {
	//TODO this file is only used for debug
	fmt.Println("--------REQUEST------------")
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}
	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	toPrint := strings.Join(request, "\n")
	fmt.Println(toPrint)
	fmt.Println("----------------------------")
}

//PrintResponse generates ascii representation of a response
func PrintResponse(r *http.Response) {
	fmt.Println("--------RESPONSE------------")
	// Create return string
	var response []string
	// Add the request string
	status := fmt.Sprintf("%v %v", r.Status, r.Proto)
	response = append(response, status)

	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			response = append(response, fmt.Sprintf("%v: %v", name, h))
		}
	}
	// If this is a POST, add post data
	response = append(response, "\n")
	var tBody []byte
	_, err := r.Body.Read(tBody)
	if err != nil {
		panic("Coult not read response body!")
	}
	response = append(response, string(tBody))
	// Return the request as a string
	toPrint := strings.Join(response, "\n")
	fmt.Println(toPrint)
	fmt.Println("----------------------------")
}
