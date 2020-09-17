package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
)

type Client struct {
	addr *string
	channel chan string
}

func (c *Client) showStatistics(statistics Statistics)  {

	fmt.Println("*****Palabras unicas*****")
	for _, word := range statistics.UniqueWordList{
		fmt.Println("	" + word)
	}

	fmt.Println("*****Contador de palabras*****")
	for word, quantity := range statistics.WordsCounter{
		fmt.Printf("	%q: %d\n", word, quantity)
	}

	fmt.Printf("* Size: %d\n", statistics.Size)

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

func (c *Client) fragmentContent(content string) []string {

	//Instancia las variables a usar
	var numParts int
	var charactersPerPart int
	firstIndex := 0

	//Calcula el numero de fragmentos en los que se dividira el mensaje
	numParts = len(content)/1000000

	//Si el numero de fragmentos es igual a 0 se entiende que solo sera un fragmento
	if numParts == 0 {
		return []string{content}
	}

	//Calcula el numero de caracteres por fragmento
	charactersPerPart = len(content)/numParts
	result := make([]string, numParts)

	//Empieza a armar el slice de fragmentos
	for i := 0; i < numParts; i++{
		if i == numParts -1{
			result[i] = content[firstIndex:]
		}else {
			result[i] = content[firstIndex:(firstIndex + charactersPerPart)]
			firstIndex += charactersPerPart
		}

	}

	return result

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

		//Manejo de error
		if error{
			fmt.Println("Client - Error leyendo el archivo")
			continue
		}

		contentParts := c.fragmentContent(fileContent)

		//Enviar al servidor el numero de fragmentos en los que se dividio el contenido del mensaje
		ws.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(len(contentParts))))

		//Envia un mensaje al servidor
		for i := 0; i < len(contentParts); i++ {
			ws.WriteMessage(websocket.TextMessage, []byte(contentParts[i]))
		}


		//Obtener la respuesta del servidor
		_, response, _ := ws.ReadMessage()

		var statistics Statistics
		err := json.Unmarshal(response, &statistics)

		if err != nil{
			log.Println(err)
		}

		c.showStatistics(statistics)

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

