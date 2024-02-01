package database

func (db *appdbimpl) ExistUsername(username string) error {
	var user User
	err := db.c.QueryRow(`SELECT u.id, u.username FROM user u WHERE u.username = ?`, username).Scan(&user.ID, &user.Username)
	if err != nil {
		return ErrDataDoesNotExist
	}
	return nil
}

func (db *appdbimpl) ExistUID(id uint64) error {
	var user User
	err := db.c.QueryRow(`SELECT u.id, u.username FROM user u WHERE u.id = ?`, id).Scan(&user.ID, &user.Username)
	if err != nil {
		return ErrDataDoesNotExist
	}
	return nil
}

func (db *appdbimpl) ExistPhoto(id uint64) error {
	var photo Photo
	err := db.c.QueryRow(`SELECT p.id FROM photo p WHERE p.id = ?`, id).Scan(&photo.ID)
	if err != nil {
		return ErrDataDoesNotExist
	}
	return nil
}
