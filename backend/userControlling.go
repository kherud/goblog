package backend

import (
	"fmt"
	"errors"
	"net/http"
	"unicode/utf8"
	"github.com/kherud/goblog/util"
	"github.com/kherud/goblog/backend/models"
	"github.com/kherud/goblog/config"
)

/**
Ensures a user exists by checking the existence of the users.json and its length.
If empty or non-existent the user is prompted to create an initial account which then is saved.
 */
func EnsureUserExists(reader util.Reader) {
	users := GetUsers()
	if users == nil || len(users) == 0 {
		user := createInitialUser(reader)
		saveUsersJson([]models.User{user})
		fmt.Printf("User '%v' successfully created.\n", user.UserName)
	}
}

/**
Returns a user by his username by iterating over all existing users and comparing their unique usernames.
If the account is not found an error message and empty instance of User is returned.
 */
func GetUser(username string) (user models.User, err error) {
	users := GetUsers()
	for _, user := range users {
		if user.UserName == username {
			return user, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

/**
Validates a transferred username and password by iterating over all existing users and comparing credentials.
Returns a boolean that represents the validity of the credentials.
 */
func AuthenticateUser(username, password string) bool {
	users := GetUsers()
	for _, user := range users {
		if compareCredentials(user, username, password) {
			return true
		}
	}
	return false
}

/**
Searches for an user by his username and stores his session string if he is found.
 */
func persistSession(username, sessionId string) {
	users := GetUsers()
	for _, user := range users {
		if user.UserName == username {
			user.Session = sessionId
			saveUser(user)
			return
		}
	}
}

/**
Checks whether the passed credential strings match the user's credentials and returns a boolean representing the result.
 */
func compareCredentials(user models.User, username, password string) bool {
	if user.UserName == username {
		return util.HashPassword(password, user.Id) == user.Password
	}
	return false
}

/**
Prompts the user to create an initial account by entering a username and password string.
Then returns the corresponding struct.
 */
func createInitialUser(reader util.Reader) models.User {
	fmt.Println("Please create an account since currently none exists.")
	fmt.Printf("- username: at least %d chars.\n", config.MIN_USERNAME_LENGTH)
	fmt.Printf("- password: at least %d chars.\n", config.MIN_PASSWORD_LENGTH)
	username := util.ReadUsername(reader, config.MIN_USERNAME_LENGTH)
	id := util.CreateHashId(username)
	password := util.ReadPassword(reader, id, config.MIN_PASSWORD_LENGTH)
	return models.User{UserName: username, Password: password, Id: id, Admin: true}
}

/**
Finds an user in a slice of multiple by comparing ids and replaces the user if it is found.
Then returns the updated users slice.
 */
func updateUsers(users []models.User, user models.User) []models.User {
	for index, record := range users {
		if record.Id == user.Id {
			users[index] = user
			break
		}
	}
	return users
}

/**
Loads all users. If none exist so far a new slice with the passed user is created and saved.
Otherwise the user is looked for in the existing slice and replaced by the new instance if found.
 */
func saveUser(user models.User) {
	users := GetUsers()
	if users == nil || len(users) == 0 {
		users = []models.User{user}
	} else {
		users = updateUsers(users, user)
	}
	saveUsersJson(users)
}

/**
Creates and saves a new account by parsing the POST form of an http(s) request.
If everything went well the username of the created account is returned.
Otherwise an error message is returned that is determined to be displayed in the front end.
 */
func CreateUser(r *http.Request) (user string, err string) {
	_, loggedIn := CheckAuthentication(r)
	if err := r.ParseForm(); err == nil && loggedIn {
		name, password, passwordConfirmation := r.FormValue("name"), r.FormValue("password"), r.FormValue("password-confirmation")
		users := GetUsers()
		for _, record := range users {
			if record.UserName == name {
				return "", "Username already exists.\n"
			}
		}
		if password != passwordConfirmation {
			return "", "Passwords don't match.\n"
		}
		if utf8.RuneCountInString(name) < config.MIN_USERNAME_LENGTH || utf8.RuneCountInString(password) < config.MIN_PASSWORD_LENGTH {
			return "", fmt.Sprintf("Username must have at least %v chars.\nPassword must have at least %v chars.\n", config.MIN_USERNAME_LENGTH, config.MIN_PASSWORD_LENGTH)
		}
		admin := false
		if r.FormValue("admin") == "on" {
			admin = true
		}
		id := util.CreateHashId(name)
		password = util.HashPassword(password, id)
		user := models.User{UserName: name, Id: id, Password: password, Admin: admin}
		users = append(users, user)
		saveUsersJson(users)
		return name, ""
	} else { // form parse error or session timeout while sending the request
		return "", "Something went wrong."
	}
}

/**
Changes the password of the currently authenticated account by parsing the POST form of an http(s) request.
Returns a string that is determined to be displayed in the frontend and that represents requirements which are not met.
If the string is empty everything went well otherwise it contains an appropriate error message.
 */
func ChangePassword(r *http.Request) string {
	user, loggedIn := CheckAuthentication(r)
	if err := r.ParseForm(); err == nil && loggedIn {
		password, passwordConfirmation := r.FormValue("password"), r.FormValue("password-confirmation")
		if password != passwordConfirmation {
			return "Passwords don't match.\n"
		}
		if utf8.RuneCountInString(password) < config.MIN_PASSWORD_LENGTH {
			return fmt.Sprintf("Password must have at least %v chars.\n", config.MIN_PASSWORD_LENGTH)
		}
		password = util.HashPassword(password, user.Id)
		users := GetUsers()
		for idx, record := range users {
			if record.Id == user.Id {
				users[idx].Password = password
			}
		}
		saveUsersJson(users)
		return ""
	} else {
		return "Something went wrong.\n"
	}
}
