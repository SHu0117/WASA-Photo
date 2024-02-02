package api

import (
	"net/http"
	"strconv"
	"strings"
)

/*
* Function that extracts the bearer token from the Authorization header

	func checkAutorization(authorization string, id uint64) int {
		var tokens = strings.Split(authorization, " ")
		if len(tokens) == 2 {
			auth := strings.Trim(tokens[1], " ")
			authid, err := strconv.Atoi(auth)
			if err != nil {
				return http.StatusInternalServerError
			}
			authId := uint64(authid)
			if auth == "" {
				return http.StatusForbidden
			}
			if id != authId {
				return http.StatusUnauthorized
			}
		}
		return 0
	}

*
*/

func getToken(authorization string) uint64 {
	tokens := strings.Split(authorization, " ")
	if len(tokens) != 2 || tokens[0] != "Bearer" {
		return http.StatusUnauthorized // or http.StatusBadRequest if you want to be more specific
	}

	token := strings.TrimSpace(tokens[1])
	if token == "" {
		return http.StatusForbidden
	}

	authID, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return http.StatusInternalServerError
	}

	return authID
}

func checkAuthorization(authorization string, id uint64) int {
	tokens := strings.Split(authorization, " ")
	if len(tokens) != 2 || tokens[0] != "Bearer" {
		return http.StatusUnauthorized // or http.StatusBadRequest if you want to be more specific
	}

	token := strings.TrimSpace(tokens[1])
	if token == "" {
		return http.StatusForbidden
	}

	authID, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return http.StatusInternalServerError
	}

	if id != authID {
		return http.StatusUnauthorized
	}

	return 0
}
