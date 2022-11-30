package uds

import (
	"log"
	"net"
	"os"
)

const SockAddr = "/tmp/led.sock"

func Run() {
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("unix", SockAddr)
	if err != nil {
		log.Fatal("Listen error", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}

		log.Printf("Client connected [%s]", conn.RemoteAddr().Network())

		data := make([]byte, 256)

		n, err := conn.Read(data)

		if err != nil {
			log.Fatal("Error while reading data from the client:", err)
		}

		log.Print("Data received: ", string(data[:n]))

		err = conn.Close()
		if err != nil {
			log.Fatal("Error while closing the connection", err)
		}
	}
}
