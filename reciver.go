package main

import (
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

	fmt.Println("Message sent!")
}
