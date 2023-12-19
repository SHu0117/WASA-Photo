package database

func (db *appdbimpl) UnfollowUser(f Following) error {
	res, err := db.c.Exec(`DELETE FROM following WHERE follower_id=? AND followed_id=?`, f.Follower_id, f.Followed_id)
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