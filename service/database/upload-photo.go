package database

// upload photo and returns photoId
func (db *appdbimpl) UploadPhoto(p Photo) (Photo, error) {

	res, err := db.c.Exec("INSERT INTO photo (user_id, file, upload_time) VALUES (?, ?, ?)",
		p.User_id, p.File, p.Upload_time)

	if err != nil {
		// Error executing query
		return p, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return p, err
	}

	p.ID = uint64(lastInsertID)
	return p, nil
}
