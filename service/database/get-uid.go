package database

func (db *appdbimpl) GetUserID(username string) (User, error) {
	var user User
	err := db.c.QueryRow(`SELECT u.id, u.username FROM user u WHERE u.username = ?`, username).Scan(&user.ID, &user.Username)
	if err != nil {
		return user, ErrDataDoesNotExist
	}
	return user, nil
}
