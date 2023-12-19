package database

func (db *appdbimpl) SetUsername(u User, username string) error {
	res, err := db.c.Exec(`UPDATE users SET username = ? WHERE id=?`,
		username, u.ID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// If we didn't delete any row, then the user didn't exist
		return ErrDataDoesNotExist
	}
	return nil
}