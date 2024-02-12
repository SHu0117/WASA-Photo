package api

import (
	"encoding/json"
	"net/http"

	"github.com/SHu0117/WASA-Photo/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	var profile Profile
	requesterID := getToken(r.Header.Get("Authorization"))
	err := rt.db.ExistUID(requesterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pathUsername := ps.ByName("username")
	dbuser, err := rt.db.GetUserID(pathUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.UserFromDatabase(dbuser)
	beingBanned, err := rt.db.CheckBeingBanned(user.UserToDatabase(), requesterID)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't get user profile")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if beingBanned {
		ctx.Logger.WithError(err).Error("can't get profile, you have been banned by the user")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	profile.Requester_id = requesterID
	profile.User_id = user.ID
	profile.Username = user.Username
	n_followers, err := rt.db.CountFollower(user.UserToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	profile.N_followers = n_followers
	n_followed, err := rt.db.CountFollowed(user.UserToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	profile.N_followed = n_followed
	n_photos, err := rt.db.CountPhotos(user.UserToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	profile.N_photo = n_photos

	n_banned, err := rt.db.CountBanned(user.UserToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	profile.N_banned = n_banned

	isbanned, err := rt.db.CheckIfBanned(user.ID, requesterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	profile.Is_banned = isbanned
	isfollowed, err := rt.db.CheckIfFollowed(user.ID, requesterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	profile.Is_followed = isfollowed
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(profile)

}
