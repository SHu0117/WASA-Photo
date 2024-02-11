package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) CheckBeingBanned(user User, requesterID uint64) (bool, error) {
	var res bool
	if err := db.c.QueryRow(`SELECT EXISTS(SELECT * FROM banning WHERE Banner_id = ? AND  Banned_id = ?)`, user.ID, requesterID).Scan(&res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
	}
	return res, nil
}

func (db *appdbimpl) CheckIfFollowed(targetID uint64, requesterID uint64) (bool, error) {
	var res bool
	if err := db.c.QueryRow(`SELECT EXISTS(SELECT * FROM following WHERE Follower_id = ? AND  followed_id = ?)`, requesterID, targetID).Scan(&res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
	}
	return res, nil
}

func (db *appdbimpl) CheckIfBanned(targetID uint64, requesterID uint64) (bool, error) {
	var res bool
	if err := db.c.QueryRow(`SELECT EXISTS(SELECT * FROM banning WHERE Banner_id = ? AND  Banned_id = ?)`, requesterID, targetID).Scan(&res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
	}
	return res, nil
}

func (db *appdbimpl) CheckIfLiked(targetID uint64, requesterID uint64) (bool, error) {
	var res bool
	if err := db.c.QueryRow(`SELECT EXISTS(SELECT * FROM like WHERE user_id = ? AND  photo_id = ?)`, requesterID, targetID).Scan(&res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
	}
	return res, nil
}
