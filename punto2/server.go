package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type Server struct {
	port string
	upgrader websocket.Upgrader
	channel chan string
}

func (s *Server) formatCommand(msg []byte) (string, []string) {

	//Elimina los saltos de linea
	command := strings.Replace(string(msg), "\n", "", 1)
	//Separa el comando estre sus espacios
	commandParts := strings.Split(command, " ")

	//Agrupar y separar el comando y sus argumentos
	for i := 2; i < len(commandParts); i++{
		commandParts[1] = commandParts[1] + " " + commandParts[i]
	}

	return command, commandParts

}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {

	//Instancia la conexión
	ws, err := s.upgrader.Upgrade(w, r, nil)

	//Manejo de error
	if err != nil{
		log.Print(err)
		return
	}

	for {

		//Obtener el comando requerido
		_, msg, err := ws.ReadMessage()

		//Menejo del error
		if err != nil {
			log.Println(err)
			break
		}

		//Llamar a la función que le da formato al comando
		command, commandParts := s.formatCommand(msg)

		//Definir si el comando se ejecyta con alguna opción
		var bash *exec.Cmd
		if len(commandParts) == 1 {
			bash = exec.Command(command)
		}else {
			bash = exec.Command(commandParts[0], commandParts[1])
		}

		//Atrapa la respuesta del comando
		out, err := bash.CombinedOutput()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}

		//Enviar la salida del comando al cliente
		ws.WriteMessage(websocket.TextMessage, []byte(string(out)))
	}

	ws.Close()

}

func (s *Server) runServer()  {

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