package api

import (
	"encoding/json"
	"net/http"

	"github.com/SHu0117/WASA-Photo/service/api/reqcontext"
	"github.com/SHu0117/WASA-Photo/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// The body was not a parseable JSON, reject it
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.ExistUsername(user.Username)
	if err == database.ErrDataDoesNotExist {
		dbuser, err := rt.db.NewUser(user.UserToDatabase())
		if err != nil {
			// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
			// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
			ctx.Logger.WithError(err).Error("can't create the user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		user.UserFromDatabase(dbuser)

		// Send the output to the user.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(user)

	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		dbuser, err := rt.db.GetUserID(user.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		user.UserFromDatabase(dbuser)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(user)
	}

}
