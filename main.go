package main

import "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/router"

func main() {
	r := router.NewRouter()
	r.Run()
}
