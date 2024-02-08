package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) CheckBanned(user User, requesterID uint64) (User, error) {
	var banned User
	err := db.c.QueryRow(`SELECT u.id, u.username FROM user u, banning b WHERE b.Banner_id = ? and b.Banned_id = ? AND u.id = ?`, user.ID, requesterID, requesterID).Scan(&banned.ID, &banned.Username)
	if err == nil {
		return banned, ErrUserHasBeenBanned
	} else if errors.Is(err, sql.ErrNoRows) {
		return banned, nil
	}
	return banned, err
}
