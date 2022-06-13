package posts

import (
	"context"
)

type Repository interface {
	GetAll(ctx context.Context, search *SearchPost) (Posts, error)
	GetOne(ctx context.Context, id string) (Post, error)
	Create(ctx context.Context, post *Post) error
	Delete(ctx context.Context, post *Post) error
	GetRelationPost(ctx context.Context, search *SearchPost) ([]RelatePost, error)
}
