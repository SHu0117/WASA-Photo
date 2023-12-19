package database

func (db *appdbimpl) CheckBanned(user User, banned User) (User, error) {
	var res User
	err := db.c.QueryRow(`SELECT u.id, u.username FROM user u, banning b WHERE b.banner_id = ? and b.banned_id = ? AND u.id = ?`,user.ID, banned.ID, banned.ID).Scan(&res.ID, &res.Username)
	if err != nil {
		return user, err
	}
	return res, nil
}