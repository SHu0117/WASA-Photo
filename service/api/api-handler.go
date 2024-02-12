package api

import (
	"net/http"
)

func (rt *_router) Handler() http.Handler {

	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.GET("/users/:username/homepage", rt.wrap(rt.getMyStream))
	rt.router.PUT("/users/:username", rt.wrap(rt.setMyUserName))
	rt.router.GET("/users/:username/profile", rt.wrap(rt.getUserProfile))

	rt.router.POST("/users/:username/photos/", rt.wrap(rt.uploadPhoto))
	rt.router.GET("/users/:username/photos/", rt.wrap(rt.getUserPhoto))
	rt.router.DELETE("/users/:username/photos/:pid", rt.wrap(rt.deletePhoto))
	// rt.router.GET("/users/:username/photos/:pid", rt.wrap(rt.getPhoto))

	rt.router.PUT("/users/:username/photos/:pid/likes/:likeUsername", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/users/:username/photos/:pid/likes/:likeUsername", rt.wrap(rt.unlikePhoto))
	rt.router.GET("/users/:username/photos/:pid/likes", rt.wrap(rt.getPhotoLikes))

	rt.router.POST("/users/:username/photos/:pid/comments/", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/users/:username/photos/:pid/comments/:cid", rt.wrap(rt.uncommentPhoto))
	rt.router.GET("/users/:username/photos/:pid/comments/", rt.wrap(rt.getPhotoComments))

	rt.router.PUT("/users/:username/following/:followUsername", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:username/following/:followUsername", rt.wrap(rt.unfollowUser))
	rt.router.GET("/users/:username/following/", rt.wrap(rt.listFollowed))

	rt.router.PUT("/users/:username/banned/:banUsername", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:username/banned/:banUsername", rt.wrap(rt.unbanUser))
	rt.router.GET("/users/:username/banned/", rt.wrap(rt.listBanned))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
