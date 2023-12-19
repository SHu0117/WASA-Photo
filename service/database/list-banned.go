package database

func (db *appdbimpl) ListBanned(user User) ([]User, error) {

	rows, err := db.c.Query("SELECT u.id, u.username FROM users u, banning b WHERE b.banner_id = ? and b.banned_id = u.id", user.ID)
	if err != nil {
		return nil, err
	}
	// Wait for the function to finish before closing rows.
	defer func() { _ = rows.Close() }()

	// Read all the users in the resulset (users that follow the requesting user)
	var bannedUsers []User
	for rows.Next() {
		var banned User
		err = rows.Scan(&banned.ID, &banned.Username)
		if err != nil {
			return nil, err
		}
		bannedUsers = append(bannedUsers, banned)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return bannedUsers, nil
}