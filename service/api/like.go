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

// Function that manages the upload of a photo
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user User
	pathUsername := ps.ByName("likeUsername")
	err := rt.db.ExistUsername(pathUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dbuser, err := rt.db.GetUserID(pathUsername)
	if errors.Is(err, database.ErrDataDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user.UserFromDatabase(dbuser)

	var like Like
	auth := checkAuthorization(r.Header.Get("Authorization"), user.ID)
	if auth != 0 {
		w.WriteHeader(auth)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	like.User_id = user.ID
	photoId, err1 := strconv.Atoi(ps.ByName("pid"))
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.db.ExistPhoto(uint64(photoId))
	if errors.Is(err, database.ErrDataDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	like.Photo_id = uint64(photoId)

	pathOwner := ps.ByName("username")
	err = rt.db.ExistUsername(pathOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = rt.db.ExistUsername(pathUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbuser, err = rt.db.GetUserID(pathOwner)
	if errors.Is(err, database.ErrDataDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user.UserFromDatabase(dbuser)
	like.Photo_user = user.ID

	// Generate a unique id for the photo
	dblike, err := rt.db.LikePhoto(like.LikeToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("there's an error executing db function call")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	like.LikeFromDatabase(dblike)
	// controllaerrore
	// _ = json.NewEncoder(w).Encode(PhotoId{IdPhoto: photoIdInt})
	_ = json.NewEncoder(w).Encode(like)
}

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user User
	pathUsername := ps.ByName("likeUsername")
	err := rt.db.ExistUsername(pathUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dbuser, err := rt.db.GetUserID(pathUsername)
	if errors.Is(err, database.ErrDataDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user.UserFromDatabase(dbuser)

	auth := checkAuthorization(r.Header.Get("Authorization"), user.ID)
	if auth != 0 {
		w.WriteHeader(auth)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	photoId, err1 := strconv.Atoi(ps.ByName("pid"))
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.db.ExistPhoto(uint64(photoId))
	if errors.Is(err, database.ErrDataDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = rt.db.UnlikePhoto(uint64(photoId), user.ID)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't delete the corrisponding like ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getPhotoLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var list []database.User
	var user User
	pathUsername := ps.ByName("username")
	dbuser, err := rt.db.GetUserID(pathUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.UserFromDatabase(dbuser)
	requesterID := getToken(r.Header.Get("Authorization"))
	err = rt.db.ExistUID(requesterID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	photoId, err1 := strconv.Atoi(ps.ByName("pid"))
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.db.ExistPhoto(uint64(photoId))
	if errors.Is(err, database.ErrDataDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	dbuser, err = rt.db.CheckBanned(user.UserToDatabase(), requesterID)
	if err != nil {
		if errors.Is(err, database.ErrUserHasBeenBanned) {
			ctx.Logger.WithError(err).Error("can't get list")
			w.WriteHeader(http.StatusForbidden)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	dblist, err := rt.db.ListLikes(uint64(photoId))
	if err != nil {
		ctx.Logger.WithError(err).Error("can't get list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	list = dblist
	// set the header and return output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}
