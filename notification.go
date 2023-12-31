package notif

import (
	"context"
	"errors"
	"fmt"

	"github.com/olympsis/models"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
Create new field service struct
  - l: logger for debugging
  - db: notifications database
  - users: main server users database
*/
func NewNotificationService(l *logrus.Logger, db *mongo.Collection, users *mongo.Collection) *Service {
	return &Service{Logger: l, Database: db, Users: users}
}

/*
Create apns client from p12 file
*/
func (p *Service) CreateNewClient(keyID string, teamID string, fileName string) error {
	key, err := token.AuthKeyFromFile(fileName)
	if err != nil {
		return err
	}

	token := token.Token{
		AuthKey: key,
		KeyID:   keyID,
		TeamID:  teamID,
	}

	p.Client = apns2.NewTokenClient(&token).Production()
	return nil
}

/*
Create push note to apns
*/
func (s *Service) PushNote(t string, b string, tk string) error {
	notification := &apns2.Notification{}
	notification.DeviceToken = tk
	notification.Topic = "com.coronislabs.Olympsis"
	notification.Payload = payload.NewPayload().AlertTitle(t).AlertBody(b).Badge(1)
	notification.Priority = 5

	res, err := s.Client.Push(notification)
	if err != nil {
		return err
	}

	if res.Sent() {
		s.Logger.Debug("Sent:", res.ApnsID)
		return nil
	} else {
		errString := fmt.Sprintf("Not Sent: %v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
		return errors.New(errString)
	}
}

/*
Send Notification to token
*/
func (s *Service) SendNotificationToToken(note *Notification, tk string) error {
	err := s.PushNote(note.Title, note.Body, tk)
	if err != nil {
		return err
	}
	return nil
}

/*
Send Notification to topic
*/
func (s *Service) SendNotificationToTopic(note *Notification) error {
	topic, err := s.GetTopic(note.Topic)
	if err != nil {
		return err
	}
	users, err := s.FindUsers(topic.Users)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return errors.New("no users found")
	}
	for i := range users {
		s.PushNote(note.Title, note.Body, users[i].DeviceToken)
	}
	return nil
}

/*
Create a topic
*/
func (s *Service) CreateTopic(name string) error {
	topic := Topic{
		Name:  name,
		Users: []string{},
	}
	_, err := s.Database.InsertOne(context.Background(), topic)
	if err != nil {
		return err
	}

	return nil
}

/*
Get a topic
*/
func (s *Service) GetTopic(name string) (Topic, error) {
	filter := bson.M{
		"name": name,
	}
	var topic Topic
	err := s.Database.FindOne(context.Background(), filter).Decode(&topic)
	if err != nil {
		return Topic{}, err
	}

	return topic, nil
}

/*
Add tokens to a topic
*/
func (s *Service) AddTokenToTopic(name string, uuid string) error {
	filter := bson.M{
		"name": name,
	}
	changes := bson.M{
		"$push": bson.M{
			"users": uuid,
		},
	}
	_, err := s.Database.UpdateOne(context.Background(), filter, changes)
	if err != nil {
		return err
	}
	return nil
}

/*
Remove user from topic
*/
func (s *Service) RemoveTokenFromTopic(name string, uuid string) error {
	filter := bson.M{
		"name": name,
	}
	changes := bson.M{
		"$pull": bson.M{
			"users": uuid,
		},
	}
	_, err := s.Database.UpdateOne(context.Background(), filter, changes)
	if err != nil {
		return err
	}
	return nil
}

/*
Delete a topic
*/
func (s *Service) DeleteTopic(name string) error {
	filter := bson.M{"name": name}
	_, err := s.Database.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

/*
Get Users Data to fetch token
*/
func (s *Service) FindUsers(arr []string) ([]models.User, error) {
	filter := bson.M{
		"uuid": bson.M{
			"$in": arr,
		},
	}
	cursor, err := s.Users.Find(context.Background(), filter)
	if err != nil {
		return []models.User{}, err
	}

	var users []models.User
	for cursor.Next(context.Background()) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			s.Logger.Error("Failed to decode user!")
		}
		users = append(users, user)
	}

	return users, nil
}
