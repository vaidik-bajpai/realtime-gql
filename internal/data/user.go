package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vaidik-bajpai/realtime-gql/internal/prisma/db"
	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type User struct {
	ID        int64    `json:"id"`
	FirstName *string  `json:"firstname"`
	LastName  *string  `json:"lastname"`
	Password  password `json:"-"`
	Email     *string  `json:"email"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type UserModel struct {
	DB *db.PrismaClient
}

func (m UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	newUser, err := m.DB.User.CreateOne(
		db.User.Firstname.Set(*user.FirstName),
		db.User.Lastname.Set(*user.LastName),
		db.User.Email.Set(*user.Email),
		db.User.Password.Set(user.Password.hash),
	).Exec(ctx)

	if err != nil {
		infoUnique, isErr := db.IsErrUniqueConstraint(err)

		switch {
		case isErr:
			for _, field := range infoUnique.Fields {
				if field == "email" {
					return ErrDuplicateEmail
				} else {
					return errors.New("unique constraint violated")
				}
			}
		default:
			return err
		}
	}

	user.ID = int64(newUser.ID)

	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err := m.DB.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	var returnUser User

	returnUser.Password.hash = user.Password
	returnUser.FirstName = &user.Firstname
	returnUser.LastName = &user.Lastname
	returnUser.Email = &user.Email
	returnUser.ID = int64(user.ID)

	fmt.Println(returnUser.ID, "from get by email")

	return &returnUser, nil
}

func (m UserModel) Get(userId int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err := m.DB.User.FindUnique(
		db.User.ID.Equals(userId),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	returnUser := &User{
		ID:        int64(user.ID),
		FirstName: &user.Firstname,
		LastName:  &user.Lastname,
		Email:     &user.Email,
	}

	returnUser.Password.hash = user.Password

	fmt.Println(returnUser)

	return returnUser, nil
}
