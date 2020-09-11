package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Server struct {
	port string
	upgrader websocket.Upgrader
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {

	ws, err := s.upgrader.Upgrade(w, r, nil)

	if err != nil{
		log.Print(err)
		return
	}

	_, msg, err := ws.ReadMessage()

	if err != nil{
		log.Println(err)
	}

	fmt.Println("Server-" + string(msg))

	ws.Close()

}

func (s *Server) runServer()  {

	http.HandleFunc("/", s.handleConnections)

	fmt.Println("Servidor escuchando el puerto " + s.port)
	err := http.ListenAndServe(s.port, nil)
	if err != nil {
		log.Println(err)
	}

}

func NewServer(port string) *Server {

	return &Server{
		port: port,
		upgrader: websocket.Upgrader{},
	}

}

