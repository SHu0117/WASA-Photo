package database

import "database/sql"

func (db *appdbimpl) CheckBanned(user User, requesterID uint64) (User, error) {
	var res User
	err := db.c.QueryRow(`SELECT u.id, u.username FROM users u, banning b WHERE b.Banner_id = ? and b.Banned_id = ? AND u.id = ?`, user.ID, requesterID, requesterID).Scan(&res.ID, &res.Username)
	if err == nil {
		return res, ErrUserHasBeenBanned
	} else if err == sql.ErrNoRows {
		return user, nil
	}
	return res, err
}
