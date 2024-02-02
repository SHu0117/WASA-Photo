package api

import (
	"net/http"
)

func (rt *_router) Handler() http.Handler {

	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/sessions", rt.wrap(rt.doLogin))
	rt.router.GET("/users/:username/homepage", rt.wrap(rt.getMyStream))

	rt.router.POST("/users/:username/photo/", rt.wrap(rt.uploadPhoto))
	rt.router.GET("/users/:username/photo/", rt.wrap(rt.getUserPhoto))
	rt.router.DELETE("/users/:username/photo/:pid", rt.wrap(rt.deletePhoto))
	rt.router.GET("/users/:username/photo/:pid", rt.wrap(rt.getPhoto))

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
