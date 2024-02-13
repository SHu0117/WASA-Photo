package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetUserPhotos(u User, requesterID uint64) ([]Photo, error) {
	var list []Photo
	rows, err := db.c.Query("SELECT p.id, p.user_id, u.username, p.file, p.upload_time FROM photos p, user u WHERE u.id = ? AND p.user_id = u.id ", u.ID)
	if err != nil {
		return nil, ErrDataDoesNotExist
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var p Photo
		err = rows.Scan(&p.ID, &p.User_id, &p.User_username, &p.File, &p.Upload_time)
		if err != nil {
			return nil, err
		}
		isliked, err := db.CheckIfLiked(p.ID, requesterID)
		if err != nil {
			return nil, err
		}
		p.IsLiked = isliked
		err = db.c.QueryRow("SELECT count(*) FROM like WHERE photo_id = ? ", p.ID).Scan(&p.N_likes)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrDataDoesNotExist
			}
		}
		err = db.c.QueryRow("SELECT count(*) FROM comment WHERE photo_id = ? ", p.ID).Scan(&p.N_comments)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrDataDoesNotExist
			}
		}
		list = append(list, p)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return list, nil
}

func (db *appdbimpl) CountPhotos(u User) (int, error) {
	var count int
	err := db.c.QueryRow("SELECT count(*) FROM photos WHERE user_id = ?", u.ID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return count, ErrDataDoesNotExist
		}
	}
	return count, nil
}

func (db *appdbimpl) GetMyStream(u User) ([]Photo, error) {
	var list []Photo
	rows, err := db.c.Query(`SELECT p.id, p.user_id, u.username, p.file, p.upload_time FROM photos p, user u
							 WHERE p.user_id = u.id AND p.user_id IN (SELECT Followed_id FROM following WHERE Follower_id = ? AND Followed_id NOT IN (SELECT Banner_id
										FROM banning WHERE Banned_id = ?)) ORDER BY upload_time DESC`, u.ID, u.ID)
	if err != nil {
		return nil, ErrDataDoesNotExist
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var p Photo
		err = rows.Scan(&p.ID, &p.User_id, &p.User_username, &p.File, &p.Upload_time)
		if err != nil {
			return nil, err
		}
		isliked, err := db.CheckIfLiked(p.ID, u.ID)
		if err != nil {
			return nil, err
		}
		p.IsLiked = isliked
		err = db.c.QueryRow("SELECT count(*) FROM like WHERE photo_id = ?", p.ID).Scan(&p.N_likes)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrDataDoesNotExist
			}
		}
		err = db.c.QueryRow("SELECT count(*) FROM comment WHERE photo_id = ? ", p.ID).Scan(&p.N_comments)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrDataDoesNotExist
			}
		}
		list = append(list, p)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return list, nil
}
