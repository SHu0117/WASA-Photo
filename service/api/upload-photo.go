package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"image/png"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/SHu0117/WASA-Photo/service/api/reqcontext"
	"github.com/SHu0117/WASA-Photo/service/database"
	"github.com/julienschmidt/httprouter"
)

// Function that manages the upload of a photo
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user User
	// w.Header().Set("Content-Type", "application/json")
	pathUsername := ps.ByName("username")
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

	var photo Photo
	photo.User_id = user.ID
	photo.Upload_time = time.Now().UTC()
	auth := checkAuthorization(r.Header.Get("Authorization"), uint64(photo.User_id))
	if auth != 0 {
		w.WriteHeader(auth)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("error reading body content")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = io.NopCloser(bytes.NewBuffer(data))

	// Check if the body content is either a png image
	_, errPng := png.Decode(r.Body)
	if errPng != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("body contains file that is not png")
		_ = json.NewEncoder(w).Encode("images must be png")
		return
	}

	// Body has been read in the previous function so it's necessary to reassign a io.ReadCloser to it
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	photo.File = data

	// Generate a unique id for the photo
	dbPhoto, err := rt.db.UploadPhoto(photo.PhotoToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("error executing db function call")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/*")
	w.WriteHeader(http.StatusCreated)
	photo.PhotoFromDatabase(dbPhoto)
	// controllaerrore
	// _ = json.NewEncoder(w).Encode(PhotoId{IdPhoto: photoIdInt})
	_ = json.NewEncoder(w).Encode(photo)
}

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user User
	// w.Header().Set("Content-Type", "application/json")
	pathUsername := ps.ByName("username")
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

	auth := checkAuthorization(r.Header.Get("Authorization"), uint64(user.ID))
	if auth != 0 {
		w.WriteHeader(auth)
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

	err = rt.db.DeletePhoto(uint64(photoId))
	if err != nil {
		ctx.Logger.WithError(err).Error("can't delete the corrisponding photo ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user User
	// w.Header().Set("Content-Type", "application/json")
	pathUsername := ps.ByName("username")
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

	auth := checkAuthorization(r.Header.Get("Authorization"), uint64(user.ID))
	if auth != 0 {
		w.WriteHeader(auth)
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

	var photo Photo
	dbphoto, err := rt.db.GetPhoto(uint64(photoId))
	if err != nil {
		ctx.Logger.WithError(err).Error("can't delete the corrisponding photo ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	photo.PhotoFromDatabase(dbphoto)
	w.Header().Set("Content-Type", "image/*")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(photo)
}

func (rt *_router) getUserPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user User
	// w.Header().Set("Content-Type", "application/json")
	pathUsername := ps.ByName("username")
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
	requesterID := getToken(r.Header.Get("Authorization"))
	err = rt.db.ExistUID(requesterID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dbuser, err = rt.db.CheckBanned(user.UserToDatabase(), requesterID)
	if err != nil {
		if errors.Is(err, database.ErrUserHasBeenBanned) {
			ctx.Logger.WithError(err).Error("can't get the list")
			w.WriteHeader(http.StatusForbidden)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	var list []database.Photo
	dblist, err := rt.db.GetUserPhotos(user.UserToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("can't get the list ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	list = dblist
	w.Header().Set("Content-Type", "image/*")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user User
	// w.Header().Set("Content-Type", "application/json")
	pathUsername := ps.ByName("username")
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

	auth := checkAuthorization(r.Header.Get("Authorization"), uint64(user.ID))
	if auth != 0 {
		w.WriteHeader(auth)
		return
	}

	var list []database.Photo
	dblist, err := rt.db.GetMyStream(user.UserToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("can't get the list ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	list = dblist
	w.Header().Set("Content-Type", "image/*")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}
