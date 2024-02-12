package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) ListBanned(user User) ([]User, error) {

	rows, err := db.c.Query("SELECT u.id, u.username FROM user u, banning b WHERE b.Banner_id = ? and b.Banned_id = u.id", user.ID)
	if err != nil {
		return nil, rows.Err()
	}
	// Wait for the function to finish before closing rows.
	defer func() { _ = rows.Close() }()

	// Read all the users in the resulset (users that follow the requesting user)
	var bannedUsers []User
	for rows.Next() {
		var banned User
		err = rows.Scan(&banned.ID, &banned.Username)
		if err != nil {
			return nil, err
		}
		bannedUsers = append(bannedUsers, banned)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return bannedUsers, nil
}

func (db *appdbimpl) CountBanned(user User) (int, error) {
	var count int
	err := db.c.QueryRow("SELECT count(*) FROM banning b WHERE b.Banner_id = ?", user.ID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return count, ErrDataDoesNotExist
		}
	}
	return count, nil
}
