package main

import (
	"log"

	"github.com/epointpayment/mloc-cpe/app/controllers"
	"github.com/epointpayment/mloc-cpe/app/database"
	"github.com/epointpayment/mloc-cpe/app/router"
	"github.com/epointpayment/mloc-cpe/app/services"

	"github.com/namsral/flag"
)

var dsn string

func init() {
	flag.StringVar(&dsn, "dsn", "", "path to sample data i.e. user:pass@/example")

	flag.Parse()
}

func main() {
	var err error

	// Create new connection handler for database
	db := database.NewDatabase("mysql", dsn)
	err = db.Open()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Setup services
	err = services.New(db)
	if err != nil {
		log.Fatalln(err)
	}

	// Setup router and run
	c := controllers.NewControllers()
	r := router.NewRouter(c)
	err = r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
