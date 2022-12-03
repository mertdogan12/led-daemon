package led

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mertdogan12/led-daemon/internal/uds"
)

var DeviceCount int = 1
var subnet string = "192.168.100."
var port string = "1337"

func Run() {
	for {
		switch uds.Mode {
		case "off":
			off()
		default:
			log.Printf("Mode %s does not exist. Mode set to default (off).", uds.Mode)
			uds.Mode = "off"
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func off() {
	for i := 2; i < 2+DeviceCount; i++ {
		addr := fmt.Sprintf("%s%d:%s", subnet, i, port)
		conn, err := net.Dial("udp", addr)

		if err != nil {
			log.Println("Could not open connection to:", addr)
			continue
		}

		_, err = conn.Write([]byte{0x00, 0x00, 0x00, 0x00})
		if err != nil {
			log.Println("Could not send a udp package to:", addr)
			conn.Close()
			continue
		}

		conn.Close()
	}
}
