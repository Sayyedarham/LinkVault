package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/hphp/linkvault/internal/scraper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hphp/linkvault/internal/model"
)

type BookmarkHandler struct {
	db        *dynamodb.Client
	tableName string
}

func NewBookmarkHandler(db *dynamodb.Client, tableName string) *BookmarkHandler {
	return &BookmarkHandler{db: db, tableName: tableName}
}

func (h *BookmarkHandler) Create(c *gin.Context) {
	var req model.CreateBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Title == "" {
		req.Title = scraper.FetchTitle(c.Request.Context(), req.URL)
	}

	userID := c.GetString("user_id") // From JWT middleware (stub for now)
	if userID == "" {
		userID = "demo-user" // TEMP: remove when JWT ready
	}

	bookmarkID := uuid.New().String()
	now := time.Now()

	bookmark := model.Bookmark{
		PK:        "USER#" + userID,
		SK:        "BOOKMARK#" + bookmarkID,
		ID:        bookmarkID,
		UserID:    userID,
		URL:       req.URL,
		Title:     req.Title,
		Tags:      req.Tags,
		CreatedAt: now,
		UpdatedAt: now,
	}

	item, err := attributevalue.MarshalMap(bookmark)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "marshal failed"})
		return
	}

	_, err = h.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(h.tableName),
		Item:      item,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bookmark)
}

func (h *BookmarkHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "demo-user"
	}

	result, err := h.db.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(h.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "USER#" + userID},
			":sk": &types.AttributeValueMemberS{Value: "BOOKMARK#"},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var bookmarks []model.Bookmark
	err = attributevalue.UnmarshalListOfMaps(result.Items, &bookmarks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unmarshal failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bookmarks": bookmarks, "count": len(bookmarks)})
}

func (h *BookmarkHandler) Delete(c *gin.Context) {
	bookmarkID := c.Param("id")
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "demo-user"
	}

	_, err := h.db.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(h.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "USER#" + userID},
			"SK": &types.AttributeValueMemberS{Value: "BOOKMARK#" + bookmarkID},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}