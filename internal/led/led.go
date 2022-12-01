package led

import (
	"log"
	"time"

	"github.com/mertdogan12/led-daemon/internal/uds"
)

func Run() {
	for {
		switch uds.Mode {
		case "off":
			off()
		default:
			log.Printf("Mode %s does not exist. Mode set to default (off).", uds.Mode)
			uds.Mode = "off"
		}
	}
}

func off() {
	time.Sleep(5 * time.Second)
}
