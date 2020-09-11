package main

import (
	"fmt"
	"log"
)

func main()  {

	channel :=  make(chan string)

	server := NewServer(":6789", channel)
	go server.runServer()

	fmt.Println(<- channel)

	client := NewClient("127.0.0.1:6789", channel)
	go client.runClient()

	log.Fatal(<- channel)

}
