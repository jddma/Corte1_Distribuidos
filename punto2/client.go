package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
)

type Client struct {
	addr *string
	channel chan string
}

func (c *Client) runClient()  {

	//Preparar la serialización de la URL
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *c.addr, Path: "/"}
	fmt.Println("Client -  estableciendo conexión con " + u.String())

	//Establecer la conexión con el servidor
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	//Manejo de errores
	if err != nil{
		log.Println(err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {

		//Leer el comando requerido
		fmt.Print(">")
		command, _ := reader.ReadString('\n')

		//En caso de que el usuario digite exit terminara el programa
		if command == "exit\n"{
			c.channel <- "ready"
		}

		ws.WriteMessage(websocket.TextMessage, []byte(command))

		_, response, _ := ws.ReadMessage()
		fmt.Println(string(response))
	}

}

func NewClient(addr string, channel chan string) *Client {

	return &Client{
		addr: flag.String("addr", addr, "http service address"),
		channel: channel,
	}
}