package backend

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"net/url"
	"net/http"
	"encoding/json"
	"path/filepath"
	"github.com/kherud/goblog/backend/models"
	"github.com/kherud/goblog/config"
)

var testEntry = models.Entry{
	Title:    "Test",
	Text:     "Test",
	Author:   "Test",
	AuthorId: 689017489,
	Date:     "Test",
	Id:       976620356,
	Comments: []models.Comment{{Text: "cTest1", Author: "cTest2", Date: "cTest3", Verified: false, Id: 489017489}},
	Keywords: []string{"abd", "def"},
}

func TestSaveCommentInvalid(t *testing.T) {
	config.ENTRIES_FILE_PATH = config.TEST_TEMP_PATH
	saveEntriesJson([]models.Entry{testEntry})
	tests := []struct {
		Params url.Values
	}{{Params: url.Values{"text": {""}, "name": {"TestTestTest"}}},
		{Params: url.Values{}},
	}
	for _, test := range tests {
		req := &http.Request{
			Form:   test.Params,
			Header: http.Header{},
		}
		SaveComment(req, "976620356")
		post, err := GetPost("976620356")
		assert.Nil(t, err)
		assert.True(t, post.Id > 0)
		assert.True(t, len(post.Comments) == 1)
		assert.True(t, post.Comments[0].Id == 489017489)
	}
	os.Remove(config.TEST_TEMP_PATH)
	config.ENTRIES_FILE_PATH = filepath.Join("test_data", "entries.json")
}

func TestSaveCommentValid(t *testing.T) {
	config.ENTRIES_FILE_PATH = config.TEST_TEMP_PATH
	saveEntriesJson([]models.Entry{testEntry})
	tests := []struct {
		Params url.Values
	}{{Params: url.Values{"text": {"Test"}, "name": {"TestTestTest"}}},
		{Params: url.Values{"text": {"Test"}, "name": {""}}},
	}
	expectedNames := []string{"TestTestTest", "Anonymous"}
	for idx, test := range tests {
		req := &http.Request{
			Form:   test.Params,
			Header: http.Header{},
		}
		SaveComment(req, "976620356")
		post, err := GetPost("976620356")
		assert.Nil(t, err)
		assert.True(t, post.Id > 0)
		assert.True(t, len(post.Comments) == 2+idx)
		// Index 0 since every new comments must be prepended to be displayed at the top
		assert.True(t, post.Comments[0].Id != 489017489)
		assert.True(t, post.Comments[0].Id > 0)
		assert.EqualValues(t, post.Comments[0].Author, expectedNames[idx])
		assert.EqualValues(t, post.Comments[0].Text, "Test")
		assert.NotEmpty(t, post.Comments[0].Date)
		assert.False(t, post.Comments[0].Verified)
	}
	os.Remove(config.TEST_TEMP_PATH)
	config.ENTRIES_FILE_PATH = filepath.Join("test_data", "entries.json")
}

func TestGetPost(t *testing.T) {
	post, err := GetPost("708643541")
	assert.Nil(t, err)
	assert.True(t, post.Id == 708643541)
	post, err = GetPost("")
	assert.True(t, err != nil)
	assert.True(t, post.Id == 0)
}

func TestVerifyCommentInvalid(t *testing.T) {
	config.ENTRIES_FILE_PATH = config.TEST_TEMP_PATH
	saveEntriesJson([]models.Entry{testEntry})
	tests := []struct {
		Params url.Values
	}{{Params: url.Values{"postId": {"976620356"}, "commentId": {"489017489"}}},
		{Params: url.Values{"postId": {"976620356"}, "commentId": {}}},
		{Params: url.Values{"postId": {""}, "commentId": {"489017489"}}},
		{Params: url.Values{}},
	}
	for idx, test := range tests {
		req := &http.Request{
			Form:   test.Params,
			Header: http.Header{},
		}
		if idx > 0 {
			cookie := &http.Cookie{Name: "Session", Value: "Konstantin#Test"}
			req.AddCookie(cookie)
		}
		verified := VerifyComment(req)
		assert.False(t, verified)
		post, err := GetPost("976620356")
		assert.Nil(t, err)
		assert.False(t, post.Comments[0].Verified)
	}
	os.Remove(config.TEST_TEMP_PATH)
	config.ENTRIES_FILE_PATH = filepath.Join("test_data", "entries.json")
}

func TestVerifyCommentValid(t *testing.T) {
	config.ENTRIES_FILE_PATH = config.TEST_TEMP_PATH
	saveEntriesJson([]models.Entry{testEntry})
	req := &http.Request{
		Form:   url.Values{"postId": {"976620356"}, "commentId": {"489017489"}},
		Header: http.Header{},
	}
	cookie := &http.Cookie{Name: "Session", Value: "Konstantin#Test"}
	req.AddCookie(cookie)
	verified := VerifyComment(req)
	assert.True(t, verified)
	post, err := GetPost("976620356")
	assert.Nil(t, err)
	assert.True(t, post.Comments[0].Verified)
	os.Remove(config.TEST_TEMP_PATH)
	config.ENTRIES_FILE_PATH = filepath.Join("test_data", "entries.json")
}

func TestCreatePostInvalid(t *testing.T) {
	config.ENTRIES_FILE_PATH = config.TEST_TEMP_PATH
	saveEntriesJson([]models.Entry{testEntry})
	tests := []struct {
		Params url.Values
	}{{Params: url.Values{"text": {"asd"}, "title": {"asd"}, "tag": {"asd"}}},
		{Params: url.Values{"text": {""}, "title": {"asd"}, "tag": {"asd"}}},
		{Params: url.Values{}},
	}
	for idx, test := range tests {
		req := &http.Request{
			Form:   test.Params,
			Header: http.Header{},
		}
		if idx > 0 {
			cookie := &http.Cookie{Name: "Session", Value: "Konstantin#Test"}
			req.AddCookie(cookie)
		}
		postId := CreatePost(req)
		assert.True(t, postId == 0)
		posts := GetEntries()
		assert.True(t, len(posts) == 1)
	}
	os.Remove(config.TEST_TEMP_PATH)
	config.ENTRIES_FILE_PATH = filepath.Join("test_data", "entries.json")
}

func TestCreatePostValid(t *testing.T) {
	config.ENTRIES_FILE_PATH = config.TEST_TEMP_PATH
	saveEntriesJson([]models.Entry{testEntry})
	tests := []struct {
		Params url.Values
	}{{Params: url.Values{"text": {"Test"}, "title": {"Test"}, "tag": {"Test"}}},
		{Params: url.Values{"text": {"Test"}, "title": {""}, "tag": {"Test", "Test2"}}},
	}
	expectedTitles := []string{"Test", "Post #3"}
	for idx, test := range tests {
		req := &http.Request{
			Form:   test.Params,
			Header: http.Header{},
		}
		cookie := &http.Cookie{Name: "Session", Value: "Konstantin#Test"}
		req.AddCookie(cookie)
		postId := CreatePost(req)
		assert.True(t, postId > 0)
		posts := GetEntries()
		assert.True(t, len(posts) == 2 + idx)
		assert.NotEmpty(t, posts[0].Date)
		assert.True(t, posts[0].AuthorId > 0)
		assert.EqualValues(t, posts[0].Id, postId)
		assert.EqualValues(t, posts[0].Text, "Test")
		assert.EqualValues(t, posts[0].Title, expectedTitles[idx])
		assert.EqualValues(t, posts[0].Keywords, tests[idx].Params["tag"])
	}
	os.Remove(config.TEST_TEMP_PATH)
	config.ENTRIES_FILE_PATH = filepath.Join("test_data", "entries.json")
}

func TestDeletePostSingle (t *testing.T) {
	config.ENTRIES_FILE_PATH = config.TEST_TEMP_PATH
	saveEntriesJson([]models.Entry{testEntry})
	req := &http.Request{
		Form:   url.Values{"postId": {"976620356"}},
		Header: http.Header{},
	}
	cookie := &http.Cookie{Name: "Session", Value: "Konstantin#Test"}
	req.AddCookie(cookie)
	DeletePost(req)
	entries := GetEntries()
	assert.True(t, len(entries) == 0)
	os.Remove(config.TEST_TEMP_PATH)
	config.ENTRIES_FILE_PATH = filepath.Join("test_data", "entries.json")
}

func TestDeletePostCorrectInMultiple(t *testing.T) {
	config.ENTRIES_FILE_PATH = config.TEST_TEMP_PATH
	entry2 := testEntry
	entry2.Id = 489017489
	saveEntriesJson([]models.Entry{testEntry, entry2})
	req := &http.Request{
		Form:   url.Values{"postId": {"976620356"}},
		Header: http.Header{},
	}
	cookie := &http.Cookie{Name: "Session", Value: "Konstantin#Test"}
	req.AddCookie(cookie)
	DeletePost(req)
	entries := GetEntries()
	assert.True(t, len(entries) == 1)
	assert.EqualValues(t, entries[0].Id, 489017489)
	os.Remove(config.TEST_TEMP_PATH)
	config.ENTRIES_FILE_PATH = filepath.Join("test_data", "entries.json")
}

func TestDeletePostInvalid(t *testing.T) {
	config.ENTRIES_FILE_PATH = config.TEST_TEMP_PATH
	saveEntriesJson([]models.Entry{testEntry})
	req := &http.Request{
		Form:   url.Values{"postId": {"976620356"}},
		Header: http.Header{},
	}
	DeletePost(req)
	entries := GetEntries()
	assert.True(t, len(entries) == 1)
	assert.EqualValues(t, entries[0].Id, testEntry.Id)
	os.Remove(config.TEST_TEMP_PATH)
	config.ENTRIES_FILE_PATH = filepath.Join("test_data", "entries.json")
}

func TestFilterPosts(t *testing.T){
	entries := GetEntries()
	assert.True(t, len(entries) == 7)
	filters := []string{"", "asd", "cde", "Test"}
	entriesFilter1 := FilterPosts(entries, filters[0])
	entriesFilter2 := FilterPosts(entries, filters[1])
	entriesFilter3 := FilterPosts(entries, filters[2])
	entriesFilter4 := FilterPosts(entries, filters[3])
	assert.True(t, len(entriesFilter1) == 0) // empty tags aren't possible -> no post should contain one
	assert.True(t, len(entriesFilter2) == 2)
	assert.True(t, len(entriesFilter3) == 1)
	assert.True(t, len(entriesFilter4) == 0)
}

func TestUpdatePostInvalid(t *testing.T){
	config.ENTRIES_FILE_PATH = config.TEST_TEMP_PATH
	saveEntriesJson([]models.Entry{testEntry})
	for idx := 0; idx < 2; idx++ {
		req := &http.Request{
			Form:   url.Values{"text": {"Test2"}, "title": {"Test2"}, "tag": {"Test2"}},
			Header: http.Header{},
		}
		if idx > 0 {
			cookie := &http.Cookie{Name: "Session", Value: "Konstanti#Test"} // different authorId for 'Konstanti'
			req.AddCookie(cookie)
		}
		updated := UpdatePost(req, "976620356")
		assert.False(t, updated)
	}
	os.Remove(config.TEST_TEMP_PATH)
	config.ENTRIES_FILE_PATH = filepath.Join("test_data", "entries.json")
}

func TestUpdatePostValid(t *testing.T){
	config.ENTRIES_FILE_PATH = config.TEST_TEMP_PATH
	saveEntriesJson([]models.Entry{testEntry})
	req := &http.Request{
		Form:   url.Values{"text": {"Test2"}, "title": {"Test2"}, "tag": {"Test2"}},
		Header: http.Header{},
	}
	cookie := &http.Cookie{Name: "Session", Value: "Konstantin#Test"}
	req.AddCookie(cookie)
	updated := UpdatePost(req, "976620356")
	assert.True(t, updated)
	post, err := GetPost("976620356")
	assert.Nil(t, err)
	assert.NotEqual(t, post.Text, testEntry.Text)
	assert.NotEqual(t, post.Title, testEntry.Title)
	assert.NotEqual(t, post.Keywords, testEntry.Keywords)
	assert.NotEqual(t, post.Date, testEntry.Date)
	assert.Equal(t, post.Id, testEntry.Id)
	os.Remove(config.TEST_TEMP_PATH)
	config.ENTRIES_FILE_PATH = filepath.Join("test_data", "entries.json")
}

func TestAssemblePost(t *testing.T) {
	tests := []struct {
		Params url.Values
	}{{Params: url.Values{"text": {"Test"}, "title": {"Test"}, "tag": {"Test"}}},
		{Params: url.Values{"text": {"Test"}, "title": {""}, "tag": {"Test", "Test2"}}},
	}
	expectedTitles := []string{"Test", "Post #8"}
	for idx, test := range tests {
		req := &http.Request{
			Form:   test.Params,
		}
		user, err := GetUser("Konstantin")
		assert.Nil(t, err)
		post := assemblePost(req, user)
		assert.True(t, post.Id > 0)
		assert.NotEmpty(t, post.Date)
		assert.EqualValues(t, post.AuthorId, user.Id)
		assert.EqualValues(t, post.Text, "Test")
		assert.EqualValues(t, post.Title, expectedTitles[idx])
		assert.EqualValues(t, post.Keywords, tests[idx].Params["tag"])
	}
}

func TestSavePost(t *testing.T){
	config.ENTRIES_FILE_PATH = config.TEST_TEMP_PATH
	savePost(testEntry)
	entries := testEntriesFileExistsGetContent(t)
	assert.EqualValues(t, entries[0].Id, 976620356)
	assert.True(t, len(entries) == 1)
	entry2 := testEntry
	entry2.Id = 489017489
	savePost(entry2)
	entries = testEntriesFileExistsGetContent(t)
	assert.EqualValues(t, entries[0].Id, 489017489)
	assert.EqualValues(t, entries[1].Id, 976620356)
	assert.True(t, len(entries) == 2)
	os.Remove(config.TEST_TEMP_PATH)
	config.ENTRIES_FILE_PATH = filepath.Join("test_data", "entries.json")
}

func testEntriesFileExistsGetContent(t *testing.T) []models.Entry {
	_, err := os.Stat(config.TEST_TEMP_PATH)
	assert.True(t, err == nil)
	raw := readFile(config.TEST_TEMP_PATH)
	var validationInstance []models.Entry
	json.Unmarshal(raw, &validationInstance)
	return validationInstance
}
