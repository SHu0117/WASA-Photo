package database

func (db *appdbimpl) NewUser(u User) (User, error) {
	res, err := db.c.Exec(`INSERT INTO user (id, username) VALUES (NULL, ?)`,
		u.Username)
	if err != nil {
		return u, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return u, err
	}

	u.ID = uint64(lastInsertID)
	return u, nil
}