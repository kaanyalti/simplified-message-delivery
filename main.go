package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type Server struct {
	Connections map[uuid.UUID]Connection
	MaxConnNum  int
}

type Connection struct {
	ID   uuid.UUID
	Conn net.Conn
}

func main() {
	addr := fmt.Sprintf("localhost:%s", os.Args[1])
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	host, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		panic(err)
	}
	fmt.Printf("listening on host: %s, port: %s\n", host, port)
	server := Server{
		Connections: make(map[uuid.UUID]Connection),
		MaxConnNum:  2,
	}
	for {
		if len(server.Connections) >= server.MaxConnNum {
			continue
		}
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		connectionId, _ := uuid.NewUUID()
		connection := Connection{
			ID:   connectionId,
			Conn: conn,
		}
		server.Connections[connectionId] = connection
		fmt.Printf("Current number of connections %s\n", strconv.Itoa(len(server.Connections)))
		fmt.Printf("Local connection address: %s\n", conn.LocalAddr().String())
		fmt.Printf("Remote connection address: %s\n", conn.RemoteAddr().String())
		if err != nil {
			panic(err)
		}

		go server.HandleConnection(connection)
	}
}

func (s *Server) HandleConnection(c Connection) {
	fmt.Printf("Serving %s\n", c.Conn.RemoteAddr().String())
L:
	for {
		data, err := bufio.NewReader(c.Conn).ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		trimmed := strings.ToUpper(strings.TrimSpace(string(data)))
		fmt.Printf("Incoming data: %s\n", trimmed)

		switch trimmed {
		case "STOP":
			fmt.Println("Closing connection")
			break L
		case "CLIENTS":
			ids := s.GetConnectionIds(c.ID)
			c.Conn.Write([]byte(fmt.Sprintf("%s\n", ids)))
		case "SELF":
			c.Conn.Write([]byte(fmt.Sprintf("Your ID: %s\n", c.ID.String())))
		}
	}
	err := c.Conn.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Connection successfully closed")
	connectionId := c.ID
	delete(s.Connections, connectionId)
}

func (s *Server) GetConnectionIds(excludeId uuid.UUID) []string {
	ids := make([]string, len(s.Connections))
	for id, _ := range s.Connections {
		if excludeId != id {
			ids = append(ids, id.String())
		}
	}
	return ids
}
