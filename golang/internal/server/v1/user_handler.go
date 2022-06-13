package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chipocrudos/microblog/pkg/response"
	"github.com/chipocrudos/microblog/pkg/token"
	"github.com/chipocrudos/microblog/pkg/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRouter struct {
	Repository users.Repository
}

func (ur *UserRouter) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var search users.SearchUser

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	userid := r.URL.Query().Get("userid")

	search.Filter = r.URL.Query().Get("filter")
	search.Type = users.FollowType(r.URL.Query().Get("type"))

	if len(page) < 1 {
		search.Page = 1
	} else {
		page, err := strconv.Atoi(page)
		if err != nil {
			page = 1
		}
		search.Page = uint64(page)
	}

	if len(limit) < 1 {
		search.Limit = 10
	} else {
		limit, err := strconv.Atoi(limit)
		if err != nil {
			limit = 10
		}
		search.Limit = uint64(limit)
	}

	if len(userid) < 1 {
		_, claims, _ := jwtauth.FromContext(ctx)
		search.UserID, _ = primitive.ObjectIDFromHex(fmt.Sprintf("%v", claims["id"]))

	} else {
		search.UserID, _ = primitive.ObjectIDFromHex(fmt.Sprintf("%v", userid))
	}

	users, err := ur.Repository.GetAll(ctx, &search)
	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"users": users})

}

func (ur *UserRouter) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var u users.User
	err := json.NewDecoder(r.Body).Decode(&u)
	defer r.Body.Close()

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	err = ur.Repository.Create(ctx, &u)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	u.Password = ""
	w.Header().Add("Location", fmt.Sprintf("%s%d", r.URL.String(), u.ID))
	response.JSON(w, r, http.StatusCreated, nil)

}

func (ur *UserRouter) GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")

	ctx := r.Context()
	u, err := ur.Repository.GetOne(ctx, idStr)

	if err != nil {
		response.HTTPError(w, r, http.StatusForbidden, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, u)
}

func (ur *UserRouter) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var u users.User
	err := json.NewDecoder(r.Body).Decode(&u)
	defer r.Body.Close()

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	usr, err := ur.Repository.AuthenticateUser(ctx, u.Email, u.Password)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tokenString := token.JwtEncode(usr)

	response.JSON(w, r, http.StatusOK, response.Map{"access": tokenString})

}

func (ur *UserRouter) GetMeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, claims, _ := jwtauth.FromContext(ctx)

	u, err := ur.Repository.GetOne(ctx, fmt.Sprintf("%v", claims["id"]))

	if err != nil {
		response.HTTPError(w, r, http.StatusForbidden, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, u)
}

func (ur *UserRouter) UpdateHandler(w http.ResponseWriter, r *http.Request) {

	var u users.User
	err := json.NewDecoder(r.Body).Decode(&u)
	defer r.Body.Close()

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	_, claims, _ := jwtauth.FromContext(ctx)

	err = ur.Repository.Update(ctx, fmt.Sprintf("%v", claims["id"]), u)
	if err != nil {
		response.HTTPError(w, r, http.StatusForbidden, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, nil)

}

func (ur *UserRouter) CreateRelationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, claims, _ := jwtauth.FromContext(ctx)
	ID := r.URL.Query().Get("id")

	if len(ID) < 1 {
		response.HTTPError(w, r, http.StatusBadRequest, "Follower id required")
	}

	var rel users.FollowUser

	rel.UserID, _ = primitive.ObjectIDFromHex(fmt.Sprintf("%v", claims["id"]))
	rel.FollowID, _ = primitive.ObjectIDFromHex(fmt.Sprintf("%v", ID))

	exist := ur.Repository.RelationStatus(ctx, &rel)

	if exist == true {
		response.HTTPError(w, r, http.StatusBadRequest, "Already follow user")
		return
	}

	err := ur.Repository.CreateRelation(ctx, &rel)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, r, http.StatusCreated, nil)
}

func (ur *UserRouter) DeleteRelationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, claims, _ := jwtauth.FromContext(ctx)
	ID := r.URL.Query().Get("id")

	if len(ID) < 1 {
		response.HTTPError(w, r, http.StatusBadRequest, "Follower id required")
	}

	var rel users.FollowUser

	rel.UserID, _ = primitive.ObjectIDFromHex(fmt.Sprintf("%v", claims["id"]))
	rel.FollowID, _ = primitive.ObjectIDFromHex(fmt.Sprintf("%v", ID))

	exist := ur.Repository.RelationStatus(ctx, &rel)

	if exist == false {
		response.HTTPError(w, r, http.StatusBadRequest, "Already not follow user")
		return
	}

	err := ur.Repository.DeleteRelation(ctx, &rel)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, nil)
}
