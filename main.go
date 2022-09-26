package main

import (
	"fmt"
	"net"
)

func main() {
	addr := "localhost:8888"
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

	for {
		conn, err := listener.Accept()
		fmt.Printf("Local connection address: %s\n", conn.LocalAddr().String())
		fmt.Printf("Remote connection address: %s\n", conn.RemoteAddr().String())
		if err != nil {
			panic(err)
		}

		go func(conn net.Conn) {
			buf := make([]byte, 1024)
			len, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("Error reading: %#v\n", err)
				return
			}
			fmt.Printf("Message received: %s\n", string(buf[:len]))
			conn.Write([]byte("Message received.\n"))
			conn.Close()
		}(conn)
	}
}
