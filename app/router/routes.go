package router

func (r *Router) appendRoutes() {
	api := r.e.Group("/api")
	api.Use(r.mwBasicAuth())

	v1 := api.Group("/v1")

	v1.GET("/login/get_customer_key", r.c.GetCustomerKey, r.mwKeyAuth("login"))
	v1.GET("/config/generate_customer_key", r.c.GenerateCustomerKey, r.mwKeyAuth("registration"))

	v1.GET("/customer/get_customer", r.c.GetCustomer, r.mwKeyAuth("default"))
	v1.POST("/customer/customer_basic", r.c.PostAddCustomer, r.mwKeyAuth("default"))
	v1.POST("/customer/add_transactions", r.c.PostAddCustomerTransactions, r.mwKeyAuth("default"))
	v1.GET("/customer/get_profile", r.c.GetCustomerProfile, r.mwKeyAuth("default"))
}
