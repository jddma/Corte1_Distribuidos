package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

type Client struct {
	addr *string
}

func (c *Client) runClient()  {

	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *c.addr, Path: "/"}
	log.Printf("Cliente conectandose a %s", u.String())

	s, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil{
		log.Println(err)
		return
	}

	s.WriteMessage(websocket.TextMessage, []byte("Esto es un mensaje"))

	s.Close()

}

func NewClient(addr string) *Client {

	return &Client{
		addr: flag.String("addr", addr, "http service address"),
	}
}

