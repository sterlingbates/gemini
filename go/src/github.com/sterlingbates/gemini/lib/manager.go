package manager

import (
    "bufio"
    "fmt"
    "net"
	"strings"
	"time"

    "github.com/goburrow/serial"
)

type OrbiterManager struct {
	conn net.Conn
	isOrbiterConnected bool
	port serial.Port
	isSerialConnected bool
}

func NewOrbiterManager() *OrbiterManager {
	tmpconn, err := net.Dial("tcp", "127.0.0.1:37777")
    ret := OrbiterManager {
		conn: tmpconn,
	}
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		ret.isOrbiterConnected = false
	} else {
		ret.isOrbiterConnected = true
	}

	c := serial.Config {
		Address: "COM4",
		BaudRate: 9600,
		DataBits: 8,
		Parity: "N",
		StopBits: 1,
		Timeout: time.Millisecond * 100,
	}
	
	tmpport, err := serial.Open(&c)
	ret.port = tmpport
	if err != nil {
		fmt.Printf("COM port error: %s\n", err)
		ret.isSerialConnected = false
	} else {
		ret.isSerialConnected = true
	}
	fmt.Printf("Initialized:\n\tOrbiter connected: %v\n\tArduino connected: %v\n\n", ret.isOrbiterConnected, ret.isSerialConnected)
	return &ret
}

func (t *OrbiterManager) IsOrbiterConnected() (bool) {
	return t.isOrbiterConnected
}

func (t *OrbiterManager) Send(content string) (string) {
	fmt.Println("Writing: " + content)
	fmt.Fprintf(t.conn, content + "\r\n")
	fmt.Println("Reading")
	s, err := bufio.NewReader(t.conn).ReadString(byte('\r'))
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Println("Done")
	return strings.TrimSpace(s)
}

func (t *OrbiterManager) ReadyComPort() (string) {
	buf := make([]byte, 128)
	n, _ := t.port.Read(buf)
	if n > 0 {
		return strings.TrimSpace(string(buf[:n]))
	}
	return ""
}

func (t *OrbiterManager) ReadComPort() (string) {
	buf := make([]byte, 128)
	result := ""
	for {
		n, err := t.port.Read(buf)
		if n == 0 || err != nil {
			return strings.TrimSpace(result)
		}
		result += string(buf[:n])
	}
	return ""
}

func (t *OrbiterManager) SendComPort(content string) () {
	fmt.Println("Writing to COM: " + content)
	
	_, err := t.port.Write([]byte(content + "\n"))
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}

func (t *OrbiterManager) Cleanup() () {
}