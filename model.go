package notif

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sideshow/apns2"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Client   *apns2.Client
	Logger   *logrus.Logger
	Database *pgxpool.Pool
	Users    *mongo.Collection
}

type Notification struct {
	Title    string      `json:"title"`
	Body     string      `json:"body"`
	Topic    string      `json:"topic"`
	Priority int         `json:"priority"`
	Data     interface{} `json:"data"`
}

type Topic struct {
	Name   string  `json:"name"`
	Tokens []Token `json:"tokens"`
}

type Token struct {
	UUID  string
	Token string
}
