package api

import (
	// "bytes"
	"encoding/base64"
	"encoding/json"
	"errors"

	// "image/png"
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
	currentTime := time.Now().UTC()
	photo.Upload_time = currentTime.Format("2006-01-02 15:04:05")
	auth := checkAuthorization(r.Header.Get("Authorization"), photo.User_id)
	if auth != 0 {
		w.WriteHeader(auth)
		return
	}

	// Parse the multipart form
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Read the file data
	data, err := io.ReadAll(file)
	if err != nil {
		// Handle error
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error reading body content")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	encodedString := base64.StdEncoding.EncodeToString(data)
	photo.File = encodedString

	// Generate a unique id for the photo
	dbPhoto, err := rt.db.UploadPhoto(photo.PhotoToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("error executing db function call")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	photo.PhotoFromDatabase(dbPhoto)
	photo.User_username = user.Username
	photo.IsLiked = false
	err = json.NewEncoder(w).Encode(photo)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error while encoding data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	err = rt.db.DeletePhoto(uint64(photoId))
	if err != nil {
		ctx.Logger.WithError(err).Error("can't delete the corrisponding photo ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getUserPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	var list []database.Photo
	dblist, err := rt.db.GetUserPhotos(dbuser, requesterID)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't get the list ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	list = dblist
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error while encoding data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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
	var user User
	user.UserFromDatabase(dbuser)

	auth := checkAuthorization(r.Header.Get("Authorization"), user.ID)
	if auth != 0 {
		w.WriteHeader(auth)
		return
	}

	var list []database.Photo
	dblist, err := rt.db.GetMyStream(user.UserToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("can't get the stream ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	list = dblist
	w.Header().Set("Content-Type", "image/*")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error while encoding data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
