package models

type User struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Id       uint32 `json:"id"`
	Session  string `json:"session"`
	Admin    bool   `json:"admin"`
}
