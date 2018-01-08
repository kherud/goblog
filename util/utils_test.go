package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"unicode/utf8"
	"time"
)

type testReader struct {
	text []rune
}

type testPasswordReader struct {
	text []rune
}

func (tr testReader) ReadString(byte) (string, error) {
	return string(tr.text[0:rand.Intn(len(tr.text))]), nil
}

func (tr testPasswordReader) ReadString(byte)(string, error){
	return string(tr.text), nil
}

func TestReadUsername(t *testing.T) {
	tr := testReader{text: []rune("  ß!§$%&/(äö `ü^°'123_  ")}
	for length := 0; length < 10; length++ {
		text := ReadUsername(tr, length)
		actualLength := utf8.RuneCountInString(text)
		assert.True(t, length <= actualLength && actualLength <= 20)
		assert.EqualValues(t, tr.text[2:2+actualLength], text)
	}
}

func TestReadPassword(t *testing.T) {
	tr := testPasswordReader{text: []rune("  ß!§$%&/(äö `ü^°'123_  ")}
	text := ReadPassword(tr, 689017489, 5)
	text2 := ReadPassword(tr, 976620356, 5)
	actualLength := utf8.RuneCountInString(text)
	actualLength2 := utf8.RuneCountInString(text2)
	assert.True(t, actualLength == actualLength2)
	assert.True(t, actualLength > len(tr.text))
	assert.True(t, actualLength2 > len(tr.text))
	assert.NotEqual(t, text, text2)
}

func TestHashPassword(t *testing.T){
	parameters := []string{"", "abc", "abcd", "abcabc", "abcabcabc"}
	for _, parameter := range parameters {
		hash := HashPassword(parameter, 689017489)
		hash2 := HashPassword(parameter, 976620356)
		assert.True(t, utf8.RuneCountInString(hash) > utf8.RuneCountInString(parameter))
		assert.True(t, utf8.RuneCountInString(hash2) > utf8.RuneCountInString(parameter))
		assert.NotEqual(t, hash, hash2)
	}
}

func TestCreateHashId(t *testing.T){
	parameters := []string{"a", "abc", "abcd", "abcabc", "abcabcabc", ""}
	ids := make(map[uint32]int)
	for idx := range parameters {
		ids[CreateHashId(parameters[:idx]...)] = 1
		time.Sleep(1)
		ids[CreateHashId(parameters[:idx]...)] = 1
	}
	// len(ids) = unique hashes -> must be two times of len(parameters) since every parameter list is hashed twice
	assert.True(t, len(ids) == len(parameters) * 2)
}

func TestCreateSessionId(t *testing.T) {
	ids := make(map[string]int)
	for idx := 0; idx < 10; idx++ {
		id := CreateSessionId()
		assert.True(t, utf8.RuneCountInString(id) == 128)
		ids[id] = 1
	}
	// len(ids) = unique ids -> must be 10
	assert.True(t, len(ids) == 10)
}
