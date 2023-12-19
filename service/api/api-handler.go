package api

import (
	"net/http"
)

// 
func (rt *_router) Handler() http.Handler {
	
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))


	rt.router.POST("/sessions", rt.wrap(rt.doLogin))

	rt.router.POST("/users/:uid/photo", rt.wrap(rt.uploadPhoto))
	
	rt.router.PUT("/users/:uid/following/:followeduid", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:uid/following/:followeduid", rt.wrap(rt.unfollowUser))

	rt.router.PUT("/users/:uid/banning/:banneduid", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:uid/banning/:banneduid", rt.wrap(rt.unbanUser))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
