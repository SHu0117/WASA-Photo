package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"WasaPhoto-1985972/service/database"
	"WasaPhoto-1985972/service/api/reqcontext"
	"strconv"
	"errors"
)

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	
	pathId, err := strconv.Atoi(ps.ByName("uid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var following Following
	following.Follower_id = uint64(pathId)

	if following.Follower_id == 0{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.db.ExistUID(following.Follower_id )
	if errors.Is(err, database.ErrDataDoesNotExist){
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pathFollowedId, err := strconv.Atoi(ps.ByName("uid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	following.Followed_id = uint64(pathFollowedId)
	if following.Followed_id == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.db.ExistUID(following.Followed_id )
	if errors.Is(err, database.ErrDataDoesNotExist){
		w.WriteHeader(http.StatusNotFound)
		return
	}

	auth := checkAutorization(r.Header.Get("Authorization"), following.Follower_id)
	if auth != 0 {
		w.WriteHeader(auth)
		return
	}

	// Check if the user is trying to ban himself/herself
	if pathId == pathFollowedId {
		w.WriteHeader(http.StatusBadRequest)
		return
	}


	err = rt.db.UnfollowUser(following.FollowingToDatabase())
	if err != nil {
		// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
		// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
		ctx.Logger.WithError(err).Error("can't follow")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	
}