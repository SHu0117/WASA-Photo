package api

import (
	"github.com/SHu0117/WASA-Photo/service/database"
)

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
}

func (u *User) UserFromDatabase(user database.User) {
	u.ID = user.ID
	u.Username = user.Username
}

func (u *User) UserToDatabase() database.User {
	return database.User{
		ID:       u.ID,
		Username: u.Username,
	}
}

type Following struct {
	Follower_id uint64 `json:"Follower_id"`
	Followed_id uint64 `json:"followed_id"`
}

func (f *Following) FollowingFromDatabase(following database.Following) {
	f.Follower_id = following.Follower_id
	f.Followed_id = following.Followed_id
}

func (following *Following) FollowingToDatabase() database.Following {
	return database.Following{
		Follower_id: following.Follower_id,
		Followed_id: following.Followed_id,
	}
}

type Banning struct {
	Banner_id uint64 `json:"Banner_id"`
	Banned_id uint64 `json:"Banned_id"`
}

func (b *Banning) BanningFromDatabase(banning database.Banning) {
	b.Banner_id = banning.Banner_id
	b.Banned_id = banning.Banned_id
}

func (banning *Banning) BanningToDatabase() database.Banning {
	return database.Banning{
		Banner_id: banning.Banner_id,
		Banned_id: banning.Banned_id,
	}
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

func (p *Photo) PhotoFromDatabase(photo database.Photo) {
	p.ID = photo.ID
	p.User_id = photo.User_id
	p.User_username = photo.User_username
	p.N_likes = photo.N_likes
	p.N_comments = photo.N_comments
	p.Upload_time = photo.Upload_time
	p.File = photo.File
	p.IsLiked = photo.IsLiked
}

func (p *Photo) PhotoToDatabase() database.Photo {
	return database.Photo{
		ID:            p.ID,
		User_id:       p.User_id,
		User_username: p.User_username,
		N_likes:       p.N_likes,
		N_comments:    p.N_comments,
		Upload_time:   p.Upload_time,
		File:          p.File,
		IsLiked:       p.IsLiked,
	}
}

type Like struct {
	ID         uint64 `json:"id"`
	User_id    uint64 `json:"userId"`
	Photo_id   uint64 `json:"photoId"`
	Photo_user uint64 `json:"photoOwner"`
}

func (l *Like) LikeFromDatabase(like database.Like) {
	l.ID = like.ID
	l.User_id = like.User_id
	l.Photo_id = like.Photo_id
	l.Photo_user = like.Photo_user

}

func (l *Like) LikeToDatabase() database.Like {
	return database.Like{
		ID:         l.ID,
		User_id:    l.User_id,
		Photo_id:   l.Photo_id,
		Photo_user: l.Photo_user,
	}
}

type Comment struct {
	ID            uint64 `json:"id"`
	User_id       uint64 `json:"user_id"`
	User_username string `json:"username"`
	Photo_id      uint64 `json:"photo_id"`
	Photo_user    uint64 `json:"photo_Owner"`
	Text          string `json:"text"`
}

func (c *Comment) CommentFromDatabase(comment database.Comment) {
	c.ID = comment.ID
	c.User_id = comment.User_id
	c.User_username = comment.User_username
	c.Photo_id = comment.Photo_id
	c.Photo_user = comment.Photo_user
	c.Text = comment.Text
}

func (c *Comment) CommentToDatabase() database.Comment {
	return database.Comment{
		ID:            c.ID,
		User_id:       c.User_id,
		User_username: c.User_username,
		Photo_id:      c.Photo_id,
		Photo_user:    c.Photo_user,
		Text:          c.Text,
	}
}

type Profile struct {
	Requester_id uint64 `json:"requester_id"`
	User_id      uint64 `json:"user_id"`
	Username     string `json:"username"`
	N_followers  int    `json:"followers"`
	N_followed   int    `json:"followed"`
	N_photo      int    `json:"photos"`
	Is_followed  bool   `json:"isFollowed"`
	Is_banned    bool   `json:"isBanned"`
}
