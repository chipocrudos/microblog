package posts

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SearchPost struct {
	UserID primitive.ObjectID `json:"userid"`
	Page   uint64             `json:"page"`
	Limit  uint64             `json:"limit"`
}

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userid" json:"userid,omitempty"`
	Body      string             `bson:"body" json:"body,omitempty"`
	CreatedAt time.Time          `bson:"createdat" json:"date,omitempty"`
}

type Posts []*Post

type RelatePost struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	UserID   string             `bson:"userid" json:"userid,omitempty"`
	FollowID string             `bson:"followid" json:"followid,omitempty"`
	Post     Post
}
