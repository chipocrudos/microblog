package users

import "context"

type Repository interface {
	GetAll(ctx context.Context, search *SearchUser) (Users, error)
	GetOne(ctx context.Context, id string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, id string, user User) error
	Delete(ctx context.Context, id string) error
	AuthenticateUser(ctx context.Context, email string, password string) (User, error)
	CreateRelation(ctx context.Context, follow *FollowUser) error
	DeleteRelation(ctx context.Context, follow *FollowUser) error
	RelationStatus(ctx context.Context, follow *FollowUser) bool
}
