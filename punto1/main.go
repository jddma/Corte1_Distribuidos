package main

import (
	"log"
	"time"
)

func main()  {

	var channel chan string

	server := NewServer(":6789")
	go server.runServer()

	time.Sleep(2 * time.Second)

	client := NewClient("127.0.0.1:6789")
	go client.runClient()

	log.Fatal(<- channel)

}
