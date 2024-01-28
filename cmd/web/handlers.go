package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/andsus/snippetbox/pkg/models"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Removed because Pat matches "/" path exactly
	// r.URL.Path != "/"
	// check if current path exactly /. if not use http.NotFound
	// if r.URL.Path != "/" {
	// 	app.notFound(w) // use the notFound() helper
	// 	return
	// }

	// deliberate panic
	// panic("oops, something went bad")

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Create an instance of a templateData struct holding the slice of snippets
	// data := &templateData{Snippets: s}

	// for _, snippet := range s {
	// 	fmt.Fprintf(w, "%v\n", snippet)
	// }

	// // Initialize a slice containing the paths to the two files. Note that the
	// // home.page.tmpl file must be the *first* file in the slice.
	// files := []string{
	// 	"./ui/html/home.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partial.tmpl",
	// }

	// // Use template.ParseFiles() to read template
	// // if error send to 500 internal server error response
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// // Use Execute method on template and write into response
	// // the last parameter represent dynamic data, nil for now
	// err = ts.Execute(w, data)
	// if err != nil {
	// 	app.serverError(w, err)
	// }
	//w.Write([]byte("Hello from Snippetbox"))
	// Use the new render helper.
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

// Add showSnippet handler
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// get id on query and convert into integer
	// if less than 1 return 404
	// Pat doesn't strip the colon from the named capture key, so we need to
	// get the value of ":id" from the query string instead of "id".
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Use the new render helper.
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
	// Create an instance of a templateData struct holding the snippet data.
	// data := &templateData{Snippet: s}

	// Initialize a slice containing the paths to the show.page.tmpl file,
	// plus the base layout and footer partial that we made earlier.
	// files := []string{
	// 	"./ui/html/show.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partial.tmpl",
	// }
	// Parse the template files...
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }
	// And then execute them. Notice how we are passing in the snippet
	// data (a models.Snippet struct) as the final parameter.
	// err = ts.Execute(w, data)
	// if err != nil {
	// 	app.serverError(w, err)
	// }
	// Write the snippet data as a plain-text HTTP response body.
	// fmt.Fprintf(w, "%v", s)
	// fmt.Fprintf(w, "Displays a specific snippet id %d", id)
	// w.Write([]byte("Displays a specific snippet id %d"))
}

// Add a new createSnippetForm handler, which for now returns a placeholder response.
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
	// w.Write([]byte("Create a new snippet..."))
}

// Add showSnippet handler
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	// Checking POST request not required

	// Use r.Method to check whether the request is using POST or not. Note that
	// http.MethodPost is a constant equal to the string "POST".
	// if r.Method != http.MethodPost {
	// 	// If it's not, use the w.WriteHeader() method to send a 405 status
	// 	// code and the w.Write() method to write a "Method Not Allowed"
	// 	// response body. We then return from the function so that the
	// 	// subsequent code is not executed.
	// 	w.Header().Add("Allow", http.MethodPost)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	// title := "O snail"
	// content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	// expires := "7"

	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT and PATCH
	// requests. If there are any errors, we use our app.ClientError helper to send
	// a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Use the r.PostForm.Get() method to retrieve the relevant data fields
	// from the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	// Initialize a map to hold any validation errors.
	errors := make(map[string]string)

	// Check that the title field is not blank and is not more than 100 characters
	// long. If it fails either of those checks, add a message to the errors
	// map using the field name as the key.
	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long (maximum is 100 characters)"
	}
	// Check that the Content field isn't blank.
	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}
	// Check the expires field isn't blank and matches one of the permitted
	// values ("1", "7" or "365").
	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid"
	}
	// If there are any validation errors, re-display the create.page.tmpl
	// template passing in the validation errors and previously submitted
	// r.PostForm data.
	if len(errors) > 0 {
		app.render(w, r, "create.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})
		return
	}

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Redirect the user to the relevant page for the snippet.
	// Change the redirect to use the new semantic URL style of /snippet/:id from /snippet?id=
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

	//w.Write([]byte("Create a new specific snippet"))
}
