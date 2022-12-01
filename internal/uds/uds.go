package uds

import (
	"log"
	"net"
	"os"
	"regexp"
)

const SockAddr = "/tmp/led.sock"

var Mode string = "off"

func Run() {
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("unix", SockAddr)
	if err != nil {
		log.Fatal("Listen error", err)
	}
	defer l.Close()

	log.Println("UDS startet at /tmp/led.sock")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}

		log.Printf("Client connected [%s]", conn.RemoteAddr().Network())

		data := make([]byte, 32)

		n, err := conn.Read(data)
		if err != nil {
			log.Println("Error while reading data from the client:", err)
			conn.Close()
			continue
		}

		// Remove spaces and newlines from mode
		re := regexp.MustCompile(`[^a-z]`)
		data = re.ReplaceAll(data[:n], make([]byte, 0))

		log.Print("Mode changed to: ", string(data))

		Mode = string(data)

		err = conn.Close()
		if err != nil {
			log.Fatal("Error while closing the connection", err)
		}
	}
}
