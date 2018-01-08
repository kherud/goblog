package config

import "path/filepath"

var (
	DATA_PATH           = filepath.Join(".", "backend", "data")
	USERS_FILE_PATH     = filepath.Join("backend", "data", "users.json")
	ENTRIES_FILE_PATH   = filepath.Join("backend", "data", "entries.json")
	USERS_TEST_PATH     = filepath.Join("test_data", "users.json")
	ENTRIES_TEST_PATH   = filepath.Join("test_data", "entries.json")
	TEST_TEMP_PATH      = filepath.Join("test_data", "test.json")
	SESSION_TIME        = 15
	POSTS_PER_REQUESTS  = 5
	MIN_USERNAME_LENGTH = 6
	MIN_PASSWORD_LENGTH = 8
	STATIC_FILE_PATH    = filepath.Join("webserver", "static")
	TEMPLATE_PATH       = filepath.Join("webserver", "templates")
	CERT_FILE           = filepath.Join("server.crt")
	KEY_FILE            = filepath.Join("server.key")
	TEST_TEMPLATE_PATH  = "templates"
	DEFAULT_PORT        = "8080"
)
