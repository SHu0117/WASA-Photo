package database

// upload photo and returns photoId
func (db *appdbimpl) LikePhoto(l Like) (Like, error) {

	res, err := db.c.Exec("INSERT INTO like(user_id, photo_id) VALUES (?, ?)",
		l.User_id, l.Photo_id)

	if err != nil {
		// Error executing query
		return l, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return l, err
	}

	l.ID = uint64(lastInsertID)
	return l, nil
}

func (db *appdbimpl) UnlikePhoto(pid uint64, uid uint64) error {

	res, err := db.c.Exec(`DELETE FROM like WHERE photo_id=? AND user_id = ?`, pid, uid)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// If we didn't delete any row, then the row didn't exist
		return ErrDataDoesNotExist
	}
	return nil
}

func (db *appdbimpl) ListLikes(pid uint64) ([]User, error) {

	rows, err := db.c.Query("SELECT u.id, u.username FROM user u, photos p, like l WHERE u.id = l.user_id AND p.id = ? AND p.id = l.photo_id", pid)
	if err != nil {
		return nil, rows.Err()
	}
	// Wait for the function to finish before closing rows.
	defer func() { _ = rows.Close() }()

	// Read all the users in the resulset (users that follow the requesting user)
	var likeUsers []User
	for rows.Next() {
		var liker User
		err = rows.Scan(&liker.ID, &liker.Username)
		if err != nil {
			return nil, err
		}
		likeUsers = append(likeUsers, liker)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return likeUsers, nil
}
