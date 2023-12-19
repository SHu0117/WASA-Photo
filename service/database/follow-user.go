package database

func (db *appdbimpl) FollowUser(f Following) (Following, error) {
	res, err := db.c.Exec(`INSERT INTO following (follower_id, followed_id) VALUES (?, ?)`,
		f.Follower_id, f.Followed_id)
	if err != nil {
		return f, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return f, err
	}

	f.Follower_id = uint64(lastInsertID)
	return f, nil
}