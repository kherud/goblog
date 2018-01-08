package backend

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"os"
	"github.com/kherud/goblog/config"
	"github.com/kherud/goblog/backend/models"
)

/**
Reads and returns all users from the users.json file.
If none are found an empty slice is returned.
 */
func GetUsers() []models.User {
	raw := readFile(config.USERS_FILE_PATH)
	var entries []models.User
	json.Unmarshal(raw, &entries)
	return entries
}

/**
Reads and returns all entries from the entries.json file.
If none are found an empty slice is returned.
 */
func GetEntries() []models.Entry {
	raw := readFile(config.ENTRIES_FILE_PATH)
	var entries []models.Entry
	json.Unmarshal(raw, &entries)
	return entries
}

/**
Writes an users slice to the users.json file.
 */
func saveUsersJson(users []models.User) {
	file, err := os.Create(config.USERS_FILE_PATH)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(users); err != nil {
		panic(err)
	}
}

/**
Writes an entries slice to the entries.json file.
 */
func saveEntriesJson(entries []models.Entry) {
	file, err := os.Create(config.ENTRIES_FILE_PATH)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	encode := json.NewEncoder(file)
	if err := encode.Encode(entries); err != nil {
		panic(err)
	}
}

/**
Tries to read a file in the given path and return its raw content.
 */
func readFile(path string) []byte {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("This may not be a problem if no entries exist so far.")
		return nil
	}
	return raw
}
