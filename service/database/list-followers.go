package database

func (db *appdbimpl) ListFollowers(user User) ([]User, error) {

	rows, err := db.c.Query("SELECT u.id, u.username FROM user u, following f WHERE f.followed_id = ? and f.Follower_id = u.id", user.ID)
	if err != nil {
		return nil, err
	}
	// Wait for the function to finish before closing rows.
	defer func() { _ = rows.Close() }()

	// Read all the users in the resulset (users that follow the requesting user)
	var followers []User
	for rows.Next() {
		var follower User
		err = rows.Scan(&follower.ID, &follower.Username)
		if err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return followers, nil
}
