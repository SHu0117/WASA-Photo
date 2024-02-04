package database

func (db *appdbimpl) SetUsername(u User, username string) error {
	res, err := db.c.Exec(`UPDATE user SET username = ? WHERE id=?`,
		username, u.ID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// If we didn't changed any row, then the user didn't exist
		return ErrDataDoesNotExist
	}
	return nil
}
