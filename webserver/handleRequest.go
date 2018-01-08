package webserver

import (
	"net/http"
	"log"
	"html/template"
	"path/filepath"
	"fmt"
	"strconv"
	"github.com/kherud/goblog/config"
	"github.com/kherud/goblog/backend"
	"github.com/kherud/goblog/backend/models"
)

/**
Starts the web server using https. Declares different handlers for static files, page & login/-out requests
 */
func StartServer() {
	fs := http.FileServer(http.Dir(config.STATIC_FILE_PATH))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", returnContent)
	http.HandleFunc("/login", loginUser)
	http.HandleFunc("/logout", logoutUser)
	err := http.ListenAndServeTLS(":"+config.DEFAULT_PORT, config.CERT_FILE, config.KEY_FILE, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

/**
Routes requests appropriate to their GET parameter. If none or an arbitrary one exists the index page is returned.
Therefor hands over necessary variables to 'assembleTemplate()': responseWriter, request, loginRequired, templateName, requestedPage and parameter (GET value).
Since the pages mostly consist of two templates with one of them being the static content (header, footer, ...) only one template name is passed that determines the dynamic main content.
The parameter value is used to pass GET values (e.g. post index id)
 */
func returnContent(w http.ResponseWriter, r *http.Request) {
	parameters := r.URL.Query()
	for key := range parameters {
		switch key {
		case "id": // Display a whole post (param: post id)
			assembleTemplate(w, r, false, "post.html", "post", parameters.Get("id"))
		case "comment": // Persists a comment then displays the appropriate post (param: post id belonging to the comment)
			backend.SaveComment(r, parameters.Get("comment"))
			http.Redirect(w, r, "https://"+r.Host+"?id="+parameters.Get("comment"), http.StatusMovedPermanently)
		case "post": // Displays the post creation site
			assembleTemplate(w, r, true, "createPost.html", "create", "")
		case "account": // Displays the account page (change password / create user if admin)
			assembleTemplate(w, r, true, "user.html", "user", "")
		case "search": // Displays the index page with results filtered by a keyword (param: keyword)
			assembleTemplate(w, r, false, "postPreview.html", "index", parameters.Get("search"))
		case "more": // Responses to an ajax request to load more posts (default 5) if there are more (param: amount of posts already displayed -> index)
			assembleSingleTemplate(w, r, "postPreview.html", "more", parameters.Get("more"))
		case "newPost": // Tries to persist a new post and shows it if successful. Otherwise returns to the post creation site.
			id := backend.CreatePost(r)
			if id == 0 {
				assembleTemplate(w, r, true, "createPost.html", "create", "")
			} else {
				http.Redirect(w, r, "https://"+r.Host+"?id="+strconv.Itoa(int(id)), http.StatusMovedPermanently)
			}
		case "newUser": // Ajax request to persist an user. Returns the username or an error message (e.g. 'Konstantin#' or '#Error Message')
			name, err := backend.CreateUser(r)
			w.Write([]byte(name + "#" + err))
		case "delete": // Ajax request to delete a post. Returns a string representing the success bool value.
			w.Write([]byte(strconv.FormatBool(backend.DeletePost(r))))
		case "edit": // Shows the post editing page (param: post id)
			assembleTemplate(w, r, true, "editPost.html", "post", parameters.Get("edit"))
		case "update": // Tries to apply edits to a post and then returns it (param: post id)
			id := parameters.Get("update")
			backend.UpdatePost(r, id)
			http.Redirect(w, r, "https://"+r.Host+"?id="+id, http.StatusMovedPermanently)
		case "password": // Ajax request to change a password. Possibly returns an error message.
			err := backend.ChangePassword(r)
			w.Write([]byte(err))
		case "verify": // Ajax request to verify a comment. Returns a string representing the success bool value.
			success := backend.VerifyComment(r)
			w.Write([]byte(strconv.FormatBool(success)))
		default: // Returns the index page (listing of recent posts) for unknown GET parameters.
			assembleTemplate(w, r, false, "postPreview.html", "index", "")
		}
		return // only consider first parameter
	}
	// If no GET parameter exists the index page is returned.
	assembleTemplate(w, r, false, "postPreview.html", "index", "")
}

/**
Assembles an html page.
Checks if a login is necessary to view the page / if the user is logged in. If the user lacks access he is redirected to the index page.
Otherwise inserts the dynamic content (templateName) into the static template (header, footer, ...).
Therefor appropriate page variables are loaded that always include information about an existing authentication.
Then returns the result of the assembled html template.
 */
func assembleTemplate(w http.ResponseWriter, r *http.Request, loginRequired bool, templateName, page, parameter string) (int, error) {
	if _, loggedIn := backend.CheckAuthentication(r); loginRequired && !loggedIn {
		http.Redirect(w, r, "https://"+r.Host, http.StatusMovedPermanently)
		return 401, nil
	}
	staticContent := filepath.Join(config.TEMPLATE_PATH, "index.html")
	dynamicContent := filepath.Join(config.TEMPLATE_PATH, templateName)
	tmpl, err := template.ParseFiles(staticContent, dynamicContent)
	if err != nil {
		fmt.Println(err)
		return 500, err
	}

	entries := getPageVars(page, r, parameter)
	tmpl.Execute(w, entries)
	return 0, nil
}

/**
Inserts variables into a single template (always postPreview.html).
Only used to answer the ajax request for loading more posts.
The parameter value is used to identify the index of the requested posts.
 */
func assembleSingleTemplate(w http.ResponseWriter, r *http.Request, templateName, page, parameter string) (int, error) {
	lp := filepath.Join(config.TEMPLATE_PATH, templateName)
	tmpl, err := template.New("mainContent").ParseFiles(lp)
	if err != nil {
		fmt.Println(err)
		return 500, err
	}
	entries := getPageVars(page, r, parameter)
	tmpl.Execute(w, entries)
	return 0, nil
}

/**
Loads information from the backend appropriate to a requested site.
The parameter (GET) value is used to give required information (e.g. id of the requested post)
 */
func getPageVars(page string, r *http.Request, parameter string) map[string]interface{} {
	entries := map[string]interface{}{}
	switch page {
	case "index":
		entries = getIndexVars(parameter)
	case "more":
		entries = getLoadMoreVars(parameter)
	case "post":
		// if an error occurs this variable is empty -> 404 message displayed, e.g. /?id=0
		entries["post"], _ = backend.GetPost(parameter)
	}
	// If the request reveals an existing session further information about the user is provided
	if user, found := backend.CheckAuthentication(r); found {
		entries["user"] = user
	}
	return entries
}

/**
Loads required variables of the index page, including:
- initial: are the displayed posts the first ones? (<-> load more)
- previews: information about posts to assemble their preview within the index page.
- search: search phrase / keyword defined by the GET parameter. If set the posts are appropriately filtered.
- more/index: if more posts exist a flag is set and the index to load more content is provided.
 */
func getIndexVars(parameter string) map[string]interface{} {
	entries := map[string]interface{}{}
	entries["initial"] = true
	entries["previews"] = backend.GetEntries()
	if parameter != "" {
		entries["previews"] = backend.FilterPosts(entries["previews"].([]models.Entry), parameter)
		entries["search"] = parameter
	} else {
		if length := len(entries["previews"].([]models.Entry)); length > config.POSTS_PER_REQUESTS {
			entries["previews"] = entries["previews"].([]models.Entry)[:config.POSTS_PER_REQUESTS]
			entries["more"] = true
			entries["index"] = config.POSTS_PER_REQUESTS
		}
	}
	return entries
}

/**
Loads required variables to reponse to an ajax request to load more posts.
- initial: are the displayed posts the first ones? (<-> load more)
- previews: information about posts to assemble their preview within the index page.
- search: search phrase / keyword defined by the GET parameter. If set the posts are appropriately filtered.
- more/index: if more posts exist a flag is set and the index to load more content is provided.
 */
func getLoadMoreVars(parameter string) map[string]interface{} {
	entries := map[string]interface{}{}
	entries["previews"] = backend.GetEntries()
	index, _ := strconv.Atoi(parameter)
	lengthLeft := len(entries["previews"].([]models.Entry)) - index
	if lengthLeft >= config.POSTS_PER_REQUESTS {
		entries["previews"] = entries["previews"].([]models.Entry)[index:index+config.POSTS_PER_REQUESTS]
		if lengthLeft > config.POSTS_PER_REQUESTS {
			entries["more"] = true
			entries["index"] = index + config.POSTS_PER_REQUESTS
		}
	} else {
		entries["previews"] = entries["previews"].([]models.Entry)[index:]
	}
	return entries
}

/**
Validates the transferred credentials and sets a session if successful.
 */
func loginUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if backend.AuthenticateUser(username, password) {
		fmt.Println("Login:", username)
		backend.SetSession(username, w)
		w.Write([]byte("success"))
	} else {
		w.Write([]byte("failed to login"))
	}
}

/**
Eliminates the active session and redirects to the index page.
 */
func logoutUser(w http.ResponseWriter, r *http.Request) {
	backend.EndSession(r)
	http.Redirect(w, r, "https://"+r.Host, http.StatusMovedPermanently)
	return
}
