package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	// check if current path exactly /. if not use http.NotFound
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox"))
}

// Add showSnippet handler
func showSnippet(w http.ResponseWriter, r *http.Request) {
	// get id on query and convert into integer
	// if less than 1 return 404
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Displays a specific snippet id %d", id)
	// w.Write([]byte("Displays a specific snippet id %d"))
}

// Add showSnippet handler
func createSnippet(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not. Note that
	// http.MethodPost is a constant equal to the string "POST".
	if r.Method != http.MethodPost {
		// If it's not, use the w.WriteHeader() method to send a 405 status
		// code and the w.Write() method to write a "Method Not Allowed"
		// response body. We then return from the function so that the
		// subsequent code is not executed.
		w.Header().Add("Allow", http.MethodPost)
		w.WriteHeader(405)
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	w.Write([]byte("Create a new specific snippet"))
}
