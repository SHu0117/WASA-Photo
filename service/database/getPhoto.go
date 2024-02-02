package database

import "database/sql"

// upload photo and returns photoId
func (db *appdbimpl) GetPhoto(pid uint64) (Photo, error) {
	var p Photo
	err := db.c.QueryRow("SELECT id, user_id, file, upload_time FROM photo WHERE id = ? ", pid).Scan(&p.ID, &p.User_id, &p.File, &p.Upload_time)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, ErrDataDoesNotExist
		}
	}
	err = db.c.QueryRow("SELECT count(*) FROM photo p, like l WHERE id = ? and l.photo_id = ? ", pid, pid).Scan(&p.N_likes)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, ErrDataDoesNotExist
		}
	}
	err = db.c.QueryRow("SELECT count(*) FROM photo p, commet c WHERE id = ? and c.photo_id = ? ", pid, pid).Scan(&p.N_comments)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, ErrDataDoesNotExist
		}
	}
	return p, nil
}

func (db *appdbimpl) GetUserPhotos(u User) ([]Photo, error) {
	var list []Photo
	rows, err := db.c.Query("SELECT id, user_id, file, upload_time FROM photo WHERE user_id = ? ", u.ID)
	if err != nil {
		return nil, ErrDataDoesNotExist
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var p Photo
		err = rows.Scan(&p.ID, &p.User_id, &p.File, &p.Upload_time)
		if err != nil {
			return nil, err
		}
		err = db.c.QueryRow("SELECT count(*) FROM photo p, like l WHERE id = ? and l.photo_id = ? ", p.ID, p.ID).Scan(&p.N_likes)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrDataDoesNotExist
			}
		}
		err = db.c.QueryRow("SELECT count(*) FROM photo p, comment c WHERE id = ? and c.photo_id = ? ", p.ID, p.ID).Scan(&p.N_comments)
		if err != nil {
			if err == sql.ErrNoRows {
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

func (db *appdbimpl) GetMyStream(u User) ([]Photo, error) {
	var list []Photo
	rows, err := db.c.Query(`SELECT * FROM photo 
							 WHERE user_id IN (SELECT Followed_id FROM following WHERE Follower_id = ? AND Followed_id NOT IN (SELECT Banner_id
										FROM banning WHERE Banned_id = ?)) ORDER BY upload_time DESC`, u.ID, u.ID)
	if err != nil {
		return nil, ErrDataDoesNotExist
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var p Photo
		err = rows.Scan(&p.ID, &p.User_id, &p.File, &p.Upload_time)
		if err != nil {
			return nil, err
		}
		err = db.c.QueryRow("SELECT count(*) FROM photo p, like l WHERE id = ? and l.photo_id = ? ", p.ID, p.ID).Scan(&p.N_likes)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrDataDoesNotExist
			}
		}
		err = db.c.QueryRow("SELECT count(*) FROM photo p, comment c WHERE id = ? and c.photo_id = ? ", p.ID, p.ID).Scan(&p.N_comments)
		if err != nil {
			if err == sql.ErrNoRows {
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
