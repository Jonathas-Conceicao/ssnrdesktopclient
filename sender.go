package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	ssnr "github.com/Jonathas-Conceicao/ssnrgo"
)

func main() {
	l := len(os.Args)
	if l < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {

	case "list":
		if l != 2 {
			printHelp()
			return
		}
		requestUsers()

	case "send":
		switch l {
		case 4:
			target, err := strconv.ParseInt(os.Args[2], 0, 16)
			if err != nil {
				printHelp()
				return
			}
			sendMessage(uint16(target), nil, &os.Args[3])
		case 5:
			target, err := strconv.ParseInt(os.Args[2], 0, 16)
			if err != nil {
				printHelp()
				return
			}
			sendMessage(uint16(target), &os.Args[3], &os.Args[4])
		default:
			printHelp()
			return
		}

	default:
		printHelp()
	}
}

// TODO: Add help message for sender
func printHelp() {
	helpMessage := "No help message just yet"
	fmt.Println(helpMessage)
}

func sendMessage(recv uint16, sndr *string, content *string) {
	cn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic("Failed to dial host")
	}

	var message *ssnr.Notification
	if sndr == nil {
		message = ssnr.NewAnonymousNotification(recv, *content)
	} else {
		message = ssnr.NewNotification(recv, *sndr, *content)
	}
	cn.Write(message.Encode())
	fmt.Println("Message sent!")
}

func requestUsers() {
	cn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic("Failed to dial host")
	}

	listing := ssnr.NewListingRequestAll()
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
