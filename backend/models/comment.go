package models

type Comment struct {
	Text     string `json:"text"`
	Author   string `json:"author"`
	Date     string `json:"date"`
	Verified bool   `json:"verified"`
	Id       uint32 `json:"id"`
}

