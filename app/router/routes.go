package router

// appendRoutes registers routes in the router
func (r *Router) appendRoutes() {
	// API
	api := r.e.Group("/api")
	api.Use(r.mwBasicAuth())

	// API, Version 1
	v1 := api.Group("/v1")

	// Endpoints for auth
	v1.GET("/login/get_customer_key", r.c.GetCustomerKey, r.mwKeyAuth("login"))
	v1.GET("/config/generate_customer_key", r.c.GenerateCustomerKey, r.mwKeyAuth("registration"))

	// Endpoints for customer
	v1.GET("/customer/get_customer", r.c.GetCustomer, r.mwKeyAuth("customer"), r.mwCUIDAuth("cust_unique_id"))
	v1.POST("/customer/customer_basic", r.c.PostAddCustomer, r.mwKeyAuth("customer"), r.mwCUIDAuth("cust_unique_id"))
	v1.GET("/customer/get_transactions", r.c.GetCustomerTransactions, r.mwKeyAuth("customer"), r.mwCUIDAuth("cust_unique_id"))
	v1.POST("/customer/add_transactions", r.c.PostAddCustomerTransactions, r.mwKeyAuth("customer"), r.mwCUIDAuth("cust_unique_id"))
	v1.GET("/customer/get_profile", r.c.GetCustomerProfile, r.mwKeyAuth("customer"), r.mwCUIDAuth("cust_unique_id"))
}
