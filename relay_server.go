package main

import (
	"io"
	"log"
	"net"
)

func main() {
	controlLn, _ := net.Listen("tcp", ":7000") // Host control port
	dataLn, _ := net.Listen("tcp", ":7002")    // Host data pool port
	playerLn, _ := net.Listen("tcp", ":7001")  // Public port

	log.Println("relay listening on :7000, 7001, 7002")

	for {
		//wait for host to register
		hostCtrl, err := controlLn.Accept()
		if err != nil {
			continue
		}
		log.Println("Host connected!")

		//handle players while this host is connected
		for {
			playerConn, err := playerLn.Accept()
			if err != nil {
				break
			}

			//tell host to create new data pipe
			_, err = hostCtrl.Write([]byte{1})
			if err != nil {
				playerConn.Close()
				break
			}

			//accept fresh data connection
			hostData, err := dataLn.Accept()
			if err != nil {
				playerConn.Close()
				continue
			}

			//pipe bytes in both directions
			go func(p, h net.Conn) {
				defer p.Close()
				defer h.Close()
				go io.Copy(p, h)
				io.Copy(h, p)
			}(playerConn, hostData)
		}
	}
}
