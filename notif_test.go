package notif

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var keyID string = ""
var teamID string = ""
var fileName string = ""

var user string = ""
var password string = ""
var addr string = ""
var dbName string = ""

var DB_USR = ""
var DB_PASS = ""
var DB_NAME = ""

func TestCreateClient(t *testing.T) {
	l := logrus.New()

	// sql database connection
	poolStr := "postgres://" + user + ":" + password + "@" + addr + "/" + dbName + "?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), poolStr)
	if err != nil {
		t.Error("failed to connect to database: ", err.Error())
		return
	}

	// document database connection
	credential := options.Credential{
		AuthSource: "admin",
		Username:   DB_USR,
		Password:   DB_PASS,
	}
	dbLoc := os.Getenv("DB_ADDR")
	opts := options.Client().ApplyURI("mongodb://" + dbLoc + ":27017")
	opts.Auth = &credential

	client, _ := mongo.Connect(context.Background(), opts)
	UserCol := client.Database(DB_NAME).Collection("users")

	// create notif service and client
	n := NewNotificationService(l, pool, UserCol)
	err = n.CreateNewClient("JN25FUC9X2", "5A6H49Q85D", "./AuthKey_JN25FUC9X2.p8")
	if err != nil {
		t.Error("failed to create notif client: ", err.Error())
	}
}

func TestSendNotificationToToken(t *testing.T) {
	l := logrus.New()

	// sql database connection
	poolStr := "postgres://" + user + ":" + password + "@" + addr + "/" + dbName + "?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), poolStr)
	if err != nil {
		t.Error("failed to connect to database: ", err.Error())
		return
	}

	// document database connection
	credential := options.Credential{
		AuthSource: "admin",
		Username:   DB_USR,
		Password:   DB_PASS,
	}
	dbLoc := os.Getenv("DB_ADDR")
	opts := options.Client().ApplyURI("mongodb://" + dbLoc + ":27017")
	opts.Auth = &credential

	client, _ := mongo.Connect(context.Background(), opts)
	UserCol := client.Database(DB_NAME).Collection("users")

	// create notif service and client
	n := NewNotificationService(l, pool, UserCol)
	err = n.CreateNewClient("JN25FUC9X2", "5A6H49Q85D", "./AuthKey_JN25FUC9X2.p8")
	if err != nil {
		t.Error("failed to create notif client: ", err.Error())
	}

	// test notification to device
	note := Notification{
		Title: "Test Notification",
		Body:  "This is a notification test",
	}

	// send notification to token
	err = n.SendNotificationToToken(&note, "8f9348fc8ce61585b9ba805bd960dbeefac1015eb71dbefad50d3c94e4588adf")
	if err != nil {
		t.Error("failed to send token: ", err.Error())
	}
}

func TestSendNotificationToTopic(t *testing.T) {
	l := logrus.New()

	// sql database connection
	poolStr := "postgres://" + user + ":" + password + "@" + addr + "/" + dbName + "?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), poolStr)
	if err != nil {
		t.Error("failed to connect to database: ", err.Error())
		return
	}

	// document database connection
	credential := options.Credential{
		AuthSource: "admin",
		Username:   DB_USR,
		Password:   DB_PASS,
	}
	opts := options.Client().ApplyURI("mongodb://" + addr + ":27017")
	opts.Auth = &credential

	client, _ := mongo.Connect(context.Background(), opts)
	UserCol := client.Database(DB_NAME).Collection("users")

	// create notif service and client
	n := NewNotificationService(l, pool, UserCol)
	err = n.CreateNewClient(keyID, teamID, fileName)
	if err != nil {
		t.Error("failed to create notif client: ", err.Error())
	}

	// test notification to device
	note := Notification{
		Title: "Test Notification",
		Body:  "This is a notification test",
		Topic: "64945010d69bf7890c997866",
	}

	err = n.SendNotificationToTopic(&note)
	if err != nil {
		t.Error("failed to send token: ", err.Error())
	}
}
