package main

import (
	"bufio"
	"fmt"
	"net"

	ssnr "github.com/Jonathas-Conceicao/ssnrgo"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic("Failed to dial host")
	}

	message := ssnr.TestingThings()
	conn.Write(message.Encode())

	// listen for reply
	r, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message from server: " + r)
}
