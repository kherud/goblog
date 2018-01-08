package webserver

import (
	"testing"
	"net/http"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"os"
	"path/filepath"
	"crypto/tls"
	"net/url"
	"net/http/cookiejar"
	"github.com/kherud/goblog/config"
)

func TestMain(m *testing.M) {
	config.TEMPLATE_PATH = config.TEST_TEMPLATE_PATH
	config.ENTRIES_FILE_PATH = filepath.Join("..", "backend", "test_data", "entries.json")
	config.USERS_FILE_PATH = filepath.Join("..", "backend", "test_data", "users.json")
	os.Exit(m.Run())
}

// Backend logic is separately unit tested, thus only the delivery of proper content is tested
// Requests that require authentication are tested twice

func TestReturnContentId(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080", false)
	assert.True(t, strings.Contains(string(body), "Recent posts"))
}

func TestReturnContentComment(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?comment", false)
	assert.True(t, strings.Contains(string(body), "404: Post not found."))
}

func TestReturnContentPost(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?post", true)
	assert.True(t, strings.Contains(string(body), "Create an entry..."))
}

func TestReturnContentPostInvalid(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?post", false)
	assert.True(t, strings.Contains(string(body), "Recent posts"))
}

func TestReturnContentUser(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?account", true)
	assert.True(t, strings.Contains(string(body), "Change password..."))
	assert.True(t, strings.Contains(string(body), "Create an user..."))
}

func TestReturnContentUserInvalid(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?account", false)
	assert.True(t, strings.Contains(string(body), "Recent posts"))
}

func TestReturnContentSearch(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?search=asd", false)
	assert.True(t, strings.Contains(string(body), "Search results for '"))
	assert.True(t, strings.Contains(string(body), "Post #60"))
	assert.True(t, strings.Contains(string(body), "Hi friend"))
}

func TestReturnContentMore(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?more=0", false)
	assert.True(t, strings.Contains(string(body), "Post #41"))
	assert.True(t, strings.Contains(string(body), "Post #60"))
	assert.True(t, strings.Contains(string(body), "Post #10"))
	assert.True(t, strings.Contains(string(body), "Hi friend"))
	assert.True(t, strings.Contains(string(body), "Post #5"))
}

func TestReturnContentMoreEnd(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?more=5", false)
	assert.True(t, strings.Contains(string(body), "Post #4"))
	assert.True(t, strings.Contains(string(body), "Post #3"))
}

func TestReturnContentNewPost(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?newPost", true)
	assert.True(t, strings.Contains(string(body), "Create an entry..."))
}

func TestReturnContentNewPostInvalid(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?newPost", false)
	assert.True(t, strings.Contains(string(body), "Recent posts"))
}

func TestReturnContentNewUser(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?newUser", true)
	assert.True(t, strings.Contains(string(body), "#Username must have at least 6 chars."))
}

func TestReturnContentNewUserInvalid(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?newUser", false)
	assert.True(t, strings.Contains(string(body), "#Something went wrong."))
}

func TestReturnContentDelete(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?delete=0", true)
	assert.True(t, strings.Contains(string(body), "false"))
}

func TestReturnContentDeleteInvalid(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?delete=0", false)
	assert.True(t, strings.Contains(string(body), "false"))
}

func TestReturnContentEdit(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?edit=0", true)
	assert.True(t, strings.Contains(string(body), "404: Post not found."))
}

func TestReturnContentEditInvalid(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?edit=0", false)
	assert.True(t, strings.Contains(string(body), "Recent posts"))
}

func TestReturnContentUpdate(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?update=0", true)
	assert.True(t, strings.Contains(string(body), "404: Post not found."))
}

func TestReturnContentUpdateInvalid(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?update=0", false)
	assert.True(t, strings.Contains(string(body), "404: Post not found."))
}

func TestReturnContentPassword(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?password", true)
	assert.True(t, strings.Contains(string(body), "Password must have at least 8 chars."))
}

func TestReturnContentPasswordInvalid(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?password", false)
	assert.True(t, strings.Contains(string(body), "Something went wrong."))
}

func TestReturnContentVerify(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?verify", true)
	assert.True(t, strings.Contains(string(body), "false"))
}

func TestReturnContentVerifyInvalid(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?verify", false)
	assert.True(t, strings.Contains(string(body), "false"))
}

func TestReturnContentVerifyDefault(t *testing.T) {
	body := testServerRequest(t, "https://localhost:8080?Test", false)
	assert.True(t, strings.Contains(string(body), "Recent posts"))
}

func TestLoginUserValid(t *testing.T) {
	srv, cli := getHTTPSServerClient(false)
	defer srv.Close()
	srv.Handler = http.HandlerFunc(loginUser)
	form := url.Values{}
	form.Add("username", "Konstant")
	form.Add("password", "12345678")
	res, err := cli.PostForm("https://localhost:8080", form)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.EqualValues(t, body, "success")
}

func TestLoginUserInvalid(t *testing.T) {
	srv, cli := getHTTPSServerClient(false)
	defer srv.Close()
	srv.Handler = http.HandlerFunc(loginUser)
	form := url.Values{}
	form.Add("username", "Konstant")
	form.Add("password", "")
	res, err := cli.PostForm("https://localhost:8080", form)
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.EqualValues(t, body, "failed to login")
}

// The actual logout process is separately tested
// No login to keep session information for future tests
func TestLogoutUser(t *testing.T) {
	srv, cli := getHTTPSServerClient(false)
	defer srv.Close()
	// Necessary to not cyclically redirect to /logout
	srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/logout" {
			logoutUser(w, r)
		}
		returnContent(w, r)
	})
	req, err := http.NewRequest("POST", "https://localhost:8080/logout", nil)
	assert.NoError(t, err)
	res, err := cli.Do(req)
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(body), "Recent posts"))
}

func testServerRequest(t *testing.T, url string, login bool) []byte {
	srv, client := getHTTPSServerClient(login)
	defer srv.Close()
	res, err := client.Get(url)
	assert.NoError(t, err)
	assert.EqualValues(t, res.StatusCode, http.StatusOK)
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.True(t, len(body) > 0)
	return body
}

func getHTTPSServerClient(login bool) (srv *http.Server, client *http.Client) {
	srv = &http.Server{
		Addr:    ":" + config.DEFAULT_PORT,
		Handler: http.HandlerFunc(returnContent),
	}
	go srv.ListenAndServeTLS("../server.crt", "../server.key")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}
	if login {
		testUrl, _ := url.Parse("https://localhost:" + config.DEFAULT_PORT)
		cookie := &http.Cookie{Name: "Session", Value: "Konstantin#Test"}
		cookies := []*http.Cookie{cookie}
		newcookiejar, _ := cookiejar.New(nil)
		newcookiejar.SetCookies(testUrl, cookies)
		client.Jar = newcookiejar
	}
	return srv, client
}
