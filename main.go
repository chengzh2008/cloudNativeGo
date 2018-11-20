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
	confPath := flag.String("conf", `./eventsservice/lib/configuration/config.json`, "flag to set the path to the configuration json file")
	flag.Parse()
	config, _ := configuration.ExtractConfiguration(*confPath)
	fmt.Println("Connecting to database")
	dbhandler, err := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	if err != nil {
		fmt.Println("something erro happened during data connection")
		return
	}
	httpErrChan, httpTLSErrChan := rest.ServeAPI(config.RestfulEndpoint, config.RestfulTLSEndpoint, dbhandler)
	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP error: ", err)
	case err := <- httpTLSErrChan:
		log.Fatal("HTTPS error: ", err)
	}
}
