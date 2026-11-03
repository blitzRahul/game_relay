package main

import (
	"io"
	"log"
	"net"
	"time"
)

// Helper that reads a 4-byte magic string to filter out Azure health probes
func verifyHandshake(conn net.Conn, expected string) bool {
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, len(expected))
	_, err := io.ReadFull(conn, buf)
	conn.SetReadDeadline(time.Time{}) // Clear deadline
	return err == nil && string(buf) == expected
}

func main() {
	controlLn, _ := net.Listen("tcp", ":7000")
	dataLn, _ := net.Listen("tcp", ":7002")
	playerLn, _ := net.Listen("tcp", ":7001")

	log.Println("Relay listening on :7000 (ctrl), :7001 (player), :7002 (data)")

	for {
		hostCtrl, err := controlLn.Accept()
		if err != nil {
			continue
		}

		// Ignore Azure health probes and scanners on :7000
		if !verifyHandshake(hostCtrl, "CTRL") {
			hostCtrl.Close()
			continue
		}
		log.Println("Real Host authenticated and connected from:", hostCtrl.RemoteAddr())

		for {
			playerConn, err := playerLn.Accept()
			if err != nil {
				break
			}

			// Tell Host to open a data pipe
			_, err = hostCtrl.Write([]byte{1})
			if err != nil {
				log.Println("Host disconnected.")
				playerConn.Close()
				break
			}

			// Wait for a valid DATA connection from the Host on :7002
			var hostData net.Conn
			for {
				cand, err := dataLn.Accept()
				if err != nil {
					break
				}
				if verifyHandshake(cand, "DATA") {
					hostData = cand
					break
				}
				cand.Close() // Discard health probe on :7002
			}

			if hostData == nil {
				playerConn.Close()
				continue
			}

			go func(p, h net.Conn) {
				defer p.Close()
				defer h.Close()
				go io.Copy(p, h)
				io.Copy(h, p)
			}(playerConn, hostData)
		}
	}
}
