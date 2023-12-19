package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"WasaPhoto-1985972/service/database"
	"WasaPhoto-1985972/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// Function that manages the upload of a photo
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")
	pathId = ps.ByName("uid")
	auth := checkAutorization(r.Header.Get("Authorization"), pathId)
	if auth != 0 {
		w.WriteHeader(auth)
		return
	}

	var photo Photo
	photo.User_id = pathId
	photo.Upload_time = time.Now().UTC()
	

	// Create a copy of the body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-upload: error reading body content")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// After reading the body we won't be able to read it again. We'll reassign a "fresh" io.ReadCloser to the body
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	// Check if the body content is either a png or a jpeg image
	err = checkFormatPhoto(r.Body, io.NopCloser(bytes.NewBuffer(data)), ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("photo-upload: body contains file that is neither jpg or png")
		// controllaerrore
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: IMG_FORMAT_ERROR_MSG})
		return
	}

	// Body has been read in the previous function so it's necessary to reassign a io.ReadCloser to it
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	// Generate a unique id for the photo
	photoIdInt, err := rt.db.CreatePhoto(photo.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-upload: error executing db function call")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	photoId := strconv.FormatInt(photoIdInt, 10)

	// Create the user's folder locally to save his/her images
	PhotoPath, err := getUserPhotoFolder(auth)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-upload: error getting user's photo folder")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create an empty file for storing the body content (image)
	out, err := os.Create(filepath.Join(PhotoPath, photoId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("photo-upload: error creating local photo file")
		//  = json.NewEncoder(w).Encode(JSONErrorMsg{Message: INTERNAL_ERROR_MSG})
		return
	}

	// Copy body content to the previously created file
	_, err = io.Copy(out, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("photo-upload: error copying body content into file photo")
		// controllaerrore
		// _ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: INTERNAL_ERROR_MSG})
		return
	}

	// Close the created file
	out.Close()

	w.WriteHeader(http.StatusCreated)
	// controllaerrore
	// _ = json.NewEncoder(w).Encode(PhotoId{IdPhoto: photoIdInt})
	_ = json.NewEncoder(w).Encode(Photo{
		Comments: nil,
		Likes:    nil,
		Owner:    photo.Owner,
		Date:     photo.Date,
		PhotoId:  int(photoIdInt),
	})

}