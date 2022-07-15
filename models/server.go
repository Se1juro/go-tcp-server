package models

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type Server struct {
	channels []int
	clients  []Clients
}

// Accept new connection from client
func (s *Server) NewClient(listener net.Listener) {
	for {
		newClient, err := listener.Accept() // Accept and return connection from new client
		log.Println("New connection from", newClient.RemoteAddr())

		if err != nil {
			fmt.Println(err)
			continue
		}

		go s.handleClient(newClient)
	}
}

// Receive messages from client
func (s *Server) handleClient(client net.Conn) {
	for {
		var data Messages
		gob.NewDecoder(client).Decode(&data)

		s.RunCommand(data.Message, data.Args, client)
	}
}

func (s *Server) NewChannel(newChannel int) {
	s.channels = append(s.channels, newChannel)
}

func (s *Server) AddClientToChannel(client Clients) {
	s.clients = append(s.clients, client)
}

func (s *Server) GetChannels() []int {
	return s.channels
}
