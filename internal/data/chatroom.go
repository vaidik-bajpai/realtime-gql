package data

import (
	"context"
	"time"

	"github.com/vaidik-bajpai/realtime-gql/internal/prisma/db"
)

type ChatRoom struct {
	ID       int
	Name     string
	Messages []Message
}

type ChatRoomModel struct {
	DB *db.PrismaClient
}

func (m ChatRoomModel) CreateChatRoom(chatRoom *ChatRoom) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	newChatRoom, err := m.DB.ChatRoom.CreateOne(
		db.ChatRoom.Name.Set(chatRoom.Name),
	).Exec(ctx)
	if err != nil {
		return err
	}

	chatRoom.ID = newChatRoom.ID

	return nil
}
