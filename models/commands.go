package models

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/upload-files-go/exceptions"
)

func (s *Server) RunCommand(data string, args []byte, client net.Conn) {
	if data == "" {
		return
	}
	if strings.HasPrefix(data, "subscribe channel") {
		s.SubscribeClient(data, client)
	}

	if strings.HasPrefix(data, "send") {
		commands := strings.Split(data, " ")
		fileName := commands[1]

		s.SendData(Messages{Message: data + " " + fileName, Args: args}, client)
	}

	if strings.HasPrefix(data, "change status") {
		args := strings.Split(data, "change status")
		status := args[1]

		s.ChangeStatusClient(client, status)
	}
}

func (s *Server) SubscribeClient(data string, client net.Conn) {
	commands := strings.Split(data, " ")
	channelId, err := strconv.Atoi(commands[2])

	for _, c := range s.Clients {
		if c.Conn.RemoteAddr() == client.RemoteAddr() {
			err := gob.NewEncoder(client).Encode(&Messages{Message: "The client is already logged in"})
			exceptions.ManageError(err, "The client is already logged in")
			return
		}
	}

	exceptions.ManageError(err, "The channel is unavailable")

	newChannel := findChannel(channelId, s)

	if newChannel > 0 {
		s.AddClientToChannel(Clients{Conn: client, CurrentChannel: channelId, Status: "receiver"})
	}

	fmt.Printf("New client subscribed to the channel %d\n", channelId)
}

func (s *Server) SendData(data Messages, client net.Conn) {
	success := s.validateSubscription(client)
	indexClient := s.findClient(client)
	if indexClient < 0 {
		fmt.Println("The client doesn't exist")
		return
	}
	currentClient := &s.Clients[indexClient]

	commands := strings.Split(data.Message, " ")
	fileName := commands[1]

	if currentClient.Status == "receiver" {
		gob.NewEncoder(client).Encode(&Messages{Message: "The client is in receiver status"})
		return
	}
	if !success {
		msg := "The client cannot send data since it is not subscribed to a channel"
		err := gob.NewEncoder(client).Encode(&Messages{Message: msg})
		exceptions.ManageError(err, msg)
		return
	}
	for _, connection := range s.Clients {
		if connection.CurrentChannel == currentClient.CurrentChannel && connection.Conn.RemoteAddr() != client.RemoteAddr() {
			gob.NewEncoder(connection.Conn).Encode(&data)
			s.FileHistory = append(s.FileHistory, FileHistory{Channel: currentClient.CurrentChannel, Sender: currentClient.Conn.RemoteAddr().String(), Receiver: connection.Conn.RemoteAddr().String(), NameFile: fileName})
		}
	}

}

func (s *Server) ChangeStatusClient(client net.Conn, status string) {
	indexClient := s.findClient(client)

	if indexClient < 0 {
		gob.NewEncoder(client).Encode(&Messages{Message: "Client not found, subscribe first"})
		return
	}

	s.Clients[indexClient].Status = status
	msg := "Your status is " + status
	gob.NewEncoder(client).Encode(&Messages{Message: msg})
}

func (s Server) validateSubscription(client net.Conn) bool {
	for _, cnn := range s.Clients {
		if client.RemoteAddr() == cnn.Conn.RemoteAddr() {
			return true
		}
	}
	return false
}

func (s Server) findClient(client net.Conn) int {
	for index, cnn := range s.Clients {
		if client.RemoteAddr() == cnn.Conn.RemoteAddr() {
			return index
		}
	}
	return -1
}

func findChannel(id int, s *Server) int {
	if id < -1 {
		return id
	}
	existsChannel := false
	for _, channel := range s.channels {
		if channel == id {
			existsChannel = true
			return id
		}
	}

	if !existsChannel {
		s.NewChannel(id)
	}
	return id
}
