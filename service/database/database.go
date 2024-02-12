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
)

var ErrDataDoesNotExist = errors.New("data does not exist")
var ErrUserHasBeenBanned = errors.New("you are banned by the user, can't get any information")

type User struct {
	ID       uint64
	Username string
}

type Photo struct {
	ID            uint64 `json:"id"`
	User_id       uint64 `json:"userId"`
	User_username string `json:"username"`
	N_likes       int64  `json:"likesN"`
	N_comments    int64  `json:"commentsN"`
	Upload_time   string `json:"uploadtime"`
	File          string `json:"file"`
	IsLiked       bool   `json:"isliked"`
}

type Following struct {
	Follower_id uint64
	Followed_id uint64
}

type Banning struct {
	Banner_id uint64
	Banned_id uint64
}

type Like struct {
	ID         uint64 `json:"id"`
	User_id    uint64 `json:"userId"`
	Photo_id   uint64 `json:"photoId"`
	Photo_user uint64 `json:"photoOwner"`
}

type Comment struct {
	ID            uint64 `json:"id"`
	User_id       uint64 `json:"user_id"`
	User_username string `json:"username"`
	Photo_id      uint64 `json:"photo_id"`
	Photo_user    uint64 `json:"photo_Owner"`
	Text          string `json:"text"`
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	NewUser(User) (User, error)
	UploadPhoto(Photo) (Photo, error)
	SetUsername(user User, username string) error
	FollowUser(Following) (Following, error)
	UnfollowUser(Following) error
	BanUser(Banning) (Banning, error)
	UnbanUser(Banning) error
	ExistUsername(username string) error
	ExistUID(id uint64) error
	ExistPhoto(id uint64) error
	ListFollowers(u User) ([]User, error)
	ListFollowed(u User) ([]User, error)
	ListBanned(u User) ([]User, error)
	CheckBeingBanned(user User, requesterID uint64) (bool, error)
	DeletePhoto(id uint64) error
	GetUserID(username string) (User, error)
	GetUsername(id uint64) (User, error)
	GetMyStream(u User) ([]Photo, error)
	GetUserPhotos(u User, requesterID uint64) ([]Photo, error)
	// GetPhoto(pid uint64) (Photo, error)
	LikePhoto(l Like) (Like, error)
	UnlikePhoto(pid uint64, uid uint64) error
	ListLikes(pid uint64) ([]User, error)
	CommentPhoto(c Comment) (Comment, error)
	UncommentPhoto(pid uint64, uid uint64, cid uint64) error
	ListComment(pid uint64) ([]Comment, error)
	CountFollowed(user User) (int, error)
	CountPhotos(u User) (int, error)
	CountFollower(user User) (int, error)
	CheckIfFollowed(targetID uint64, requesterID uint64) (bool, error)
	CheckIfBanned(targetID uint64, requesterID uint64) (bool, error)
	CheckIfLiked(targetID uint64, requesterID uint64) (bool, error)
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

	// Enable foreign keys
	_, err := db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='user';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {

		userTable := `CREATE TABLE IF NOT EXISTS user(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			username TEXT NOT NULL UNIQUE);`

		followingTable := `CREATE TABLE IF NOT EXISTS following(
			Follower_id INTEGER,
			followed_id INTEGER,
			FOREIGN KEY(Follower_id) REFERENCES user(id),
			FOREIGN KEY(followed_id) REFERENCES user(id));`

		banningTable := `CREATE TABLE IF NOT EXISTS banning(
			Banner_id INTEGER,
			Banned_id INTEGER,
			FOREIGN KEY(Banner_id) REFERENCES user(id),
			FOREIGN KEY(Banned_id) REFERENCES user(id));`

		photoTable := `CREATE TABLE IF NOT EXISTS photos(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			user_id INTEGER NOT NULL,
			file TEXT NOT NULL,
			upload_time TEXT,
			FOREIGN KEY(user_id) REFERENCES user(id));`

		likeTable := `CREATE TABLE IF NOT EXISTS like(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			user_id INTEGER NOT NULL,
			photo_id INTEGER NOT NULL,
			FOREIGN KEY(user_id) REFERENCES user(id),
			FOREIGN KEY(photo_id) REFERENCES photos(id));`

		commentTable := `CREATE TABLE IF NOT EXISTS comment(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			texts TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			photo_id INTEGER NOT NULL,
			FOREIGN KEY(user_id) REFERENCES user(id),
			FOREIGN KEY(photo_id) REFERENCES photos(id));`

		_, err = db.Exec(userTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(photoTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(likeTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(commentTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(banningTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(followingTable)
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
