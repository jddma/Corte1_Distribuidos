package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/url"
)

type Client struct {
	addr *string
	channel chan string
}

func (c *Client) getFileContent(path string) (string, bool) {

	var result string

	file, err := ioutil.ReadFile(path)
	if err != nil{
		return result, true
	}

	result = string(file)

	return result, false

}

func (c *Client) runClient()  {


	//Leer la ruta del archivo
	var filePath string
	fmt.Print("Client - Digite la ruta del archivo a enviar: ")
	fmt.Scanf("%s", &filePath)

	//Obtener el contenido del archivo
	fileContent, error := c.getFileContent(filePath)

	//Manejo del error
	if error{
		c.channel <- "Error al leer el archivo"
	}

	//Preparar la serialización de la URL
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *c.addr, Path: "/"}
	fmt.Println("Cliente conectandose a " + u.String())

	//Establecer la conexión con el servidor
	s, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	//Manejo de errores
	if err != nil{
		log.Println(err)
		return
	}

	//Envia un mensaje al servidor
	s.WriteMessage(websocket.TextMessage, []byte(fileContent))

	//Obtener la respuesta del servidor
	_, msg, err := s.ReadMessage()
	fmt.Println("CLiente - " + string(msg))

	s.Close()
	c.channel <- "ready"

}

func NewClient(addr string, channel chan string) *Client {

	return &Client{
		addr: flag.String("addr", addr, "http service address"),
		channel: channel,
	}
}

