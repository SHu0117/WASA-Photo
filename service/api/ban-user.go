package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SHu0117/WASA-Photo/service/api/reqcontext"
	"github.com/SHu0117/WASA-Photo/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	pathUsername := ps.ByName("username")
	dbuser, err := rt.db.GetUserID(pathUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var banning Banning
	banning.Banner_id = dbuser.ID

	pathBanUsername := ps.ByName("banUsername")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dbuser, err = rt.db.GetUserID(pathBanUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	banning.Banned_id = dbuser.ID

	ifBanned, err := rt.db.CheckIfBanned(banning.Banned_id, banning.Banner_id)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't execute")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ifBanned {
		ctx.Logger.WithError(err).Error("you have already banned the user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	auth := checkAuthorization(r.Header.Get("Authorization"), banning.Banner_id)
	if auth != 0 {
		w.WriteHeader(auth)
		return
	}

	// Check if the user is trying to ban himself/herself
	if pathUsername == pathBanUsername {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbbanning, err := rt.db.BanUser(banning.BanningToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("can't ban")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	banning.BanningFromDatabase(dbbanning)

	// Send the output to the user.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(banning)
}

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	pathUsername := ps.ByName("username")
	dbuser, err := rt.db.GetUserID(pathUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var banning Banning
	banning.Banner_id = dbuser.ID

	pathBanUsername := ps.ByName("banUsername")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dbuser, err = rt.db.GetUserID(pathBanUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	banning.Banned_id = dbuser.ID

	ifBanned, err := rt.db.CheckIfBanned(banning.Banned_id, banning.Banner_id)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't execute")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ifBanned {
		ctx.Logger.WithError(err).Error("you have not banned the user before")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	auth := checkAuthorization(r.Header.Get("Authorization"), banning.Banner_id)
	if auth != 0 {
		w.WriteHeader(auth)
		return
	}

	// Check if the user is trying to ban himself/herself
	if pathUsername == pathBanUsername {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.UnbanUser(banning.BanningToDatabase())
	if errors.Is(err, database.ErrDataDoesNotExist) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("can't unban")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) listBanned(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var list []database.User
	pathUsername := ps.ByName("username")
	dbuser, err := rt.db.GetUserID(pathUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
	dblist, err := rt.db.ListBanned(dbuser)
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
