package led

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mazznoer/colorgrad"
	"github.com/mertdogan12/led-daemon/internal/uds"
)

var DeviceCount int = 1
var subnet string = "192.168.100."
var port string = "1337"

var _red uint16 = 255
var _green uint16 = 0
var _blue uint16 = 0

var grad = colorgrad.Rainbow()

func Run() {
	for {
		switch uds.Mode {
		case "off":
			off()

		case "color":
			colorEffect()

		case "fade":
			fade()

		default:
			log.Printf("Mode %s does not exist. Mode set to default (off).", uds.Mode)
			uds.Mode = "off"
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func off() {
	changeColor(0, 0, 0)

	time.Sleep(time.Second)
}

func colorEffect() {
	changeColor(uds.Color.Red, uds.Color.Green, uds.Color.Blue)

	time.Sleep(time.Second)
}

var speed float64 = 0.01
var pos float64 = 0

func fade() {
	red := grad.At(pos).R * 100
	green := grad.At(pos).G * 100
	blue := grad.At(pos).B * 100

	changeColor(
		uint16(red),
		uint16(green),
		uint16(blue),
	)

	pos += speed
}

func changeColor(red uint16, green uint16, blue uint16) {
	fmt.Println(red, green, blue)

	multi := uint16(1024 / 255)

	data := make([]byte, 0)
	tmp := make([]byte, 2)

	binary.LittleEndian.PutUint16(tmp, 1024-red*multi) // Red
	data = append(data, tmp...)

	binary.LittleEndian.PutUint16(tmp, 1024-green*multi) // Green
	data = append(data, tmp...)

	binary.LittleEndian.PutUint16(tmp, 1024-blue*multi) // Blue
	data = append(data, tmp...)

	println(data[0], data[1])

	sendData(data)
}

func sendData(data []byte) {
	if len(data) > 6 {
		log.Fatal("Send data is to big")
	}

	for i := 2; i < 2+DeviceCount; i++ {
		addr := fmt.Sprintf("%s%d:%s", subnet, i, port)
		conn, err := net.Dial("udp", addr)

		if err != nil {
			log.Println("Could not open connection to:", addr)
			continue
		}

		_, err = conn.Write(data)
		if err != nil {
			log.Println("Could not send a udp package to:", addr)
			conn.Close()
			continue
		}

		conn.Close()
	}

}
