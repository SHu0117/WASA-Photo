package api

import(
	"time"
	"WasaPhoto-1985972/service/database"
)

type User struct {
	ID        uint64  `json:"id"`
	Username  string  `json:"username"`
}

func (u *User) UserFromDatabase(user database.User) {
	u.ID = user.ID
	u.Username = user.Username
}

func (u *User) UserToDatabase() database.User {
	return database.User{
		ID:        u.ID,
		Username:  u.Username,
	}
}


type Following struct {
	Follower_id    uint64  `json:"follower_id"`
	Followed_id    uint64 `json:"followed_id"`
}

func (f *Following) FollowingFromDatabase(following database.Following) {
	f.Follower_id = following.Follower_id
	f.Followed_id = following.Followed_id
}

func (following *Following) FollowingToDatabase() database.Following {
	return database.Following{
		Follower_id:        following.Follower_id,
		Followed_id:  		following.Followed_id,
	}
}

type Banning struct {
	Banner_id    uint64  `json:"banner_id"`
	Banned_id    uint64  `json:"banned_id"`
}

func (b *Banning) BanningFromDatabase(banning database.Banning) {
	b.Banner_id = banning.Banner_id
	b.Banned_id = banning.Banned_id
}

func (banning *Banning) BanningToDatabase() database.Banning {
	return database.Banning{
		Banner_id:          banning.Banner_id,
		Banned_id:  		banning.Banned_id,
	}
}

type Photo struct{
	ID 			uint64			
	User_id 	uint64			
	N_likes 	int64
	N_comments 	int64
	Upload_time	time.Time
}

func (p *Photo) PhotoFromDatabase(photo database.Photo){
	p.ID = photo.id,				
	p.User_id = photo.User_id,
	p.N_likes = photo.N_likes,
	p.N_comments = photo.N_comments,
	p.Upload_time = phot.Upload_time,
}


func (p *Photo) PhotoToDatabase() database.Photo{
	return database.Photo{
		p.ID :			photo.ID,
		p.User_id :		photo.User_id,
		p.N_likes 	: 	photo.N_likes,
		p.N_comments :	photo.N_comments,
		p.Upload_time : photo.Upload_time,
	}
}


