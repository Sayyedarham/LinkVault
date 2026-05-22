package handler

import (
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db        *dynamodb.Client
	tableName string
	jwtSecret string
}

func NewAuthHandler(db *dynamodb.Client, tableName, jwtSecret string) *AuthHandler {
	return &AuthHandler{db: db, tableName: tableName, jwtSecret: jwtSecret}
}

type authRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check duplicate
	existing, _ := h.db.GetItem(c, &dynamodb.GetItemInput{
		TableName: &h.tableName,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "USER#" + req.Email},
			"SK": &types.AttributeValueMemberS{Value: "PROFILE"},
		},
	})
	if existing.Item != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	userID := uuid.New().String()
	_, err = h.db.PutItem(c, &dynamodb.PutItemInput{
		TableName: &h.tableName,
		Item: map[string]types.AttributeValue{
			"PK":           &types.AttributeValueMemberS{Value: "USER#" + req.Email},
			"SK":           &types.AttributeValueMemberS{Value: "PROFILE"},
			"userId":       &types.AttributeValueMemberS{Value: userID},
			"email":        &types.AttributeValueMemberS{Value: req.Email},
			"passwordHash": &types.AttributeValueMemberS{Value: string(hash)},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	token, err := h.signToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user":  gin.H{"id": userID, "email": req.Email},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.db.GetItem(c, &dynamodb.GetItemInput{
		TableName: &h.tableName,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "USER#" + req.Email},
			"SK": &types.AttributeValueMemberS{Value: "PROFILE"},
		},
	})
	if err != nil || result.Item == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	storedHash := result.Item["passwordHash"].(*types.AttributeValueMemberS).Value
	userID := result.Item["userId"].(*types.AttributeValueMemberS).Value

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := h.signToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  gin.H{"id": userID, "email": req.Email},
	})
}

func (h *AuthHandler) signToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.jwtSecret))
}