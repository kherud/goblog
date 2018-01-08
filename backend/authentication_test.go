package backend

import (
	"testing"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"time"
	"strings"
	"unicode/utf8"
	"net/http"
)

func TestSetSession(t *testing.T){
	recorder := httptest.NewRecorder()
	SetSession("Konstant", recorder)
	cookie, exists := recorder.HeaderMap["Set-Cookie"]
	assert.True(t, exists)
	parts := strings.Split(cookie[0], ";")
	assert.True(t, len(parts) == 2)
	sessionParts := strings.Split(parts[0], "=")
	sessionParts = strings.Split(sessionParts[1], "#")
	assert.EqualValues(t, sessionParts[0], "Konstant")
	assert.True(t, utf8.RuneCountInString(sessionParts[1]) == 128)
	dateLayout := "Mon, 02 Jan 2006 15:04:05 MST"
	expirationParts := strings.Split(parts[1], "=")
	expiration, err := time.Parse(dateLayout, expirationParts[1])
	assert.Nil(t, err)
	before := time.Now().Add(time.Minute * 14).Add(time.Second * 70)
	after := time.Now().Add(time.Minute * 14).Add(time.Second * 50)
	assert.True(t, expiration.After(after))
	assert.True(t, expiration.Before(before))
	user, err := GetUser("Konstant")
	assert.Nil(t, err)
	assert.EqualValues(t, user.Session, sessionParts[1])
}

func TestEndSessionValidUser(t *testing.T){
	recorder := httptest.NewRecorder()
	SetSession("Konstant", recorder)
	user, err := GetUser("Konstant")
	assert.Nil(t, err)
	assert.NotEmpty(t, user.Session)
	cookieString, exists := recorder.HeaderMap["Set-Cookie"]
	assert.True(t, exists)
	parts := strings.Split(cookieString[0], ";")
	sessionParts := strings.Split(parts[0], "=")
	cookie := &http.Cookie{Name: sessionParts[0], Value: sessionParts[1]}
	req := &http.Request{
		Header: http.Header{},
	}
	req.AddCookie(cookie)
	EndSession(req)
	user, err = GetUser("Konstant")
	assert.Nil(t, err)
	assert.Empty(t, user.Session)
}

func TestEndSessionInvalidUser(t *testing.T){
	user, err := GetUser("Test")
	assert.NotNil(t, err)
	assert.Empty(t, user.Session)
	cookie := &http.Cookie{Name: "Session", Value: "Test#Test"}
	req := &http.Request{
		Header: http.Header{},
	}
	req.AddCookie(cookie)
	EndSession(req)
	user, err = GetUser("Test")
	assert.NotNil(t, err)
	assert.Empty(t, user.Session)
	assert.True(t, user.Id == 0)
}

func TestEndSessionInvalidCookie(t *testing.T){
	user, err := GetUser("Konstantin")
	assert.Nil(t, err)
	assert.NotEmpty(t, user.Session)
	req := &http.Request{}
	EndSession(req)
	user, err = GetUser("Konstantin")
	assert.Nil(t, err)
	assert.NotEmpty(t, user.Session)
}

func TestCheckAuthenticationInvalid(t *testing.T) {
	tests := []struct {name string; value string}{
		{"", ""},
		{"", "Konstantin#Test"},
		{"Session", "Test#Test"},
	}
	for _, test := range tests {
		req := &http.Request{
			Header: http.Header{},
		}
		cookie := &http.Cookie{Name: test.name, Value: test.value}
		req.AddCookie(cookie)
		user, authenticated := CheckAuthentication(req)
		assert.False(t, authenticated)
		assert.EqualValues(t, user.Id, uint32(0))
	}
}

func TestCheckAuthenticationValid(t *testing.T) {
	tests := []struct {name string; value string; exptectedUsername string; expectedUserId uint32; expectedAuthentication bool}{
		{"Session", "Konstantin#Test", "Konstantin", 689017489,true},
		{"Session", "Konstanti#Test", "Konstanti", 3876830309, false},
	}
	for _, test := range tests {
		req := &http.Request{
			Header: http.Header{},
		}
		cookie := &http.Cookie{Name: test.name, Value: test.value}
		req.AddCookie(cookie)
		user, authenticated := CheckAuthentication(req)
		assert.EqualValues(t, authenticated, test.expectedAuthentication)
		assert.EqualValues(t, user.Id, test.expectedUserId)
		assert.EqualValues(t, user.UserName, test.exptectedUsername)
	}
}

func TestProcessCookie(t *testing.T){
	tests := []struct {test string; expectedA string; expectedB string}{
		{"Konstantin#12345678", "Konstantin", "12345678"},
		{"#12345678", "", "12345678"},
		{"Konstantin#", "Konstantin", ""},
	}
	for _, test := range tests {
		resultA, resultB := processCookie(test.test)
		assert.EqualValues(t, resultA, test.expectedA)
		assert.EqualValues(t, resultB, test.expectedB)
	}
}
