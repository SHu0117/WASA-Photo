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

	w.Header().Set("Content-Type", "application/json")
	pathId, err := strconv.Atoi(ps.ByName("uid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = rt.db.ExistUID(uint64(pathId))
	if errors.Is(err, database.ErrDataDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var photo Photo
	photo.User_id = uint64(pathId)
	photo.Upload_time = time.Now().UTC()
	auth := checkAuthorization(r.Header.Get("Authorization"), photo.User_id)
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

	w.WriteHeader(http.StatusCreated)
	photo.PhotoFromDatabase(dbPhoto)
	// controllaerrore
	// _ = json.NewEncoder(w).Encode(PhotoId{IdPhoto: photoIdInt})
	_ = json.NewEncoder(w).Encode(photo)

}
