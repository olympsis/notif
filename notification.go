package notif

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"
	"github.com/sirupsen/logrus"
)

/*
Create new field service struct
*/
func NewNotificationService(l *logrus.Logger, db *pgxpool.Pool) *Service {
	return &Service{Logger: l, Database: db}
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
	for i := range topic.Tokens {
		s.PushNote(note.Title, note.Body, topic.Tokens[i].Token)
	}
	return nil
}

/*
Create a topic
*/
func (s *Service) CreateTopic(name string) error {
	createStmt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
		ID SERIAL PRIMARY KEY, 
		uuid VARCHAR(50),
		token VARCHAR(100)
	)`, name)

	_, err := s.Database.Query(context.TODO(), createStmt)

	return err
}

/*
Get a topic
*/
func (s *Service) GetTopic(name string) (Topic, error) {
	queryStmt := fmt.Sprintf(`SELECT uuid, token FROM "%s"`, name)
	rows, err := s.Database.Query(context.TODO(), queryStmt)
	if err != nil {
		return Topic{}, err
	}
	defer rows.Close()

	var tokens []Token
	for rows.Next() {
		var ud string
		var tk string
		err = rows.Scan(&ud, &tk)
		if err != nil {
			return Topic{}, err
		}
		token := Token{
			UUID:  ud,
			Token: tk,
		}
		tokens = append(tokens, token)
	}

	topic := Topic{
		Name:   name,
		Tokens: tokens,
	}
	return topic, nil
}

/*
Add tokens to a topic
*/
func (s *Service) AddTokenToTopic(name string, uuid string, token string) error {

	insertStmt := fmt.Sprintf(`INSERT INTO "%s" (uuid, token) VALUES ($1, $2)`, name)

	// add token to topic table
	_, err := s.Database.Exec(context.TODO(), insertStmt, uuid, token)

	return err
}

/*
Remove token from topic
*/
func (s *Service) RemoveTokenFromTopic(name string, token string) error {
	updateStmt := fmt.Sprintf(`DELETE FROM "%s" WHERE token=$1`, name)

	// remove entry from table
	_, err := s.Database.Exec(context.TODO(), updateStmt, token)
	if err != nil {
		return err
	}

	return nil
}

/*
Delete a topic
*/
func (s *Service) DeleteTopic(name string) error {
	deleteSmt := fmt.Sprintf(`DROP TABLE IF EXISTS "%s"`, name)
	_, err := s.Database.Exec(context.TODO(), deleteSmt)

	return err
}

/*
Update User token across all topics
*/
func (s *Service) UpdateUserToken(uuid string, topics []string, token string) error {

	// update topics with the new token
	for _, t := range topics {
		updateStmt := fmt.Sprintf(`
			UPDATE "%s"
			SET token=$1
			WHERE uuid=$2`, t)
		_, err := s.Database.Exec(context.TODO(), updateStmt, token, uuid)
		if err != nil {
			s.Logger.Error(err.Error())
		}
	}

	return nil
}
