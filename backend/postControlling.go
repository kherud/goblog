package backend

import (
	"time"
	"errors"
	"strconv"
	"net/http"
	"fmt"
	"unicode/utf8"
	"github.com/kherud/goblog/backend/models"
	"github.com/kherud/goblog/util"
)

/**
Extracts a comment from the POST form of an http(s) request.
Saves the comment within the post of the transferred post id if all requirements are met.
Comments are always prepended to the comment slice of the post. Thus they are chronologically displayed.
 */
func SaveComment(r *http.Request, postId string) {
	if r.FormValue("text") == "" {
		return
	}
	var author string
	if r.FormValue("name") == "" {
		author = "Anonymous"
	} else {
		author = r.FormValue("name")
	}
	uintPostId, _ := strconv.ParseUint(postId, 10, 32) // parse string to uint32
	entries := GetEntries()
	for idx, post := range entries {
		if post.Id == uint32(uintPostId) {
			date := time.Now().Local().Format("02.01.2006 - 15:04")
			comment := models.Comment{
				Text:   r.FormValue("text"),
				Author: author,
				Date:   date,
				Id:     util.CreateHashId(date, author, r.FormValue("text")),
			}
			entries[idx].Comments = append([]models.Comment{comment}, post.Comments...) // prepend
			saveEntriesJson(entries)
		}
	}
}

/**
Returns a single post by iterating over all existing posts and comparing their ids.
If the post is found it is returned without an error.
Otherwise an empty instance with an appropriate error is returned.
 */
func GetPost(id string) (post models.Entry, err error) {
	entries := GetEntries()
	for _, entry := range entries {
		if strconv.Itoa(int(entry.Id)) == id {
			return entry, nil
		}
	}
	return models.Entry{}, errors.New("entry not found")
}

/**
Changes the status of a comment to verified if the request is authenticated and all requirements are met.
Does so by parsing the POST form of a request and extracting the comment id and its affiliated post id.
Returns a boolean that represents whether the comment was successfully verified.
 */
func VerifyComment(r *http.Request) bool {
	_, loggedIn := CheckAuthentication(r)
	var postId, commentId string
	if err := r.ParseForm(); err == nil && loggedIn {
		postId = r.FormValue("postId")
		commentId = r.FormValue("commentId")
	} else {
		return false
	}
	entries := GetEntries()
	for entryIdx, entry := range entries { // search for entry
		if strconv.Itoa(int(entry.Id)) == postId {
			for commentIdx, comment := range entry.Comments { // search for comment
				if strconv.Itoa(int(comment.Id)) == commentId {
					entries[entryIdx].Comments[commentIdx].Verified = true
					saveEntriesJson(entries)
					return true
				}
			}
		}
	}
	return false
}

/**
Extracts a post from the POST form of an http(s) request if the request is authenticated.
If the post was successfully created its id is returned otherwise 0.
 */
func CreatePost(r *http.Request) uint32 {
	user, loggedIn := CheckAuthentication(r)
	if !loggedIn {
		return 0
	}
	post := assemblePost(r, user)
	if post.Id != 0 {
		savePost(post)
	}
	return post.Id
}

/**
Deletes a post by extracting the requested id from the POST form of a http(s) request.
Only does so if the request is authenticated and the user is author of the post.
 */
func DeletePost(r *http.Request) bool {
	user, loggedIn := CheckAuthentication(r)
	if err := r.ParseForm(); err == nil && loggedIn {
		postId := r.FormValue("postId")
		entries := GetEntries()
		for idx, entry := range entries {
			if strconv.Itoa(int(entry.Id)) == postId && entry.AuthorId == user.Id { // check if it is the author's post
				entries = append(entries[:idx], entries[idx+1:]...) // create new slice without element
				saveEntriesJson(entries)
				return true
			}
		}
	}
	return false
}

/**
Filters a slice of posts by checking if their keywords contain a search phrase (keyword).
The filtered slice then is returned.
 */
func FilterPosts(entries []models.Entry, searchterm string) (result []models.Entry) {
	for _, element := range entries {
		for _, keyword := range element.Keywords {
			if keyword == searchterm {
				result = append(result, element)
			}
		}
	}
	return
}

/**
Applies edits to an existing post affiliated to the passed id if the request is authenticated.
Does so by parsing the POST form of an http(s) request.
Also updates the creation date of the post and therefore prepends it to all other existing posts.
 */
func UpdatePost(r *http.Request, postId string) bool {
	user, loggedIn := CheckAuthentication(r)
	if err := r.ParseForm(); err == nil && loggedIn {
		entries := GetEntries()
		for idx, entry := range entries {
			if strconv.Itoa(int(entry.Id)) == postId && entry.AuthorId == user.Id {
				newPost := assemblePost(r, user)
				newPost.Comments = entry.Comments // keep the old comments
				newPost.Id = entry.Id // keep the id
				subSlice := append(entries[:idx], entries[idx+1:]...) // delete old post
				entries = append([]models.Entry{newPost}, subSlice...) // prepend updated post
				saveEntriesJson(entries)
				return true
			}
		}
	}
	return false
}

/**
Creates and returns an instance of Entry by parsing the POST form of an http(s) request.
If no text was transferred an empty instance is returned.
 */
func assemblePost(r *http.Request, user models.User) models.Entry {
	entries := GetEntries()
	date := time.Now().Local().Format("02.01.2006 - 15:04")
	var title string
	if title = r.FormValue("title"); title == "" { // if no title was passed automatically create one
		title = fmt.Sprintf("Post #%v", len(entries)+1)
	}
	if err := r.ParseForm(); err != nil || utf8.RuneCountInString(r.FormValue("text")) == 0 {
		return models.Entry{}
	}
	postId := util.CreateHashId(date, user.UserName, r.FormValue("text"))
	entry := models.Entry{
		Text:     r.FormValue("text"),
		Title:    title,
		Author:   user.UserName,
		AuthorId: user.Id,
		Date:     date,
		Id:       postId,
		Keywords: r.Form["tag"],
	}
	return entry
}

/**
Loads all entries. If none exist so far a new slice with the passed entry is created and saved.
Otherwise the entry is prepended to all other existing posts. Thus all posts are always chronologically displayed.
 */
func savePost(entry models.Entry) {
	entries := GetEntries()
	if entries == nil || len(entries) == 0 {
		entries = []models.Entry{entry}
	} else {
		entries = append([]models.Entry{entry}, entries...) // prepend
	}
	saveEntriesJson(entries)
}
