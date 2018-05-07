package router

func (r *Router) appendRoutes() {
	r.e.GET("/ping", r.c.Ping)

	r.e.POST("/customers/add", r.c.PostAddCustomer)
	r.e.POST("/customers/:customerID/transactions/add", r.c.PostAddCustomerTransactions)
	r.e.GET("/customers/:customerID/profile", r.c.GetCustomerProfile)
}
