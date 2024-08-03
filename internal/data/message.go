package data

import (
	"context"
	"errors"
	"time"

	"github.com/vaidik-bajpai/realtime-gql/graph/model"
	"github.com/vaidik-bajpai/realtime-gql/internal/prisma/db"
)

type Message struct {
	ID        int
	Content   string
	Sender    User
	Timestamp time.Time
	Chat      ChatRoom
}

type MessageModel struct {
	DB *db.PrismaClient
}

func (m MessageModel) CreateMessage(message *Message) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	createdMessage, err := m.DB.Message.CreateOne(
		db.Message.Content.Set(message.Content),
		db.Message.Sender.Link(
			db.User.ID.Equals(int(message.Sender.ID)),
		),
		db.Message.ChatRoom.Link(
			db.ChatRoom.ID.Equals(int(message.Chat.ID)),
		),
	).Exec(ctx)

	if err != nil {
		return errors.New("error creating the message")
	}

	message.ID = createdMessage.ID
	message.Timestamp = createdMessage.Timestamp

	return nil
}

func (m MessageModel) GetAChat(chatRoomID int) ([]*model.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	allMessages, err := m.DB.Message.FindMany(
		db.Message.ChatRoomID.Equals(chatRoomID),
	).With(
		db.Message.Sender.Fetch(),
	).Exec(ctx)

	if err != nil {
		return nil, errors.New("error fetching messages")
	}

	var returnMessages []*model.Message
	for _, message := range allMessages {

		user := message.Sender()
		returnMessages = append(returnMessages, &model.Message{
			ID:      message.ID,
			Content: message.Content,
			Sender: &model.User{
				ID:        user.ID,
				Email:     user.Email,
				Firstname: user.Firstname,
				Lastname:  user.Lastname,
			},
			Timestamp: message.Timestamp.GoString(),
		})
	}
	return returnMessages, nil
}
