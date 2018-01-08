package backend

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"encoding/json"
	"unicode/utf8"
	"net/url"
	"net/http"
	"path/filepath"
	"github.com/kherud/goblog/backend/models"
	"github.com/kherud/goblog/config"
	"github.com/kherud/goblog/util"
)

type testReader struct {
	text []rune
}

func (tr testReader) ReadString(byte) (string, error) {
	return string(tr.text), nil
}

var testUser = models.User{
	UserName: "Test",
	Password: "inGPp5bFPgeeB6vp6p3_ECLirbGb9LKNeFPS9tAuAW8=",
	Id:       689017489,
	Session:  "Test",
	Admin:    false,
}

func TestEnsureUserExists(t *testing.T) {
	config.USERS_FILE_PATH = config.TEST_TEMP_PATH
	reader := testReader{text: []rune("TestTestTest")}
	EnsureUserExists(reader)
	users := testFileExistsGetContent(t)
	assert.True(t, len(users) == 1)
	assert.True(t, users[0].Admin)
	assert.True(t, utf8.RuneCountInString(users[0].Password) > 30)
	assert.EqualValues(t, users[0].UserName, "TestTestTest")
	assert.True(t, users[0].Id > 0)
	id := users[0].Id
	EnsureUserExists(reader)
	users = testFileExistsGetContent(t)
	assert.True(t, len(users) == 1)
	assert.True(t, users[0].Id == id)
	os.Remove(config.TEST_TEMP_PATH)
	config.USERS_FILE_PATH = filepath.Join("test_data", "users.json")
}

func TestGetUser(t *testing.T) {
	user, err := GetUser("Konstantin")
	assert.Nil(t, err)
	assert.True(t, user.Id == 689017489)
	user, err = GetUser("")
	assert.True(t, err != nil)
	assert.EqualValues(t, user.Id,0)
}

func TestAuthenticateUser(t *testing.T) {
	authentication_valid := AuthenticateUser("Konstantin", "12345678")
	authentication_invalid := AuthenticateUser("nitnatsnoK", "12345678")
	authentication_invalid2 := AuthenticateUser("Konstantin", "87654321")
	authentication_invalid3 := AuthenticateUser("", "")
	assert.True(t, authentication_valid)
	assert.False(t, authentication_invalid)
	assert.False(t, authentication_invalid2)
	assert.False(t, authentication_invalid3)
}

func TestPersistSession(t *testing.T) {
	config.USERS_FILE_PATH = config.TEST_TEMP_PATH
	EnsureUserExists(testReader{text: []rune("TestTestTest")})
	users := testFileExistsGetContent(t)
	assert.True(t, users[0].Session == "")
	persistSession("TestTestTest", "1234567890")
	users = testFileExistsGetContent(t)
	assert.True(t, users[0].Session == "1234567890")
	os.Remove(config.TEST_TEMP_PATH)
	config.USERS_FILE_PATH = filepath.Join("test_data", "users.json")
}

func TestCompareCredentials(t *testing.T) {
	authentication_valid := compareCredentials(testUser, "Test", "Test")
	authentication_invalid := compareCredentials(testUser, "Test", "")
	authentication_invalid2 := compareCredentials(testUser, "", "Test")
	authentication_invalid3 := compareCredentials(models.User{}, "Test", "Test")
	assert.True(t, authentication_valid)
	assert.False(t, authentication_invalid)
	assert.False(t, authentication_invalid2)
	assert.False(t, authentication_invalid3)
}

func TestCreateInitialUser(t *testing.T) {
	user := createInitialUser(testReader{text: []rune("TestTestTest")})
	assert.True(t, user.Id > 0)
	assert.True(t, user.Admin)
	assert.EqualValues(t, user.UserName, "TestTestTest")
	assert.EqualValues(t, user.Password, util.HashPassword("TestTestTest", user.Id))
	assert.Empty(t, user.Session)
}

func TestUpdateUser(t *testing.T) {
	users := []models.User{testUser}
	users = updateUsers(users, models.User{Id: 689017489})
	assert.Empty(t, users[0].UserName)
	assert.Empty(t, users[0].Password)
	assert.Empty(t, users[0].Session)
	assert.False(t, users[0].Admin)
}

func TestSaveUser(t *testing.T) {
	config.USERS_FILE_PATH = config.TEST_TEMP_PATH
	saveUser(testUser)
	validationInstance := testFileExistsGetContent(t)
	assert.True(t, len(validationInstance) == 1)
	testUser.UserName = "TestTest"
	saveUser(testUser)
	validationInstance = testFileExistsGetContent(t)
	assert.EqualValues(t, validationInstance[0].UserName, "TestTest")
	os.Remove(config.TEST_TEMP_PATH)
	config.USERS_FILE_PATH = filepath.Join("test_data", "users.json")
}

func TestCreateUserInvalidForm(t *testing.T) {
	tests := []struct {
		Params url.Values
	}{{Params: url.Values{"name": {"TestTestTest"}, "password": {"TestTestTest"}, "password-confirmation": {"TestTestTest"}}},
		{Params: url.Values{"name": {""}, "password": {"TestTestTest"}, "password-confirmation": {"TestTestTest"}}},
		{Params: url.Values{"name": {"TestTestTest"}, "password": {"TestTestTest"}, "password-confirmation": {""}}},
		{Params: url.Values{"name": {"TestTestTest"}, "password": {""}, "password-confirmation": {"TestTestTest"}}},
		{Params: url.Values{"name": {"Konstantin"}, "password": {""}, "password-confirmation": {"TestTestTest"}}},
		{Params: url.Values{}},
	}

	for idx, test := range tests {
		req := &http.Request{
			Form:   test.Params,
			Header: http.Header{},
		}
		// check login validation in first case
		if idx > 0 {
			cookie := &http.Cookie{Name: "Session", Value: "Konstantin#Test"}
			req.AddCookie(cookie)
		}
		user, err := CreateUser(req)
		assert.Empty(t, user)
		assert.NotEmpty(t, err)
	}
}

func TestCreateUserValidForm(t *testing.T) {
	req := &http.Request{
		Form:   url.Values{"name": {"TestTestTest"}, "password": {"TestTestTest"}, "password-confirmation": {"TestTestTest"}},
		Header: http.Header{},
	}
	cookie := &http.Cookie{Name: "Session", Value: "Konstantin#Test"}
	req.AddCookie(cookie)
	user, err := CreateUser(req)
	assert.NotEmpty(t, user)
	assert.Empty(t, err)
	// delete user for future tests
	users := GetUsers()
	users = users[:len(users) - 1]
	saveUsersJson(users)
}

func TestCreateUserValidFormAdmin(t *testing.T) {
	req := &http.Request{
		Form:   url.Values{"name": {"TestTestTest"}, "password": {"TestTestTest"}, "password-confirmation": {"TestTestTest"}, "admin": {"on"}},
		Header: http.Header{},
	}
	cookie := &http.Cookie{Name: "Session", Value: "Konstantin#Test"}
	req.AddCookie(cookie)
	userName, errMsg := CreateUser(req)
	assert.NotEmpty(t, userName)
	assert.Empty(t, errMsg)
	user, err := GetUser(userName)
	assert.Nil(t, err)
	assert.True(t, user.Admin)
	// delete user for future tests
	users := GetUsers()
	users = users[:len(users) - 1]
	saveUsersJson(users)
}

func TestChangePasswordInvalid(t *testing.T){
	tests := []struct {
		Params url.Values
	}{{Params: url.Values{"password": {"TestTestTest"}, "password-confirmation": {"TestTestTest"}}},
		{Params: url.Values{"password": {"Test"}, "password-confirmation": {"Test"}}},
		{Params: url.Values{"password": {"TestTestTest"}, "password-confirmation": {""}}},
		{Params: url.Values{"password": {""}, "password-confirmation": {"TestTestTest"}}},
		{Params: url.Values{}},
	}

	for idx, test := range tests {
		req := &http.Request{
			Form:   test.Params,
			Header: http.Header{},
		}
		// check login validation in first case
		if idx > 0 {
			cookie := &http.Cookie{Name: "Session", Value: "Konstantin#Test"}
			req.AddCookie(cookie)
		}
		err := ChangePassword(req)
		assert.NotEmpty(t, err)
	}
}

func TestChangePasswordValidForm(t *testing.T) {
	req := &http.Request{
		Form:   url.Values{"password": {"TestTestTest"}, "password-confirmation": {"TestTestTest"}},
		Header: http.Header{},
	}
	cookie := &http.Cookie{Name: "Session", Value: "Konstantin#Test"}
	req.AddCookie(cookie)
	err := ChangePassword(req)
	assert.Empty(t, err)
	user, _ := GetUser("Konstantin")
	oldPassword := util.HashPassword("12345678", user.Id)
	assert.NotEqual(t, user.Password, oldPassword)
	user.Password = oldPassword
	users := GetUsers()
	users = updateUsers(users, user)
	saveUsersJson(users)
}

func testFileExistsGetContent(t *testing.T) []models.User {
	_, err := os.Stat(config.TEST_TEMP_PATH)
	assert.True(t, err == nil)
	raw := readFile(config.TEST_TEMP_PATH)
	var validationInstance []models.User
	json.Unmarshal(raw, &validationInstance)
	return validationInstance
}

