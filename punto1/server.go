package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	port string
	upgrader websocket.Upgrader
	channel chan string
}

func (s *Server) getData(requestBody string) string {

	text := strings.Replace(requestBody, "\n", "", -1)

	//Separar cada una de las palabras en un slice
	words := strings.Split(text, " ")
	//Crear un map para hacer el contador de palabras donde la clave es la palabra y el valor las veces que se repite
	counter := make(map[string]int)

	//Itera las palabras obtenidas y cuenta cada una de ellas asignandole una clave en el contador
	for _, word := range words{
		counter[word]++
	}

	//Crear un slice para alacenar las palabras unicas
	uniqueWordList := []string{}

	//Itera el contador de lapabras y agrega las palabras que solo se encontraron una vez al slice
	for word, quantity := range counter{
		if quantity == 1 {
			uniqueWordList = append(uniqueWordList, word)
		}
	}

	//Serializar las estadisticas obtenidas en JSON
	data := NewStatistics(counter, uniqueWordList, len(requestBody))
	decode, _ := json.Marshal(data)
	statistics := string(decode)

	return statistics

}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Server - Nueva conexión establecida con cliente")
	//Instancia la conexión
	ws, err := s.upgrader.Upgrade(w, r, nil)

	//Manejo de error
	if err != nil{
		log.Print(err)
		return
	}

	for {

		//Recibe el numero de fragmentos en los que llegara el contenido
		_, strNumFragments, err := ws.ReadMessage()

		//Manejo del error
		if err != nil {
			log.Println(err)
			break
		}

		//Inicia las variables de control
		numFragments, _ := strconv.Atoi(string(strNumFragments))
		content := ""

		//Recibe los fragmentos del contenido y los concatena
		for i := 0; i < numFragments; i++{
			//Lee el mensaje enviado por el socket
			_, msg, err := ws.ReadMessage()

			//Manejo de error
			if err != nil {
				log.Println(err)
				break
			}

			content += string(msg)
		}

		fmt.Println("Server - mensaje recibido")

		//Envio de la respuesta
		ws.WriteMessage(websocket.TextMessage, []byte(s.getData(content)))
	}
	ws.Close()

}

func (s *Server) runServer()  {

	//Definir la el path donde llegaran las solicitudes
	http.HandleFunc("/", s.handleConnections)

	//Envia al channel la señal de que el servidor ya esta listo
	s.channel <- "Server - Servidor escuchando el puerto " + s.port

	//Inicializa en servidor
	err := http.ListenAndServe(s.port, nil)


	if err != nil {
		log.Println(err)
	}

}

func NewServer(port string, channel chan string) *Server {

	return &Server{
		port: port,
		upgrader: websocket.Upgrader{},
		channel: channel,
	}

}

