package util

import (
	"fmt"
	"strings"
	"unicode/utf8"
	"hash/fnv"
	"time"
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strconv"
)

/**
Using a wrapper for proper testing
 */
type Reader interface {
	ReadString(byte) (string, error)
}

/**
Repeats to ask for a username until a string of proper length (param: minimumLength) is entered.
 */
func ReadUsername(reader Reader, minimumLength int) string {
	var username string
	for utf8.RuneCountInString(username) < minimumLength {
		fmt.Print("Enter a valid username: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		text = strings.TrimSpace(text)
		username = text
	}
	return username
}

/**
Repeats to ask for a password until a string of proper length (param: minimumLength) is entered.
 */
func ReadPassword(reader Reader, id uint32, minimumLength int) string {
	var password string
	for utf8.RuneCountInString(password) < minimumLength {
		fmt.Print("Enter a valid password: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		text = strings.TrimSpace(text)
		password = text
	}
	return HashPassword(password, id)
}

/**
Takes the plain password and an id for salting which is put in the middle of the plain password.
The resulting string then is hashed with SHA256 and returned as base64 encoding.
 */
func HashPassword(password string, id uint32) string {
	shaHash := sha256.New()
	saltIndex := utf8.RuneCountInString(password) / 2
	saltedPassword := password[:saltIndex] + strconv.Itoa(int(id)) + password[saltIndex:]
	shaHash.Write([]byte(saltedPassword))
	return base64.URLEncoding.EncodeToString(shaHash.Sum(nil))
}

/**
Takes a list of arguments that are used in combination with the current time stamp to create an unique string.
The result then is hashed with FNV.
 */
func CreateHashId(prehash ... string) uint32 {
	preHash := strconv.FormatInt(time.Now().UnixNano(), 10)
	for _, element := range prehash {
		preHash += element
	}
	hash32 := fnv.New32()
	hash32.Write([]byte(preHash))
	return hash32.Sum32()
}

/**
Creates an 128 char string of randomly chosen defined chars that is used for identifying sessions.
 */
func CreateSessionId() string {
	// code used from https://stackoverflow.com/a/31832326/4812335
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	)

	b := make([]byte, 128)
	for i := 0; i < 128; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return string(b)
}
