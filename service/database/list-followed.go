package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) ListFollowed(user User) ([]User, error) {

	rows, err := db.c.Query("SELECT u.id, u.username FROM user u, following f WHERE f.Follower_id = ? and f.followed_id = u.id", user.ID)
	if err != nil {
		return nil, rows.Err()
	}
	// Wait for the function to finish before closing rows.
	defer func() { _ = rows.Close() }()

	// Read all the users in the resulset (users that follow the requesting user)
	var followedUsers []User
	for rows.Next() {
		var followed User
		err = rows.Scan(&followed.ID, &followed.Username)
		if err != nil {
			return nil, err
		}
		followedUsers = append(followedUsers, followed)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return followedUsers, nil
}

func (db *appdbimpl) CountFollowed(user User) (int, error) {
	var count int
	err := db.c.QueryRow("SELECT count(*) FROM following f WHERE f.Follower_id = ?", user.ID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return count, ErrDataDoesNotExist
		}
	}
	return count, nil
}

func (db *appdbimpl) CountFollower(user User) (int, error) {
	var count int
	err := db.c.QueryRow("SELECT count(*) FROM following f WHERE f.followed_id = ?", user.ID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return count, ErrDataDoesNotExist
		}
	}
	return count, nil
}
