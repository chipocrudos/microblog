package data

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/chipocrudos/microblog/pkg/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	Data       *Data
	Collection string
	Follow     string
}

func (ur *UserRepository) GetCollection() *mongo.Collection {
	return ur.Data.GetCollection(ur.Collection)
}

func (ur *UserRepository) GetCollectionFollow() *mongo.Collection {
	return ur.Data.GetCollection(ur.Follow)

}

func (ur *UserRepository) GetAll(ctx context.Context, search *users.SearchUser) (users.Users, error) {

	log.Println(search.UserID, search.Filter, search.Type, search.Page, search.Limit)

	us := users.Users{}
	cnxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	col := ur.GetCollection()

	option := options.Find()
	option.SetSkip((int64(search.Page) - 1) * int64(search.Limit))
	option.SetLimit(int64(search.Limit))

	query := bson.M{
		"firstname": bson.M{"$regex": `(?i)` + search.Filter},
		"lastname":  bson.M{"$regex": `(?i)` + search.Filter},
	}

	cur, err := col.Find(cnxt, query, option)

	if err != nil {
		return us, err
	}

	var follower, include bool

	for cur.Next(cnxt) {
		var s users.User

		err := cur.Decode(&s)

		if err != nil {
			return us, err
		}

		var follow users.FollowUser
		follow.UserID = search.UserID
		follow.FollowID = s.ID

		include = false

		follower = ur.RelationStatus(ctx, &follow)

		if search.Type == "follow" && follower == true {
			include = true
		} else if search.Type == "unfollow" && follower == false {
			include = true
		} else if s.ID == search.UserID {
			include = false
		} else if search.Type == "" {
			include = true
		}

		if include == true {
			s.Password = ""
			s.Picture = ""
			s.Email = ""

			us = append(us, &s)
		}

	}

	err = cur.Err()
	if err != nil {
		return us, err
	}
	cur.Close(cnxt)

	return us, nil
}

func (ur *UserRepository) Create(ctx context.Context, u *users.User) error {

	if err := u.HashPassword(); err != nil {
		return err
	}

	cnxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	exist, _ := ur.GetByEmail(ctx, u.Email)
	if len(exist.Email) > 0 {
		return errors.New("User exist")

	}

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	col := ur.GetCollection()
	_, err := col.InsertOne(cnxt, u)

	return err
}

func (ur *UserRepository) Update(ctx context.Context, id string, u users.User) error {

	cnxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	col := ur.GetCollection()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": bson.M{"$eq": objID}}

	updateUser := make(map[string]interface{})
	updateUser["updatedat"] = time.Now()

	if len(u.Email) > 0 {
		updateUser["email"] = u.Email
	}
	if len(u.FirstName) > 0 {
		updateUser["firstname"] = u.FirstName
	}
	if len(u.LastName) > 0 {
		updateUser["lastname"] = u.LastName
	}
	if len(u.Picture) > 0 {
		updateUser["picture"] = u.Picture
	}

	updateString := bson.M{
		"$set": updateUser,
	}

	_, err := col.UpdateOne(cnxt, filter, updateString)

	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) Delete(ctx context.Context, id string) error {

	return nil
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (users.User, error) {
	cnxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var u users.User

	col := ur.GetCollection()

	condition := bson.M{"email": email}

	err := col.FindOne(cnxt, condition).Decode(&u)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (ur *UserRepository) GetOne(ctx context.Context, id string) (users.User, error) {
	cnxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	u := users.User{}

	col := ur.GetCollection()

	idU, _ := primitive.ObjectIDFromHex(id)

	condition := bson.M{"_id": idU}

	err := col.FindOne(cnxt, condition).Decode(&u)
	if err != nil {
		return u, err
	}

	u.Password = ""

	return u, nil
}

func (ur *UserRepository) AuthenticateUser(ctx context.Context, email string, password string) (users.User, error) {

	u, err := ur.GetByEmail(ctx, email)

	if err != nil {
		return u, errors.New("User or password error")
	}

	if u.Email == email && u.PasswordMatch(password) {
		return u, nil
	} else {
		return u, errors.New("User or password error")
	}

}

func (ur *UserRepository) CreateRelation(ctx context.Context, follow *users.FollowUser) error {

	cntx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	col := ur.GetCollectionFollow()

	_, err := col.InsertOne(cntx, follow)

	return err

}

func (ur *UserRepository) DeleteRelation(ctx context.Context, follow *users.FollowUser) error {
	cntx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	col := ur.GetCollectionFollow()

	_, err := col.DeleteOne(cntx, follow)

	return err

}

func (ur *UserRepository) RelationStatus(ctx context.Context, follow *users.FollowUser) bool {

	cntx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	col := ur.GetCollectionFollow()

	filter := bson.M{
		"userid":   follow.UserID,
		"followid": follow.FollowID,
	}

	var result users.FollowUser
	err := col.FindOne(cntx, filter).Decode(&result)
	if err != nil {
		return false
	}
	return true
}
