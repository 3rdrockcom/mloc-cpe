package router

func (r *Router) appendRoutes() {
	r.e.GET("/ping", r.c.Ping)
}
