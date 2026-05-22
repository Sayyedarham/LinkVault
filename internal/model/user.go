package model

type User struct {
	PK           string `dynamodbav:"PK"`
	SK           string `dynamodbav:"SK"`
	Email        string `dynamodbav:"email"`
	PasswordHash string `dynamodbav:"passwordHash"`
	UserID       string `dynamodbav:"userId"`
}
