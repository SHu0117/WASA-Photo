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

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
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
	user.ID = uint64(pathId)

	var new string
	err = json.NewDecoder(r.Body).Decode(&new)
	if err != nil {
		ctx.Logger.WithError(err).Error("setting username: error decoding json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.db.ExistUsername(new)
	if err != database.ErrDataDoesNotExist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.ID = uint64(pathId)

	auth := checkAuthorization(r.Header.Get("Authorization"), user.ID)
	if auth != 0 {
		w.WriteHeader(auth)
		return
	}

	err = rt.db.SetUsername(user.UserToDatabase(), new)
	if err != nil {
		ctx.Logger.WithError(err).Error("setting username: error executing update query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
