package database

// upload photo and returns photoId
func (db *appdbimpl) CommentPhoto(c Comment) (Comment, error) {

	res, err := db.c.Exec("INSERT INTO comment(texts, user_id, photo_id) VALUES (?, ?, ?)",
		c.Text, c.User_id, c.Photo_id)

	if err != nil {
		// Error executing query
		return c, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return c, err
	}

	c.ID = uint64(lastInsertID)
	return c, nil
}

func (db *appdbimpl) UncommentPhoto(pid uint64, uid uint64, cid uint64) error {

	res, err := db.c.Exec(`DELETE FROM comment WHERE photo_id=? AND user_id = ? AND id = ?`, pid, uid, cid)
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

func (db *appdbimpl) ListComment(pid uint64, requesterID uint64) ([]Comment, error) {

	rows, err := db.c.Query(`SELECT c.id, c.user_id, u.username, c.photo_id, p.user_id, c.texts  FROM photos p, comment c, user u 
				WHERE p.id = ? AND p.id = c.photo_id AND c.user_id = u.id AND c.user_id NOT IN (SELECT Banner_id
					FROM banning WHERE Banned_id = ?)`, pid, requesterID)
	if err != nil {
		return nil, rows.Err()
	}
	// Wait for the function to finish before closing rows.
	defer func() { _ = rows.Close() }()

	// Read all the users in the resulset (users that follow the requesting user)
	var listComment []Comment
	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.ID, &comment.User_id, &comment.User_username, &comment.Photo_id, &comment.Photo_user, &comment.Text)
		if err != nil {
			return nil, err
		}
		listComment = append(listComment, comment)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return listComment, nil
}
