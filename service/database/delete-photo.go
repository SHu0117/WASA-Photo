package database

func (db *appdbimpl) DeletePhoto(id uint64) error {

	_, err1 := db.c.Exec(`DELETE FROM likes WHERE photoId=?`, id)
	if err1 != nil {
		return err1
	}

	_, err2 := db.c.Exec(`DELETE FROM comments WHERE photoId=?`, id)
	if err2 != nil {
		return err2
	}

	res, err3 := db.c.Exec(`DELETE FROM photos WHERE id=?`, id)
	if err3 != nil {
		return err3
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
