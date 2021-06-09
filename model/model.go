package model

type InBoxResp struct {
	Posts []Post `json:"posts"`
	Users []User `json:"users"`
}

type Post struct {
	ID int64 `json:"id"`
	Content string `json:"content"`
	UserID int64 `json:"user_id"`
	CreatedAt int64 `json:"created_at"`
}

type User struct {
	ID int64 `json:"id"`
	Username string `json:"username"`
	DisplayName string `json:"display_name"`
}
