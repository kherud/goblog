package backend

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"unicode/utf8"
	"os"
	"encoding/json"
	"path/filepath"
	"github.com/kherud/goblog/config"
	"github.com/kherud/goblog/backend/models"
)

func TestMain(m *testing.M) {
	config.USERS_FILE_PATH = config.USERS_TEST_PATH
	config.ENTRIES_FILE_PATH = config.ENTRIES_TEST_PATH
	config.MIN_USERNAME_LENGTH = 6
	config.MIN_PASSWORD_LENGTH = 8
	os.Exit(m.Run())
}

func TestGetUsersSafe(t *testing.T) {
	users := GetUsers()
	assert.True(t, len(users) == 3)
	for _, user := range users {
		assert.True(t, utf8.RuneCountInString(user.UserName) > config.MIN_USERNAME_LENGTH)
		assert.True(t, utf8.RuneCountInString(user.Password) > 30)
		assert.True(t, user.Id != 0)
	}
}

func TestGetUsersCorrupted(t *testing.T) {
	config.USERS_FILE_PATH = filepath.Join("test_data", "users_corrupted.json")
	users := GetUsers()
	assert.True(t, len(users) == 0)
	config.USERS_FILE_PATH = config.USERS_TEST_PATH
}

func TestGetEntriesSafe(t *testing.T) {
	entries := GetEntries()
	assert.True(t, len(entries) == 7)
	for _, entry := range entries {
		assert.True(t, utf8.RuneCountInString(entry.Author) > 0)
		assert.True(t, utf8.RuneCountInString(entry.Date) > 0)
		assert.True(t, entry.AuthorId != 0)
		assert.True(t, entry.Id != 0)
	}
}

func TestGetEntriesCorrupted(t *testing.T) {
	config.ENTRIES_FILE_PATH = filepath.Join("test_data", "entries_corrupted.json")
	entries := GetEntries()
	assert.True(t, len(entries) == 0)
	config.ENTRIES_FILE_PATH = config.ENTRIES_TEST_PATH
}

func TestSaveUsersJson(t *testing.T) {
	config.USERS_FILE_PATH = config.TEST_TEMP_PATH
	testUser := models.User{
		UserName: "Test1",
		Password: "Test2",
		Id:       689017489,
		Session:  "Test3",
		Admin:    true,
	}
	saveUsersJson([]models.User{testUser})
	_, err := os.Stat(config.TEST_TEMP_PATH)
	assert.Nil(t, err)
	raw := readFile(config.TEST_TEMP_PATH)
	var validationInstance []models.User
	json.Unmarshal(raw, &validationInstance)
	assert.EqualValues(t, testUser.UserName, validationInstance[0].UserName)
	assert.EqualValues(t, testUser.Password, validationInstance[0].Password)
	assert.EqualValues(t, testUser.Id, validationInstance[0].Id)
	assert.EqualValues(t, testUser.Session, validationInstance[0].Session)
	assert.EqualValues(t, testUser.Admin, validationInstance[0].Admin)
	os.Remove(config.TEST_TEMP_PATH)
	config.USERS_FILE_PATH = filepath.Join("test_data", "users.json")
}

func TestSaveEntriesJson(t *testing.T) {
	config.ENTRIES_FILE_PATH = config.TEST_TEMP_PATH
	testEntry := models.Entry{
		Title:    "Test1",
		Text:     "Test2",
		Author:   "Test3",
		AuthorId: 689017489,
		Date:     "Test4",
		Id:       589017489,
		Comments: []models.Comment{{Text: "cTest1", Author: "cTest2", Date: "cTest3", Verified: false, Id: 489017489}},
		Keywords: []string{"abc", "def"},
	}
	saveEntriesJson([]models.Entry{testEntry})
	_, err := os.Stat(config.TEST_TEMP_PATH)
	assert.Nil(t, err)
	raw := readFile(config.TEST_TEMP_PATH)
	var validationInstance []models.Entry
	json.Unmarshal(raw, &validationInstance)
	assert.EqualValues(t, testEntry.Title, validationInstance[0].Title)
	assert.EqualValues(t, testEntry.Text, validationInstance[0].Text)
	assert.EqualValues(t, testEntry.Author, validationInstance[0].Author)
	assert.EqualValues(t, testEntry.AuthorId, validationInstance[0].AuthorId)
	assert.EqualValues(t, testEntry.Date, validationInstance[0].Date)
	assert.EqualValues(t, testEntry.Id, validationInstance[0].Id)
	assert.EqualValues(t, testEntry.Comments, validationInstance[0].Comments)
	assert.EqualValues(t, testEntry.Keywords, validationInstance[0].Keywords)
	os.Remove(config.TEST_TEMP_PATH)
	config.USERS_FILE_PATH = filepath.Join("test_data", "entries.json")
}
