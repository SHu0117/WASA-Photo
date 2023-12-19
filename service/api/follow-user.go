package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	
	pathId, err := strconv.Atoi(ps.ByName("uid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var following Following
	following.Follower_id := uint64(pathId)

	if follwing.Follower_id == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbfollowing, err := rt.db.ExistUID(follwing.Follower_id )
	if errors.Is(err, database.ErrDataDoesNotExist){
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pathFollowedId, err := strconv.Atoi(ps.ByName("uid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	following.Followed_id := uint64(pathFollowedId)
	if follwing.Followed_id == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbfollowing, err := rt.db.ExistUID(following.Followed_id )
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


	dbfollowing, err := rt.db.FollowUser(Following.FollowingToDatabase())
	if err != nil {
		// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
		// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
		ctx.Logger.WithError(err).Error("can't follow")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Here we can re-use `fountain` as FromDatabase is overwriting every variabile in the structure.
	following.FollowingFromDatabase(dbfollowing)

	// Send the output to the user.
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(following)
}