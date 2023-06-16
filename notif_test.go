package notif

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var keyID string = ""
var teamID string = ""
var fileName string = ""

var user string = ""
var password string = ""
var addr string = ""
var dbName string = ""

func TestCreateClient(t *testing.T) {
	l := logrus.New()

	// database connection
	poolStr := "postgres://" + user + ":" + password + "@" + addr + "/" + dbName + "?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), poolStr)
	if err != nil {
		t.Error("failed to connect to database: ", err.Error())
		return
	}

	// create notif service and client
	n := NewNotificationService(l, pool)
	err = n.CreateNewClient("JN25FUC9X2", "5A6H49Q85D", "./AuthKey_JN25FUC9X2.p8")
	if err != nil {
		t.Error("failed to create notif client: ", err.Error())
	}
}

func TestSendNotificationToToken(t *testing.T) {
	l := logrus.New()

	// database connection
	poolStr := "postgres://" + user + ":" + password + "@" + addr + "/" + dbName + "?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), poolStr)
	if err != nil {
		t.Error("failed to connect to database: ", err.Error())
		return
	}

	// create notif service and client
	n := NewNotificationService(l, pool)
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

	// database connection
	poolStr := "postgres://" + user + ":" + password + "@" + addr + "/" + dbName + "?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), poolStr)
	if err != nil {
		t.Error("failed to connect to database: ", err.Error())
		return
	}

	// create notif service and client
	n := NewNotificationService(l, pool)
	err = n.CreateNewClient(keyID, teamID, fileName)
	if err != nil {
		t.Error("failed to create notif client: ", err.Error())
	}

	// test notification to device
	note := Notification{
		Title: "Test Notification",
		Body:  "This is a notification test",
		Topic: "647f712d40f9d342d309496f_admin",
	}

	err = n.SendNotificationToTopic(&note)
	if err != nil {
		t.Error("failed to send token: ", err.Error())
	}
}

func TestUpdateUserToken(t *testing.T) {
	l := logrus.New()

	// database connection
	poolStr := "postgres://" + user + ":" + password + "@" + addr + "/" + dbName + "?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), poolStr)
	if err != nil {
		t.Error("failed to connect to database: ", err.Error())
		return
	}

	// create notif service and client
	n := NewNotificationService(l, pool)
	err = n.UpdateUserToken("e57673a3-9bd5-4f80-a8b4-4a0fab316b00", []string{"647f712d40f9d342d309496f_admin", "6488a0fcc534706f8c817b76", "647f712d40f9d342d309496f", "647f6e7e40f9d342d3094969_admin", "647f6e7e40f9d342d3094969"}, "8f9348fc8ce61585b9ba805bd960dbeefac1015eb71dbefad50d3c94e4588adf")
	if err != nil {
		t.Error("failed to connect to update topics: ", err.Error())
		return
	}
}
