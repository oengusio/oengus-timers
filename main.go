package main

import (
	"log"
	"oengus-timers/rabbitmq"
)

func main() {
	log.Println("Connecting to RabbitMQ.....")
	rabbitmq.SetupAMQP()
	log.Println("Connection successful!")

	log.Println("Starting timers.....")
	StartTimers()
	log.Println("Timers have finished!")
}
