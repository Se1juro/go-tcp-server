package models

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type Server struct {
	channels    []int
	Clients     []Clients
	FileHistory []FileHistory
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
		if data.Message == "Disconnect" {
			fmt.Printf("Client %s disconnected\n", client.RemoteAddr().String())
			s.RemoveClientFromServer(client)
			client.Close()
		}
		s.RunCommand(data.Message, data.Args, client)
	}
}

func (s *Server) NewChannel(newChannel int) {
	s.channels = append(s.channels, newChannel)
}

func (s *Server) AddClientToChannel(client Clients) {
	s.Clients = append(s.Clients, client)
}

func (s *Server) RemoveClientFromServer(conn net.Conn) {
	for index, client := range s.Clients {
		if client.Conn.RemoteAddr() == conn.RemoteAddr() {
			s.Clients = append(s.Clients[:index], s.Clients[index+1:]...)
			fmt.Println(s.Clients)
		}
	}
}

func (s *Server) GetChannels() []Channels {
	var channels []Channels
	for index, channel := range s.channels {
		channels = append(channels, Channels{ChannelId: channel})
		for _, client := range s.Clients {
			if client.CurrentChannel == channel {
				channels[index].Clients = append(channels[index].Clients, ChannelClients{Conn: client.Conn.RemoteAddr().String(), CurrentChannel: client.CurrentChannel, Status: client.Status})
			}
		}
	}
	return channels
}

func (s *Server) GetAllHistoryFiles() []FileHistory {
	return s.FileHistory
}

func (s *Server) GetAllFilesFromChannel(channelId int) []FileHistory {
	var newHistoryFiles []FileHistory
	for _, channel := range s.FileHistory {
		if channel.Channel == channelId {
			newHistoryFiles = append(newHistoryFiles, channel)
		}
	}
	return newHistoryFiles
}
