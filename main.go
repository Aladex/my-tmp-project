package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	// Listen for incoming connections.
	l, err := net.Listen("tcp", "0.0.0.0:514")
	if err != nil {
		log.Fatalln("Error listening:", err.Error())
	}
	// Close the listener when the application closes.
	defer l.Close()
	log.Println("Listening on 0.0.0.0:514")
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatalln("Error accepting: ", err.Error())
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	defer conn.Close()

	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(conn)
	)

ILOOP:
	for {
		n, err := r.Read(buf)
		data := string(buf[:n])

		switch err {
		case io.EOF:
			break ILOOP
		case nil:
			// log.Println("Receive:", data)
			go writeToLogFile(data, conn.RemoteAddr().String())
			if isTransportOver(data) {
				break ILOOP
			}

		default:
			log.Fatalf("Receive data failed:%s", err)
			return
		}

	}
}

func writeToLogFile(mess, ip string) {
	thisTime := time.Now()
	fmt.Println(fmt.Sprintf("%s: %s - %s", thisTime.Format("2006-01-02 03:04:05 MST"), ip, mess))
	f, err := os.OpenFile("text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%s: %s - %s\n", thisTime.Format("2006-01-02 03:04:05 MST"), ip, mess)); err != nil {
		log.Println(err)
	}
}

func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, "\r")
	return
}
