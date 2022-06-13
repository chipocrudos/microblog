package users

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type FollowType string

const (
	Any      FollowType = ""
	Follow   FollowType = "follow"
	Unfollow FollowType = "unfollow"
)

type SearchUser struct {
	UserID primitive.ObjectID `json:"userid"`
	Filter string             `json:"search"`
	Type   FollowType         `json:"type"`
	Page   uint64             `json:"page"`
	Limit  uint64             `json:"limit"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName string             `json:"first_name,omitempty"`
	LastName  string             `json:"last_name,omitempty"`
	Email     string             `json:"email,omitempty"`
	Picture   string             `json:"picture,omitempty"`
	Password  string             `json:"password,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty"`
}

func (u *User) HashPassword() error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(passwordHash)

	return nil
}

func (u User) PasswordMatch(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err == nil
}

type Users []*User

type FollowUser struct {
	UserID   primitive.ObjectID `bson:"userid" json:"userid,omitempty"`
	FollowID primitive.ObjectID `bson:"followid" json:"followid,omitempty"`
}
