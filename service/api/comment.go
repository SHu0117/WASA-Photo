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
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	requesterID := getToken(r.Header.Get("Authorization"))
	err := rt.db.ExistUID(requesterID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var comment Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	auth := checkAuthorization(r.Header.Get("Authorization"), requesterID)
	if auth != 0 {
		w.WriteHeader(auth)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	comment.User_id = requesterID
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
	comment.Photo_id = uint64(photoId)

	pathOwner := ps.ByName("username")
	err = rt.db.ExistUsername(pathOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user User
	dbuser, err := rt.db.GetUserID(pathOwner)
	if errors.Is(err, database.ErrDataDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user.UserFromDatabase(dbuser)

	comment.Photo_user = user.ID

	// Generate a unique id for the photo
	dbcomment, err := rt.db.CommentPhoto(comment.CommentToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("error executing db function call")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	comment.CommentFromDatabase(dbcomment)
	// controllaerrore
	// _ = json.NewEncoder(w).Encode(PhotoId{IdPhoto: photoIdInt})
	_ = json.NewEncoder(w).Encode(comment)
}

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	pathUsername := ps.ByName("commentUsername")
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
	var user User
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

	err = rt.db.UncommentPhoto(uint64(photoId), user.ID)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't delete the corrisponding comment ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getPhotoComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var list []database.Comment
	pathUsername := ps.ByName("username")
	dbuser, err := rt.db.GetUserID(pathUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user User
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

	beingBanned, err := rt.db.CheckBeingBanned(user.UserToDatabase(), requesterID)
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
	dblist, err := rt.db.ListComment(uint64(photoId))
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
