package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SHu0117/WASA-Photo/service/api/reqcontext"
	"github.com/SHu0117/WASA-Photo/service/database"
	"github.com/julienschmidt/httprouter"
)

// Function that manages the upload of a photo
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")
	pathId, err := strconv.Atoi(ps.ByName("uid"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.db.ExistUID(uint64(pathId))
	if errors.Is(err, database.ErrDataDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	auth := checkAuthorization(r.Header.Get("Authorization"), uint64(pathId))
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
