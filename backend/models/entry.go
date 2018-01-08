package models

type Entry struct {
	Title    string    `json:"title"`
	Text     string    `json:"text"`
	Author   string    `json:"author"`
	AuthorId uint32    `json:"author_id"`
	Date     string    `json:"date"`
	Id       uint32    `json:"id"`
	Comments []Comment `json:"comments"`
	Keywords []string  `json:"keywords"`
}
