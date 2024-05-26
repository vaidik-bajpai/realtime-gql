// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type ChatRoom struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	Messages []*Message `json:"messages,omitempty"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Message struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	Sender    *User  `json:"sender"`
	Timestamp string `json:"timestamp"`
}

type Mutation struct {
}

type NewUser struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Query struct {
}

type Subscription struct {
}

type User struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}
