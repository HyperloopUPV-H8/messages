package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/fatih/color"
)

func main() {
	listener := createListener("127.0.0.4", 50500)
	defer listener.Close()

	for {

		color.White("Waiting for connection...")
		conn, err := listener.AcceptTCP()

		if err != nil {
			log.Fatal("Error accepting connection")
		}

		color.Cyan(fmt.Sprintf("Accepted conn > %s", conn.RemoteAddr().String()))

		var wg sync.WaitGroup

		wg.Add(1)
		go writeLoop(conn, &wg)
		wg.Add(1)
		go readLoop(conn, &wg)
		wg.Wait()

		color.Red("Connection closed")
	}

}

func createListener(l_ip string, l_port int) *net.TCPListener {
	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", l_ip, l_port))

	if err != nil {
		log.Fatalf("ResolveTCPAddr: %s\n", err)
	}

	listener, err := net.ListenTCP("tcp", laddr)

	if err != nil {
		log.Fatal("Error listening TCP")
	}

	return listener
}

func writeLoop(conn *net.TCPConn, wg *sync.WaitGroup) {

	generator := NewMessageGenerator()
	ticker := time.NewTicker(time.Millisecond * 50)

	for range ticker.C {
		message, id := generator.New()
		payload, err := json.Marshal(message)

		if err != nil {
			fmt.Printf("Error marshaling message %v %v\n", err, message)
			continue
		}

		idBuf := make([]byte, 2)
		binary.LittleEndian.PutUint16(idBuf, uint16(id))

		fullMessage := append(idBuf, '\n', '\n')
		fullMessage = append(fullMessage, payload...)

		_, err = conn.Write(fullMessage)

		color.White("Send message")
		if err != nil {
			fmt.Printf("Error sending message %v", err)
			break
		}

	}

	wg.Done()
	ticker.Stop()
}

func readLoop(conn *net.TCPConn, wg *sync.WaitGroup) {
	for {
		buf := make([]byte, 100)
		_, err := conn.Read(buf)
		fmt.Println(buf)
		if err != nil {
			break
		}
	}

	wg.Done()
}
