package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"chatz.com/api/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	mongoClient *mongo.Client
)

// Get the messages related to a certain author in a certain chat room, since forever.
func getMessage(c *gin.Context) {

	reqParams := c.Request.URL.Query()

	// reqParams keys are associated with an array of strings
	// because in a request, regarding query paranms, there can exist multiple equal keys
	// so all the values would be stored in an array of values

	room_id, _ := strconv.Atoi(reqParams["room_id"][0])
	author_id, _ := strconv.Atoi(reqParams["author_id"][0])

	messagesColl := mongoClient.Database("chats").Collection("messages")

	var result bson.M
	err := messagesColl.FindOne(c, bson.D{
		{Key: "author_id", Value: author_id},
		{Key: "room_id", Value: room_id},
	}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		fmt.Printf("No messages were found with author_id: %d or with room_id: %d\n", author_id, room_id)
		return
	}

	c.JSON(http.StatusOK, result)
}

func getUserValueFromName(c *gin.Context) {

	/*mongoClient := config.InitClient(context.TODO())
	defer config.DisconnectClient()

	username := c.Params.ByName("name")

	// Comprobar si existe un value para dicho username en mongoDB

	if ok {
		c.JSON(http.StatusOK, gin.H{"user": username, "value": value})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"user": username, "status": "no value"})
	}*/
}

func setupRouter() *gin.Engine {

	r := gin.Default()

	r.GET("/msg", getMessage)
	r.GET("/user/:name", getUserValueFromName)

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // foo:bar
		"manu": "123", // manu:123
	}))

	/* example basicauth header:
	authorization: Basic Zm9vOmJhcg==
	Body: {"value":"bar"}
	*/
	authorized.POST("admin", func(c *gin.Context) {
		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {

	cwd, _ := os.Getwd()

	err := godotenv.Load(filepath.Join(cwd, "./config/.env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoClient = config.InitClient(context.TODO())
	defer config.DisconnectClient()

	r := setupRouter()
	r.Run(":8080")
}
