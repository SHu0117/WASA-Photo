package api

import (
	"net/http"
	"strings"
	"errors"
	"strconv"
)


// Function that extracts the bearer token from the Authorization header
func checkAutorization(authorization string, id int) int{
	var tokens = strings.Split(authorization, " ")
	if len(tokens) == 2 {
		auth := strings.Trim(tokens[1], " ")
		authid, err := strconv.Atoi(auth);
		if err != nil {
			return http.StatusInternalServerError
		}
		if auth == ""{
			return http.StatusForbidden
		}
		if id != authid {
			return http.StatusUnauthorized
		}
	}
	return 0
}

