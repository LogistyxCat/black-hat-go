package main

import (
	"io"
	"log"
	"net"
)

// Old  echo(conn net.Conn function; poorly designed, updated below
//func echo(conn net.Conn) {
//	defer conn.Close()
//
//	// Create a buffer to store received data
//	b := make([]byte, 512)
//	for {
//		// Receive data via conn.Read
//		size, err := conn.Read(b[1:])
//		if err == io.EOF {
//			log.Println("Client disconnected")
//			break
//		}
//		if err != nil {
//			log.Println("Unexpected error")
//			break
//		}
//		log.Printf("Received %d bytes: %s", size, string(b))
//
//		// Send data via conn.Write
//		log.Println("Writing data")
//		if _, err := conn.Write(b[0:size]); err != nil {
//			log.Fatalln("Unable to write data")
//		}
//	}
//}

// Improved echo(conn net.Conn); utilizes bufio to simplify io buffer
//func echo(conn net.Conn) {
//	defer conn.Close()
//
//	reader := bufio.NewReader(conn)
//	s, err := reader.ReadString('\n')
//	if err != nil {
//		log.Fatalln("Unable to read data")
//	}
//	log.Printf("Read %d bytes: %s", len(s), s)
//
//	log.Println("Writing data")
//	writer := bufio.NewWriter(conn)
//	if _, err := writer.WriteString(s); err != nil {
//		log.Fatalln("Unable to write data")
//	}
//	writer.Flush()
//}

// further improved echo(conn net.Conn) function; uses io.Copy to copy from
// reader to writer directly
func echo(conn net.Conn) {
	defer conn.Close()
	// Copy from io.Reader to io.Writer via io.Copy()
	if _, err := io.Copy(conn, conn); err != nil {
		log.Fatalln("Unable to read/write data")
	}
}

func main() {
	// Bind to TCP 20080 on all interfaces
	listener, err := net.Listen("tcp", ":20080")
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}
	log.Println("Listening on port 20080")
	for {
		// Wait for a connection. Create net.Conn on connection established.
		conn, err := listener.Accept()
		log.Println("Received connection")
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}
		// Handle the connection; uses goroutine for concurrency
		go echo(conn)
	}
}
