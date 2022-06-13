package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chipocrudos/microblog/pkg/posts"
	"github.com/chipocrudos/microblog/pkg/response"
	"github.com/go-chi/jwtauth/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostRouter struct {
	Repository posts.Repository
}

func (pr *PostRouter) GetAllHandler(w http.ResponseWriter, r *http.Request) {

	var search posts.SearchPost

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	userid := r.URL.Query().Get("userid")
	ctx := r.Context()

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

	rposts, err := pr.Repository.GetAll(ctx, &search)

	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"posts": rposts})

}

func (pr *PostRouter) CreateHandler(w http.ResponseWriter, r *http.Request) {

	var post posts.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	defer r.Body.Close()

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if len(post.Body) < 1 {
		response.HTTPError(w, r, http.StatusBadRequest, "Body post required")
		return
	}

	ctx := r.Context()
	_, claims, _ := jwtauth.FromContext(ctx)

	post.UserID, _ = primitive.ObjectIDFromHex(fmt.Sprintf("%v", claims["id"]))

	err = pr.Repository.Create(ctx, &post)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Add("Location", fmt.Sprintf("%s%d", r.URL.String(), post.ID))
	response.JSON(w, r, http.StatusCreated, nil)

}

func (pr *PostRouter) DeleteHandler(w http.ResponseWriter, r *http.Request) {

	ID := r.URL.Query().Get("id")

	if len(ID) < 1 {
		response.HTTPError(w, r, http.StatusBadRequest, "Post id required")
	}

	var post posts.Post
	ctx := r.Context()

	_, claims, _ := jwtauth.FromContext(ctx)
	post.ID, _ = primitive.ObjectIDFromHex(fmt.Sprintf("%v", ID))

	post.UserID, _ = primitive.ObjectIDFromHex(fmt.Sprintf("%v", claims["id"]))

	err := pr.Repository.Delete(ctx, &post)

	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return

	}

	response.JSON(w, r, http.StatusOK, nil)
}

func (pr *PostRouter) GetRelatePostHandler(w http.ResponseWriter, r *http.Request) {

	var search posts.SearchPost

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	userid := r.URL.Query().Get("userid")
	ctx := r.Context()

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

	rposts, err := pr.Repository.GetRelationPost(ctx, &search)

	if err != nil {
		response.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, response.Map{"posts": rposts})

}
