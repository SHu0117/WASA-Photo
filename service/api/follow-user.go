package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SHu0117/WASA-Photo/service/api/reqcontext"
	"github.com/SHu0117/WASA-Photo/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	pathUsername := ps.ByName("username")
	dbuser, err := rt.db.GetUserID(pathUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user User
	user.UserFromDatabase(dbuser)

	var following Following
	following.Follower_id = user.ID

	pathFollowUsername := ps.ByName("followUsername")
	dbuser, err = rt.db.GetUserID(pathFollowUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.UserFromDatabase(dbuser)
	following.Followed_id = user.ID

	ifFollowed, err := rt.db.CheckIfFollowed(following.Followed_id, following.Follower_id)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't execute")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ifFollowed {
		ctx.Logger.WithError(err).Error("you have already followed the user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	auth := checkAuthorization(r.Header.Get("Authorization"), following.Follower_id)
	if auth != 0 {
		w.WriteHeader(auth)
		return
	}

	// Check if the user is trying to ban himself/herself
	if pathUsername == pathFollowUsername {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbfollowing, err := rt.db.FollowUser(following.FollowingToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("can't follow")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Here we can re-use `follow` as FromDatabase is overwriting every variabile in the structure.
	following.FollowingFromDatabase(dbfollowing)

	// Send the output to the user.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(following)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error while encoding data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	pathUsername := ps.ByName("username")
	dbuser, err := rt.db.GetUserID(pathUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user User
	user.UserFromDatabase(dbuser)
	var following Following
	following.Follower_id = user.ID

	pathFollowUsername := ps.ByName("followUsername")
	dbuser, err = rt.db.GetUserID(pathFollowUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.UserFromDatabase(dbuser)
	following.Followed_id = user.ID

	ifFollowed, err := rt.db.CheckIfFollowed(following.Followed_id, following.Follower_id)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't execute")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ifFollowed {
		ctx.Logger.WithError(err).Error("you have not followed the user before")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	auth := checkAuthorization(r.Header.Get("Authorization"), following.Follower_id)
	if auth != 0 {
		w.WriteHeader(auth)
		return
	}

	// Check if the user is trying to ban himself/herself
	if pathUsername == pathFollowUsername {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.UnfollowUser(following.FollowingToDatabase())
	if errors.Is(err, database.ErrDataDoesNotExist) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("can't unfollow")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

func (rt *_router) listFollowed(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var list []database.User
	pathUsername := ps.ByName("username")
	dbuser, err := rt.db.GetUserID(pathUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requesterID := getToken(r.Header.Get("Authorization"))
	err = rt.db.ExistUID(requesterID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	beingBanned, err := rt.db.CheckBeingBanned(dbuser, requesterID)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't get list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if beingBanned {
		ctx.Logger.WithError(err).Error("can't get list, you have been banned by the user")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	dblist, err := rt.db.ListFollowed(dbuser, requesterID)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't get list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	list = dblist
	// set the header and return output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error while encoding data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rt *_router) listFollower(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var list []database.User
	pathUsername := ps.ByName("username")
	dbuser, err := rt.db.GetUserID(pathUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requesterID := getToken(r.Header.Get("Authorization"))
	err = rt.db.ExistUID(requesterID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	beingBanned, err := rt.db.CheckBeingBanned(dbuser, requesterID)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't get list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if beingBanned {
		ctx.Logger.WithError(err).Error("can't get list, you have been banned by the user")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	dblist, err := rt.db.ListFollower(dbuser, requesterID)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't get list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	list = dblist
	// set the header and return output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error while encoding data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
