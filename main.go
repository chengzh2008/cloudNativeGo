package main

import (
	"flag"
	"fmt"
	"log"
	"./eventsservice/rest"
	"./eventsservice/lib/configuration"
	"./eventsservice/lib/persistence/dblayer"
)

func main() {
	confPath := flag.String("conf", `./lib/configuration/config.json`, "flag to set the path to the configuration json file")
	flag.Parse()
	config, _ := configuration.ExtractConfiguration(*confPath)
	fmt.Println("Connecting to database")
	dbhandler, err := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	if err != nil {
		fmt.Println("something erro happened during data connection")
		return
	}
	log.Fatal(rest.ServeAPI(config.RestfulEndpoint, dbhandler))
}
