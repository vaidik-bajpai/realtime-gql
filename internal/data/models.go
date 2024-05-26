package data

import (
	"errors"
	"log"

	"github.com/vaidik-bajpai/realtime-gql/internal/prisma/db"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateEmail = errors.New("error duplicate email")
)

type Models struct {
	Users     UserModel
	Messages  MessageModel
	ChatRooms ChatRoomModel
}

func NewModels() Models {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Fatal(err)
	}
	return Models{
		Users:     UserModel{DB: client},
		Messages:  MessageModel{DB: client},
		ChatRooms: ChatRoomModel{DB: client},
	}
}

var Model = NewModels()
