package router

func (r *Router) appendRoutes() {
	r.e.GET("/ping", r.c.Ping)

	api := r.e.Group("/api")
	api.Use(r.mwBasicAuth())

	v1 := api.Group("/v1")

	v1.GET("/login/get_customer_key", r.c.GetCustomerKey, r.mwKeyAuth("login"))
	v1.GET("/config/generate_customer_key", r.c.GenerateCustomerKey, r.mwKeyAuth("registration"))

	v1.POST("/customers/add", r.c.PostAddCustomer, r.mwKeyAuth("default"))
	v1.POST("/customers/:customerID/transactions/add", r.c.PostAddCustomerTransactions, r.mwKeyAuth("default"))
	v1.GET("/customers/:customerID/profile", r.c.GetCustomerProfile, r.mwKeyAuth("default"))
}
