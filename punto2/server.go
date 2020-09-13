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

func (s *Server) formatCommand(msg []byte) (string, string, bool) {

	//Elimina los saltos del comando enviado
	strMsg := strings.Replace(string(msg), "\n", "", 1)

	//Variable para almacenar unicamente el comando
	command := strMsg
	//Variable para almacenar las variables del comando
	var options string
	//Vaiabe que sirve de bandera para indicar si la entrada es un comando con opciones adicionales
	useOptions := false

	//Itera entre la entrada enviada por el cliente
	for i := 0; i < len(strMsg); i++ {

		//En caso de encontrar un espacio separa el comando de sus opciones
		if strMsg[i] == ' '{
			letters := []rune(strMsg)
			command = string(letters[:i])
			options = string(letters[(i +1):])
			useOptions = true
			break
		}
	}

	return command, options, useOptions

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
		command, options, useOptions := s.formatCommand(msg)

		//Definir si el comando se ejecyta con alguna opción
		var bash *exec.Cmd
		if useOptions {
			bash = exec.Command(command, options)
		}else {
			bash = exec.Command(command)
		}

		//Atrapa la respuesta del comando
		out, err := bash.CombinedOutput()
		if err != nil {
			log.Println("cmd.Run() failed with %s\n", err)
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