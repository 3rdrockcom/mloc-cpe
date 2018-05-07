package router

func (r *Router) appendRoutes() {
	r.e.GET("/ping", r.c.Ping)

	api := r.e.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/ping", r.c.Ping)

	v1.POST("/customers/add", r.c.PostAddCustomer)
	v1.POST("/customers/:customerID/transactions/add", r.c.PostAddCustomerTransactions)
	v1.GET("/customers/:customerID/profile", r.c.GetCustomerProfile)
}
