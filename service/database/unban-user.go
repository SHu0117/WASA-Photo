package database

func (db *appdbimpl) UnbanUser(b Banning) error {
	res, err := db.c.Exec(`DELETE FROM banning WHERE Banner_id=? AND Banned_id=?`, b.Banner_id, b.Banned_id)
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
