package uds

import (
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ColorStruct struct {
	Red   uint16
	Green uint16
	Blue  uint16
}

const SockAddr = "/tmp/led.sock"

var Mode string = "off"
var Color ColorStruct = ColorStruct{
	Red:   0,
	Green: 0,
	Blue:  0,
}

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
		re := regexp.MustCompile(`[^a-z 0-9]`)
		data = re.ReplaceAll(data[:n], make([]byte, 0))

		Mode = string(data)

		println(Mode)

		for i, m := range strings.Split(Mode, " ") {
			switch i {
			case 0:
				Mode = m

			case 1:
				r, err := strconv.Atoi(m)
				if err != nil {
					answer(conn, "Could not parse color (as integer)"+m)
				}

				Color.Red = uint16(r)

			case 2:
				g, err := strconv.Atoi(m)
				if err != nil {
					answer(conn, "Could not parse color (as integer)"+m)
					break
				}

				Color.Green = uint16(g)

			case 3:
				b, err := strconv.Atoi(m)
				if err != nil {
					answer(conn, "Could not parse color (as integer)"+m)
					break
				}

				Color.Blue = uint16(b)

				println(Color.Blue)
			}
		}

		answer(conn, "Mode changed to: "+Mode)

		err = conn.Close()
		if err != nil {
			log.Fatal("Error while closing the connection", err)
		}
	}
}

func answer(c net.Conn, text string) error {
	text += "\n"

	log.Print(text)
	_, err := c.Write([]byte(text))

	return err
}
