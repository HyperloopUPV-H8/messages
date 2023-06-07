package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter"
	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

var addrFlag = flag.String("a", "127.0.0.4:50500", "address of the tcp server")
var messageDelay = flag.Duration("md", time.Second, "delay between messages")
var orderDelay = flag.Duration("od", time.Second, "delay between state orders")
var configPath = flag.String("c", "./config.toml", "path to the config")

func main() {
	flag.Parse()

	config, err := getConfig(*configPath)
	if err != nil {
		log.Fatalln(color.RedString("error reading config: %s", err))
	}

	excelAdapter := excel_adapter.New(config.Excel)
	boards := excelAdapter.GetBoards()
	globalInfo := excelAdapter.GetGlobalInfo()

	kindToId, err := parseIds(globalInfo.MessageToId)
	if err != nil {
		log.Fatalln(color.RedString("error parsing message ids: %s", err))
	}
	boardToId, err := parseIds(globalInfo.BoardToId)
	if err != nil {
		log.Fatalln(color.RedString("error parsing board ids: %s", err))
	}

	msgGenerator := NewMessageGenerator(exclude(exclude(kindToId, "state_orders"), "blcu_ack"), boardToId)

	ordGenerator := NewOrderGenerator(kindToId["state_orders"], getOrders(boards), boardToId)

	listener, err := createListener(*addrFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Printf("Listening on %s\n", *addrFlag)

	wg := &sync.WaitGroup{}
	defer wg.Wait()
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(color.RedString("Error accepting connection"))
			break
		}
		wg.Add(1)
		go handleConn(conn, wg, &msgGenerator, &ordGenerator)
	}

}

func createListener(addr string) (*net.TCPListener, error) {
	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "resolve tcp address")
	}

	listener, err := net.ListenTCP("tcp", laddr)
	return listener, errors.Wrap(err, "listen tcp")
}

func handleConn(conn *net.TCPConn, wg *sync.WaitGroup, msgGenerator *MessageGenerator, ordGenerator *OrderGenerator) {
	defer wg.Done()
	defer conn.Close()
	defer log.Println(color.RedString("[%s] Disconnect", conn.RemoteAddr()))
	log.Println(color.CyanString("[%s] Connect", conn.RemoteAddr()))

	go handleRead(conn)

	msg_ticker := time.NewTicker(*messageDelay)
	defer msg_ticker.Stop()
	ord_ticker := time.NewTicker(*orderDelay)
	defer ord_ticker.Stop()
	for {
		select {
		case <-msg_ticker.C:
			err := sendMessage(conn, msgGenerator)
			if err != nil {
				log.Println(color.RedString("[%s] Error writing: %s", err))
				return
			}
		case <-ord_ticker.C:
			err := sendOrder(conn, ordGenerator)
			if err != nil {
				log.Println(color.RedString("[%s] Error writing: %s", err))
				return
			}
		}
	}
}

func sendMessage(conn *net.TCPConn, msgGenerator *MessageGenerator) error {
	message, id := msgGenerator.New()
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	buffer := make([]byte, 2)
	binary.LittleEndian.PutUint16(buffer, uint16(id))
	buffer = append(buffer, '\n', '\n')
	buffer = append(buffer, payload...)

	_, err = conn.Write(append(buffer, 0x00))
	log.Println(color.GreenString("[%s] Write (%d)", conn.RemoteAddr(), id))
	return err
}

func sendOrder(conn *net.TCPConn, ordGenerator *OrderGenerator) error {
	stateOrders, err := ordGenerator.New()
	if err != nil {
		return err
	}

	_, err = conn.Write(stateOrders)
	log.Println(color.GreenString("[%s] Write State Orders", conn.RemoteAddr()))
	return err
}

func handleRead(conn *net.TCPConn) {
	buf := make([]byte, 1500)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			log.Println(color.RedString("[%s] Error reading: %s", err))
			return
		}
		log.Println(color.BlueString("[%s] Read (%d)", conn.RemoteAddr(), binary.LittleEndian.Uint16(buf[:2])))
	}

}

func parseIds(literal map[string]string) (map[string]uint16, error) {
	result := make(map[string]uint16, len(literal))
	for key, value := range literal {
		id, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return nil, err
		}
		result[key] = uint16(id)
	}
	return result, nil
}

func getOrders(boards map[string]models.Board) map[string][]uint16 {
	orders := make(map[string][]uint16, len(boards))
	for _, board := range boards {
		stateOrders := make([]uint16, 0, len(board.Packets))
		for _, packet := range board.Packets {
			if packet.Description.Type == "stateOrder" {
				id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
				if err != nil {
					log.Fatalln(color.RedString("error parsing order id: %s", err))
				}
				stateOrders = append(stateOrders, uint16(id))
			}
		}
		orders[board.Name] = stateOrders
	}
	return orders
}

func exclude[K comparable, V any](m map[K]V, item K) map[K]V {
	new := make(map[K]V)
	for k, v := range m {
		if k != item {
			new[k] = v
		}
	}
	return new
}
