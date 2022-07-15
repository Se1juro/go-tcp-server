package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/upload-files-go/api"
	"github.com/upload-files-go/models"
)

func start(startServer string, server *models.Server, _ *sync.WaitGroup) {
	listener, err := net.Listen("tcp", startServer) // Create tcp server on 3035 port
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	log.Printf("Server listening on port %s", startServer)

	server.NewChannel(1)
	server.NewClient(listener)
}

func main() {
	var startServer string
	waitGroup := &sync.WaitGroup{}

	waitGroup.Add(1)

	flag.StringVar(&startServer, "start", "", "Start a TCP server")
	flag.Parse()
	server := models.Server{}

	if startServer != "" {
		log.Printf("Try to initialize on port %s", startServer)
		go start(":"+startServer, &server, waitGroup)

	}
	api.Api(&server)
	waitGroup.Wait()

	log.Println("Closing...")
}
