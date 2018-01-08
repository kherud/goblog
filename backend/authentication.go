package backend

import (
	"net/http"
	"time"
	"strings"
	"github.com/kherud/goblog/config"
	"github.com/kherud/goblog/util"
	"github.com/kherud/goblog/backend/models"
)

/**
Creates a cookie and saves it for a username.
 */
func SetSession(username string, w http.ResponseWriter) {
	expiration := time.Now().Add(time.Minute * time.Duration(config.SESSION_TIME))	//expiration time for cookies
	sessionId := util.CreateSessionId()	//creation of a sessionID
	persistSession(username, sessionId)
	cookie := http.Cookie{Name: "Session", Value: username + "#" + sessionId, Expires: expiration}	//specifies cookie informations
	http.SetCookie(w, &cookie)
}

/**
Searches for an active session and ends it.
 */
func EndSession(r *http.Request){
	cookie, err := r.Cookie("Session")
	if err != nil {
		return
	}

	username, _ := processCookie(cookie.Value)
	user, err := GetUser(username)

	if err != nil {
		return
	}

	user.Session = ""
	saveUser(user)
}

/**
Check if a cookie exists and if its session id is valid.
 */
func CheckAuthentication(r *http.Request) (models.User, bool) {
	// Tutorial used from http://austingwalters.com/building-a-web-server-in-go-web-cookies/
	cookie, err := r.Cookie("Session")
	if err != nil {
		return models.User{}, false
	}

	username, sessionId := processCookie(cookie.Value)
	user, err := GetUser(username)
	if err != nil {
		return models.User{}, false
	}
	return user, sessionId == user.Session && len(user.Session) > 0
}

/**
Splits the value of a session cookie into username and the correspondent id.
 */
func processCookie(cookieValue string) (string, string) {
	credentials := strings.Split(cookieValue, "#")
	username := credentials[0]
	sessionId := credentials[1]
	return username, sessionId
}

