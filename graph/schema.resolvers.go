package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/pascaldekloe/jwt"
	"github.com/vaidik-bajpai/realtime-gql/graph/model"
	"github.com/vaidik-bajpai/realtime-gql/internal/data"
)

// PostMessage is the resolver for the postMessage field.
func (r *mutationResolver) PostMessage(ctx context.Context, chatRoomID int, content string) (*model.Message, error) {
	user := data.ContextGetUser(ctx)
	if user == data.AnonymousUser {
		return &model.Message{}, errors.New("access denied")
	}
	message := data.Message{
		Content: content,
		Sender: data.User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
		Chat: data.ChatRoom{
			ID: chatRoomID,
		},
	}
	err := data.Model.Messages.CreateMessage(&message)
	if err != nil {
		return &model.Message{}, errors.New("something went wrong can't post the message")
	}

	r.subscriptionResolver.publishMessage(strconv.Itoa(chatRoomID), &model.Message{
		ID:      message.ID,
		Content: message.Content,
		Sender: &model.User{
			Firstname: *message.Sender.FirstName,
			Lastname:  *message.Sender.LastName,
			Email:     *message.Sender.Email,
		},
		Timestamp: message.Timestamp.GoString(),
	})

	return &model.Message{
		ID:      message.ID,
		Content: message.Content,
		Sender: &model.User{
			Firstname: *message.Sender.FirstName,
			Lastname:  *message.Sender.LastName,
			Email:     *message.Sender.Email,
		},
		Timestamp: message.Timestamp.GoString(),
	}, err
}

// CreateChatRoom is the resolver for the createChatRoom field.
func (r *mutationResolver) CreateChatRoom(ctx context.Context, name string) (*model.ChatRoom, error) {
	user := data.ContextGetUser(ctx)
	fmt.Println(user)
	if user == data.AnonymousUser {
		return &model.ChatRoom{}, errors.New("access denied")
	}

	chatRoom := &data.ChatRoom{
		Name: name,
	}
	err := data.Model.ChatRooms.CreateChatRoom(chatRoom)
	if err != nil {
		return &model.ChatRoom{}, err
	}

	return &model.ChatRoom{
		ID:   chatRoom.ID,
		Name: chatRoom.Name,
	}, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	user, err := data.Model.Users.GetByEmail(input.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	fmt.Println(user.ID, "from get by email")

	ok, err := user.Password.Matches(input.Password)
	if !ok || err != nil {
		return "", errors.New("invalid email or password")
	}

	if user == nil {
		return "", fmt.Errorf("user cannot be nil")
	}

	var claims jwt.Claims
	fmt.Println(user.ID)
	claims.Subject = strconv.FormatInt(user.ID, 10)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "bajpai"
	claims.Audiences = []string{"bajpai"}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set")
	}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %v", err)
	}

	return "Bearer " + string(jwtBytes), nil
}

// Signup is the resolver for the signup field.
func (r *mutationResolver) Signup(ctx context.Context, input model.NewUser) (*model.User, error) {
	user := &data.User{
		FirstName: &input.Firstname,
		LastName:  &input.Lastname,
		Email:     &input.Email,
	}

	err := user.Password.Set(input.Password)
	if err != nil {
		return &model.User{}, errors.New("error decoding the password")
	}

	err = data.Model.Users.Insert(user)
	if err != nil {
		log.Fatal(err)
	}

	return &model.User{
		ID:        int(user.ID),
		Firstname: *user.FirstName,
		Lastname:  *user.LastName,
		Email:     *user.Email,
	}, nil
}

// ChatRooms is the resolver for the chatRooms field.
func (r *queryResolver) ChatRooms(ctx context.Context) ([]*model.ChatRoom, error) {
	panic(fmt.Errorf("not implemented: ChatRooms - chatRooms"))
}

// Messages is the resolver for the messages field.
func (r *queryResolver) Messages(ctx context.Context, chatRoomID string) ([]*model.Message, error) {
	panic(fmt.Errorf("not implemented: Messages - messages"))
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	panic(fmt.Errorf("not implemented: Me - me"))
}

// MessagePosted is the resolver for the messagePosted field.
func (r *subscriptionResolver) MessagePosted(ctx context.Context, chatRoomID string) (<-chan *model.Message, error) {
	messageChannel := make(chan *model.Message, 1)

	r.mu.Lock()
	r.observers[chatRoomID] = append(r.observers[chatRoomID], messageChannel)
	r.mu.Unlock()

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		defer r.mu.Unlock()

		subscribers := r.observers[chatRoomID]
		for i := range subscribers {
			if subscribers[i] == messageChannel {
				subscribers = append(subscribers[:i], subscribers[i+1:]...)
				break
			}
		}

		r.observers[chatRoomID] = subscribers
		close(messageChannel)
	}()

	return messageChannel, nil
}

func (r *subscriptionResolver) publishMessage(chatRoomID string, message *model.Message) {
	r.mu.Lock()
	defer r.mu.Unlock()

	subscribers := r.observers[chatRoomID]
	for _, subscriber := range subscribers {
		select {
		case subscriber <- message:
		default:
			log.Println("Subscriber is not ready to receive message")
		}
	}
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver {
	return r.subscriptionResolver
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct {
	mu        sync.Mutex
	observers map[string][]chan *model.Message
}
