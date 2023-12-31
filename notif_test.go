package notif

import (
	"context"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var KEYID string = ""
var TEAMID string = ""
var FILENAME string = ""

var SERVER_USR = ""
var SERVER_PASS = ""
var SERVER_LOC = ""
var SERVER_DB = "olympsis"
var USERS_COL = "users"

var NOTIF_USR = ""
var NOTIF_PASS = ""
var NOTIF_LOC = ""
var NOTIF_DB = "notifications"
var TOPICS_COL = "topics"

var TEST_UUID = ""
var TEST_TOKEN = ""

func TestCreateClient(t *testing.T) {
	l := logrus.New()
	credential := options.Credential{
		AuthSource: "admin",
		Username:   SERVER_USR,
		Password:   SERVER_PASS,
	}
	opts := options.Client().ApplyURI("mongodb://" + SERVER_LOC + ":27017")
	opts.Auth = &credential
	client, _ := mongo.Connect(context.Background(), opts)
	UserCol := client.Database(SERVER_DB).Collection(USERS_COL)

	opts = options.Client().ApplyURI(`mongodb+srv://admin:` + NOTIF_PASS + `@database-2.pdjjqal.mongodb.net/?retryWrites=true&w=majority`)
	client2, _ := mongo.Connect(context.Background(), opts)
	Database := client2.Database(NOTIF_DB).Collection(TOPICS_COL)

	// create notif service and client
	n := NewNotificationService(l, Database, UserCol)
	err := n.CreateNewClient(KEYID, TEAMID, FILENAME)
	if err != nil {
		t.Error("failed to create notif client: ", err.Error())
	}
}

func TestHappyPath(t *testing.T) {
	l := logrus.New()
	credential := options.Credential{
		AuthSource: "admin",
		Username:   SERVER_USR,
		Password:   SERVER_PASS,
	}
	opts := options.Client().ApplyURI("mongodb://" + SERVER_LOC + ":27017")
	opts.Auth = &credential
	client, _ := mongo.Connect(context.Background(), opts)
	UserCol := client.Database(SERVER_DB).Collection(USERS_COL)

	opts2 := options.Client().ApplyURI(`mongodb+srv://admin:` + NOTIF_PASS + `@` + NOTIF_LOC + `/?retryWrites=true&w=majority`)
	client2, _ := mongo.Connect(context.Background(), opts2)
	Database := client2.Database(NOTIF_DB).Collection(TOPICS_COL)

	// create notif service and client
	n := NewNotificationService(l, Database, UserCol)
	err := n.CreateNewClient(KEYID, TEAMID, FILENAME)
	if err != nil {
		t.Error("failed to create notif client: ", err.Error())
	}

	// test create new topic
	id := primitive.NewObjectID()
	err = n.CreateTopic(id.Hex())
	if err != nil {
		t.Error("failed to create topic: ", err.Error())
		return
	}

	// test adding token to topic
	err = n.AddTokenToTopic(id.Hex(), TEST_UUID)
	if err != nil {
		t.Error("failed to add token to topic: ", err.Error())
		return
	}

	// test notification to topic
	note := Notification{
		Title: "Test Notification",
		Body:  "Happy Path Topic Test",
		Topic: id.Hex(),
	}
	err = n.SendNotificationToTopic(&note)
	if err != nil {
		t.Error("failed to send notification to topic: ", err.Error())
		return
	}

	// test remove token from topic
	err = n.RemoveTokenFromTopic(id.Hex(), TEST_UUID)
	if err != nil {
		t.Error("failed to remove token from topic: ", err.Error())
	}

	// test deleting topic
	err = n.DeleteTopic(id.Hex())
	if err != nil {
		t.Error("failed to delete topic: ", err.Error())
	}

	note = Notification{
		Title: "Test Notification",
		Body:  "Happy Path Token Test",
	}
	err = n.SendNotificationToToken(&note, TEST_TOKEN)
	if err != nil {
		t.Error("failed to send notification to token: ", err.Error())
	}

	users, err := n.FindUsers([]string{TEST_UUID})
	if err != nil {
		t.Error("failed to find users: ", err.Error())
	}
	for i := range users {
		fmt.Print(users[i].UUID)
	}
}
