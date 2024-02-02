package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SHu0117/WASA-Photo/service/api/reqcontext"
	"github.com/SHu0117/WASA-Photo/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	pathUsername := ps.ByName("username")
	err := rt.db.ExistUsername(pathUsername)
	if errors.Is(err, database.ErrDataDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	dbuser, err := rt.db.GetUserID(pathUsername)
	user.UserFromDatabase(dbuser)

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
