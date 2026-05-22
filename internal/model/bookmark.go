package model

import "time"

type Bookmark struct {
	PK        string    `json:"-" dynamodbav:"PK"`           // USER#<userID>
	SK        string    `json:"-" dynamodbav:"SK"`           // BOOKMARK#<bookmarkID>
	ID        string    `json:"id" dynamodbav:"id"`
	UserID    string    `json:"user_id" dynamodbav:"user_id"`
	URL       string    `json:"url" dynamodbav:"url" binding:"required,url"`
	Title     string    `json:"title" dynamodbav:"title"`
	Tags      []string  `json:"tags" dynamodbav:"tags"`
	CreatedAt time.Time `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt time.Time `json:"updated_at" dynamodbav:"updated_at"`
}

type CreateBookmarkRequest struct {
	URL   string   `json:"url" binding:"required,url"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}