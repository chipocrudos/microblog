package data

import (
	"context"
	"time"

	"github.com/chipocrudos/microblog/pkg/posts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostRepository struct {
	Data       *Data
	Collection string
	Follow     string
}

func (pr *PostRepository) GetCollection() *mongo.Collection {
	return pr.Data.GetCollection(pr.Collection)
}

func (pr *PostRepository) GetCollectionFollow() *mongo.Collection {
	return pr.Data.GetCollection(pr.Follow)
}

func (pr *PostRepository) GetAll(ctx context.Context, search *posts.SearchPost) (posts.Posts, error) {
	var rPosts posts.Posts

	cnxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	col := pr.GetCollection()

	filter := bson.M{
		"userid": search.UserID,
	}

	option := options.Find()
	option.SetLimit(int64(search.Limit))
	option.SetSort(bson.D{{Key: "date", Value: -1}})
	option.SetSkip((int64(search.Page) - 1) * int64(search.Limit))

	cursor, err := col.Find(cnxt, filter, option)

	if err != nil {
		return rPosts, err
	}

	for cursor.Next(context.TODO()) {
		var post posts.Post
		err := cursor.Decode(&post)

		if err != nil {
			return rPosts, nil
		}

		rPosts = append(rPosts, &post)
	}

	return rPosts, nil
}

func (pr *PostRepository) Create(ctx context.Context, post *posts.Post) error {
	post.CreatedAt = time.Now()

	cnxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	col := pr.GetCollection()

	_, err := col.InsertOne(cnxt, post)

	return err
}

func (pr *PostRepository) GetOne(ctx context.Context, id string) (posts.Post, error) {
	var rPost posts.Post

	return rPost, nil
}

func (pr *PostRepository) Delete(ctx context.Context, post *posts.Post) error {

	cnxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	col := pr.GetCollection()

	filter := bson.M{
		"_id":    post.ID,
		"userid": post.UserID,
	}

	_, err := col.DeleteOne(cnxt, filter)

	return err
}

func (pr *PostRepository) GetRelationPost(ctx context.Context, search *posts.SearchPost) ([]posts.RelatePost, error) {

	cnxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	col := pr.GetCollectionFollow()

	condition := make([]bson.M, 0)
	condition = append(condition, bson.M{"$match": bson.M{"userid": search.UserID}})
	condition = append(condition, bson.M{
		"$lookup": bson.M{
			"from":         pr.Collection,
			"localField":   "followid",
			"foreignField": "userid",
			"as":           "post",
		}})

	condition = append(condition, bson.M{"$unwind": "$post"})
	condition = append(condition, bson.M{"$sort": bson.M{"createdat": -1}})
	condition = append(condition, bson.M{"$skip": (search.Page - 1) * search.Limit})
	condition = append(condition, bson.M{"$limit": search.Limit})

	cur, err := col.Aggregate(cnxt, condition)
	var result []posts.RelatePost

	err = cur.All(cnxt, &result)

	return result, err
}
