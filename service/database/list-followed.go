package database

func (db *appdbimpl) ListFollowed(user User) ([]User, error) {

	rows, err := db.c.Query("SELECT u.id, u.username FROM users u, following f WHERE f.Follower_id = ? and f.followed_id = u.id", user.ID)
	if err != nil {
		return nil, rows.Err()
	}
	// Wait for the function to finish before closing rows.
	defer func() { _ = rows.Close() }()

	// Read all the users in the resulset (users that follow the requesting user)
	var followedUsers []User
	for rows.Next() {
		var followed User
		err = rows.Scan(&followed.ID, &followed.Username)
		if err != nil {
			return nil, err
		}
		followedUsers = append(followedUsers, followed)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return followedUsers, nil
}
