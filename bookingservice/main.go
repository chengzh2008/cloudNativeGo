package main

import (
	"github.com/streadway/amqp"
	"../lib/configuration"
	"../lib/persistence/dblayer"
	"./listener"
	"./rest"
	"flag"
	msgqueue_amqp "../lib/msgqueue/amqp"
)

func main() {
	confPath := flag.String("config", "../lib/configuration/config.json", "path to config file")
	flag.Parse()
	config, _ := configuration.ExtractConfiguration(*confPath)

	conn, err := amqp.Dial(config.AMQPMessageBroker)
	if err != nil {
		panic(err)
	}
	eventListener, err := msgqueue_amqp.NewAMQPEventListener(conn, "events", "booking")
	if err != nil {
		panic(err)
	}

	eventEmitter, err := msgqueue_amqp.NewAMQPEventEmitter(conn, "events")
	if err != nil {
		panic(err)
	}

	dbhandler, err := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	if err != nil {
		panic(err)
	}

	processor := &listener.EventProcessor{eventListener, dbhandler}
	go processor.ProcessEvents()


	rest.ServeAPI(config.RestfulEndpoint, dbhandler, eventEmitter)


}
