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

type Message struct {
	Author_id int    `json:"author_id"`
	Room_id   int    `json:"room_id"`
	Message   string `json:"message"`
}

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

func getMessagesByRoomId(c *gin.Context) {

	reqParams := c.Request.URL.Query()
	room_id, _ := strconv.Atoi(reqParams["room_id"][0])

	messagesColl := mongoClient.Database("chats").Collection("messages")

	var result bson.M
	err := messagesColl.FindOne(c, bson.D{
		{Key: "room_id", Value: room_id},
	}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		fmt.Printf("No messages were found for the specified room: %d\n", room_id)
		c.String(http.StatusNotFound, "No messages were found for the specified room: "+strconv.Itoa(room_id))
		return
	}

	c.JSON(http.StatusOK, result)

}

func getRoomById(c *gin.Context) {

	reqParams := c.Request.URL.Query()
	room_id, _ := strconv.Atoi(reqParams["room_id"][0])

	roomsColl := mongoClient.Database("chats").Collection("rooms")

	var result bson.M
	err := roomsColl.FindOne(c, bson.D{
		{Key: "room_id", Value: room_id},
	}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		fmt.Printf("No room were found for the specified id: %d\n", room_id)
		c.String(http.StatusNotFound, "No rooms were found for the specified id: "+strconv.Itoa(room_id))
		return
	}

	c.JSON(http.StatusOK, result)

}

func sendMessage(c *gin.Context) {

	messagesColl := mongoClient.Database("chats").Collection("messages")

	var payload []Message
	docs := make([]interface{}, len(payload))

	for _, element := range payload {
		docs = append(docs, element)
	}

	res, err := messagesColl.InsertMany(c, docs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "There was an error, and documents could not be created.")
	}

	c.JSON(http.StatusOK, result)

}

func setupRouter() *gin.Engine {

	r := gin.Default()

	r.GET("/user/:name", getUserValueFromName)

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"root": "root",
	}))

	authorized.GET("/msg/:room_id", getMessagesByRoomId)
	authorized.GET("/room/:room_id", getRoomById)
	authorized.POST("/msg", sendMessage)

	authorized.POST("msg", func(c *gin.Context) {
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
