package main

import (
	"log"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/controllers"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/database"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/router"

	"github.com/namsral/flag"
)

var dsn string

func init() {
	flag.StringVar(&dsn, "dsn", "", "path to sample data i.e. user:pass@/example")

	flag.Parse()
}

func main() {
	var err error

	db := database.NewDatabase("mysql", dsn)
	err = db.Open()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	c := controllers.NewControllers(db)

	r := router.NewRouter(c)
	err = r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
