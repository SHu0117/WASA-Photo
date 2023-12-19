package database

func (db *appdbimpl) ExistUsername(username string) (User, error) {
	var user User
	err := db.c.QueryRow(`SELECT u.id, u.username FROM user u WHERE u.username = ?`,username).Scan(&user.ID, &user.Username)
	if err != nil {
		return user, ErrDataDoesNotExist
	}
	return user, nil
}

func (db *appdbimpl) ExistUID(id uint64) (User, error) {
	var user User
	err := db.c.QueryRow(`SELECT u.id, u.username FROM user u WHERE u.id = ?`,id).Scan(&user.ID, &user.Username)
	if err != nil {
		return user, ErrDataDoesNotExist
	}
	return user, nil
}