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
		log.Fatal("could not connect to relay:", err)
	}
	log.Println("connected to relay")

	buf := make([]byte, 1)
	for {
		// wait for 1 byte join signal
		_, err := ctrl.Read(buf)
		if err != nil {
			log.Fatal("relay connection closed", err)
		}

		go func() {
			// bridge local and cloud data port
			gameConn, err := net.Dial("tcp", localGame)
			if err != nil {
				log.Println("local server not running on", localGame)
				return
			}
			relayData, err := net.Dial("tcp", aciIP+":7002")
			if err != nil {
				gameConn.Close()
				return
			}

			defer gameConn.Close()
			defer relayData.Close()
			go io.Copy(gameConn, relayData)
			io.Copy(relayData, gameConn)
		}()
	}
}
