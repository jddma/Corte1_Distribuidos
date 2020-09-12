package main

import (
	"fmt"
	"log"
)

func main()  {

	channel :=  make(chan string)

	server := NewServer(":15432", channel)
	go server.runServer()

	fmt.Println(<- channel)

	client := NewClient("127.0.0.1:15432", channel)
	go client.runClient()

	log.Println(<- channel)

}