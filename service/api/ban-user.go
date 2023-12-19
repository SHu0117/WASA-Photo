package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"errors"
	"WasaPhoto-1985972/service/database"
	"WasaPhoto-1985972/service/api/reqcontext"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	
	pathId, err := strconv.Atoi(ps.ByName("uid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var banning Banning
	banning.Banner_id = uint64(pathId)

	if banning.Banner_id == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = rt.db.ExistUID(banning.Banner_id )
	if errors.Is(err, database.ErrDataDoesNotExist){
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pathBannedId, err := strconv.Atoi(ps.ByName("banneduid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	banning.Banned_id = uint64(pathBannedId)
	if banning.Banned_id == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.db.ExistUID(banning.Banned_id)
	if errors.Is(err, database.ErrDataDoesNotExist){
		w.WriteHeader(http.StatusNotFound)
		return
	}

	auth := checkAutorization(r.Header.Get("Authorization"), banning.Banner_id)
	if auth != 0 {
		w.WriteHeader(auth)
		return
	}

	// Check if the user is trying to ban himself/herself
	if pathId == pathBannedId {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbbanning, err := rt.db.BanUser(banning.BanningToDatabase())
	if err != nil {
		// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
		// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
		ctx.Logger.WithError(err).Error("can't follow")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	

	// Here we can re-use `fountain` as FromDatabase is overwriting every variabile in the structure.
	banning.BanningFromDatabase(dbbanning)

	// Send the output to the user.
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(banning)
}