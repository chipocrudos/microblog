package token

import (
	"time"

	"github.com/chipocrudos/microblog/config"
	"github.com/chipocrudos/microblog/pkg/response"
	"github.com/chipocrudos/microblog/pkg/users"

	"github.com/go-chi/jwtauth/v5"
)

var TokenAuth *jwtauth.JWTAuth

func init() {
	TokenAuth = jwtauth.New("HS256", []byte(config.Config.JWT_SALT), nil)
}

func JwtEncode(u users.User) string {

	now := time.Now()

	claims := response.Map{
		"id":  u.ID,
		"iat": now.Unix(),
		"exp": now.Add(time.Duration(config.Config.EXP_TIME) * time.Minute),
	}

	_, tokenString, _ := TokenAuth.Encode(claims)
	return tokenString
}
