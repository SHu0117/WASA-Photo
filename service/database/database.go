/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrDataDoesNotExist = errors.New("data does not exist")

type User struct {
	ID        uint64
	Username  string
}

type Photo struct {
	ID           uint64
	User_id		 uint64
	N_likes	     int64
	N_comments   int64
	Upload_time  time.Time
}

type Following struct {
	Follower_id uint64
	Followed_id uint64
}

type Banning struct {
	Banner_id uint64
	Banned_id uint64
}



// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	NewUser(User) (User, error)
	UploadPhoto(Photo) (Photo, error)
	//GetMyStream() ([]Photo, error)
	SetUsername(user User, username string) error
	FollowUser(Following) (Following, error)
	UnfollowUser(Following) error
	BanUser(Banning) (Banning, error)
	UnbanUser(Banning) error
	ExistUsername(username string) (error)
	ExistUID(id uint64) (error)
	ListFollowers(u User)([]User, error)
	ListFollowed(u User)([]User, error)
	ListBanned(u User)([]User, error)
	CheckBanned(user User, banned User) (User, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE users(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			username TEXT NOT NULL,
			UNIQUE(username));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='following';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE following(
			follower_id INTEGER,
			followed_id INTEGER,
			FOREIGN KEY(follower_id) REFERENCES users(id),
			FOREIGN KEY(followed_id) REFERENCES users(id));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='banning';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE banning(
			banner_id INTEGER,
			banned_id INTEGER,
			FOREIGN KEY(banner_id) REFERENCES users(id),
			FOREIGN KEY(banned_id) REFERENCES users(id));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='photo';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE photo(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			user_id INTEGER NOT NULL,
			n_likes INTEGER NOT NULL,
			n_comments INTEGER NOT NULL,
			upload_time DATETIME,
			FOREIGN KEY(user_id) REFERENCES users(id));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='like';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE like(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			user_id INTEGER NOT NULL,
			photo_id INTEGER NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(photo_id) REFERENCES photo(id));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='comment';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE comment(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			texts TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			photo_id INTEGER NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(photo_id) REFERENCES photo(id));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
