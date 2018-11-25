package main

import (
	"fmt"
	"net"

	ssnr "github.com/Jonathas-Conceicao/ssnrgo"
)

func main() {
	testListing()
	testMessage()
}

func testMessage() {
	cn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic("Failed to dial host")
	}

	message := ssnr.TestNotification()
	cn.Write(message.Encode())
	fmt.Println("Message sent!")
}

func testListing() {
	cn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic("Failed to dial host")
	}

	listing := ssnr.TestListingRequest()
	cn.Write(listing.Encode())
	fmt.Println("Request sent!")

	tmp := make([]byte, 500)
	_, err = cn.Read(tmp)
	if err != nil {
		panic("Error found when reading message")
	}
	if tmp[0] == 0 {
		panic("Error message received back!")
	}

	listing = ssnr.DecodeListingReceived(tmp)
	fmt.Println("List received!")
	fmt.Println(listing)
}
