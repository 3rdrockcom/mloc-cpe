package router

func (r *Router) appendRoutes() {
	r.e.GET("/ping", r.c.Ping)

	r.e.POST("/customer/add", r.c.PostAddCustomer)
	r.e.POST("/customers/:customerID/transactions/add", r.c.PostAddCustomerTransactions)
}
