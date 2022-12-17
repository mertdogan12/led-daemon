package uds

import (
	"fmt"
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

var FadeSpeed float64 = 0.01
var BlinkSpeed float64 = 0.01

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

		data := make([]byte, 256)

		n, err := conn.Read(data)
		if err != nil {
			log.Println("Error while reading data from the client:", err)
			conn.Close()
			continue
		}

		re := regexp.MustCompile(`[^a-z 0-9.]`)
		data = re.ReplaceAll(data[:n], make([]byte, 0))

		dataString := strings.Split(string(data), " ")
		Mode = dataString[0]

		switch Mode {
		case "color":
			if len(dataString) != 3+1 {
				answer(conn, "No Color is given")
				break
			}

			r, err := strconv.Atoi(dataString[1])
			g, err := strconv.Atoi(dataString[2])
			b, err := strconv.Atoi(dataString[3])
			if err != nil {
				answer(
					conn,
					fmt.Sprintf("Could not parse colors %s %s %s", dataString[1], dataString[2], dataString[3]),
				)
				break
			}

			Color = ColorStruct{
				Red:   uint16(r),
				Green: uint16(g),
				Blue:  uint16(b),
			}

		case "fade", "blink":
			if len(dataString) != 1+1 {
				answer(conn, "No speed is given")
				break
			}

			speed, err := strconv.ParseFloat(dataString[1], 64)
			if err != nil {
				answer(
					conn,
					fmt.Sprint("Could not parse speed", dataString[1]),
				)
				break
			}

			if speed <= 0 || speed >= 1 {
				answer(conn, fmt.Sprintf("Speed must be between 0 and 1, %f", speed))
				break
			}

			if Mode == "fade" {
				FadeSpeed = speed
			} else {
				BlinkSpeed = speed
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
