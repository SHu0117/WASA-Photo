package database

func (db *appdbimpl) FollowUser(f Following) (Following, error) {
	_, err := db.c.Exec(`INSERT INTO following (Follower_id, followed_id) VALUES (?, ?)`,
		f.Follower_id, f.Followed_id)
	if err != nil {
		return f, err
	}

	return f, nil
}
