package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/SHu0117/WASA-Photo/service/api/reqcontext"
	"github.com/SHu0117/WASA-Photo/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	pathId, err := strconv.Atoi(ps.ByName("uid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var following Following
	following.Follower_id = uint64(pathId)

	if following.Follower_id == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.db.ExistUID(following.Follower_id)
	if errors.Is(err, database.ErrDataDoesNotExist) {
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

	err = rt.db.ExistUID(following.Followed_id)
	if errors.Is(err, database.ErrDataDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	auth := checkAuthorization(r.Header.Get("Authorization"), following.Follower_id)
	if auth != 0 {
		w.WriteHeader(auth)
		return
	}

	// Check if the user is trying to follow himself/herself
	if pathId == pathFollowedId {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbfollowing, err := rt.db.FollowUser(following.FollowingToDatabase())
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
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(following)
}
