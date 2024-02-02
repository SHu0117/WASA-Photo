package database

func (db *appdbimpl) BanUser(b Banning) (Banning, error) {
	res, err := db.c.Exec(`INSERT INTO banning (Banner_id, Banned_id) VALUES (?, ?)`,
		b.Banner_id, b.Banned_id)
	if err != nil {
		return b, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return b, err
	}

	b.Banner_id = uint64(lastInsertID)
	return b, nil
}
