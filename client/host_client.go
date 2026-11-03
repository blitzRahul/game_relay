package main

import (
	"io"
	"log"
	"net"
)

const (
	aciIP     = "4.187.182.59"
	localGame = "127.0.0.1:6767"
)

func main() {
	ctrl, err := net.Dial("tcp", aciIP+":7000")
	if err != nil {
		log.Fatal("Failed to connect to Azure relay:", err)
	}
	// Send the 4-byte magic handshake so Azure knows we aren't a health probe
	ctrl.Write([]byte("CTRL"))
	log.Println("Connected and authenticated with Azure Relay!")

	buf := make([]byte, 1)
	for {
		_, err := ctrl.Read(buf)
		if err != nil {
			log.Fatal("Relay control connection closed:", err)
		}

		go func() {
			gameConn, err := net.Dial("tcp", localGame)
			if err != nil {
				log.Println("Local server not running on", localGame)
				return
			}
			relayData, err := net.Dial("tcp", aciIP+":7002")
			if err != nil {
				gameConn.Close()
				return
			}
			// Authenticate the data pipe
			relayData.Write([]byte("DATA"))

			defer gameConn.Close()
			defer relayData.Close()
			go io.Copy(gameConn, relayData)
			io.Copy(relayData, gameConn)
		}()
	}
}
