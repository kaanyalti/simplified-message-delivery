package main

import (
	"bufio"
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
	}
	for {
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
	for {
		data, err := bufio.NewReader(c.Conn).ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println(err.Error())
			break
		}

		fmt.Printf("Incoming data: %s\n", data)
		c.Conn.Write([]byte(fmt.Sprintf("Your ID: %s\n", c.ID.String())))

		trimmed := strings.TrimSpace(string(data))

		switch trimmed {
		case "STOP":
			fmt.Println("Closing connection")
			break
		case "CLIENTS":
			ids := s.GetConnectionIds()
			c.Conn.Write([]byte(fmt.Sprintf("%s", ids)))
		}

	}
	c.Conn.Close()
	connectionId := c.ID
	delete(s.Connections, connectionId)
}

func (s *Server) GetConnectionIds() []string {
	ids := make([]string, len(s.Connections))
	for id, _ := range s.Connections {
		ids = append(ids, id.String())
	}
	return ids
}
