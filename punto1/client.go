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

	for {

		//Leer la ruta del archivo
		var filePath string
		fmt.Print("Client - Digite la ruta del archivo a enviar: ")
		fmt.Scanf("%s", &filePath)
		if filePath == ""{
			break
		}

		//Obtener el contenido del archivo
		fileContent, error := c.getFileContent(filePath)

		//Mnejo de error
		if error{
			fmt.Println("Client - Error leyendo el archivo")
			continue
		}

		//Envia un mensaje al servidor
		ws.WriteMessage(websocket.TextMessage, []byte(fileContent))

		//Obtener la respuesta del servidor
		_, msg, _ := ws.ReadMessage()
		fmt.Println("Client - respuesta: " + string(msg))
	}

	ws.Close()
	c.channel <- "ready"

}

func NewClient(addr string, channel chan string) *Client {

	return &Client{
		addr: flag.String("addr", addr, "http service address"),
		channel: channel,
	}
}

