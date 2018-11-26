package main

import (
	"fmt"
	"go/build"
	"net"
	"os"
	"strconv"

	ssnr "github.com/Jonathas-Conceicao/ssnrgo"
)

func main() {
	confFile := build.Default.GOPATH + "/configs/ssnr_sender_config.json"
	config := ssnr.NewConfig(confFile)

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
		requestUsers(config)

	case "send":
		if l != 4 {
			printHelp()
			return
		}
		target, err := strconv.ParseInt(os.Args[2], 0, 16)
		if err != nil {
			printHelp()
			return
		}
		sendMessage(config, uint16(target), &os.Args[3])

	default:
		printHelp()
	}
}

// TODO: Add help message for sender
func printHelp() {
	helpMessage := "No help message just yet"
	fmt.Println(helpMessage)
}

func sendMessage(config *ssnr.Config, recv uint16, content *string) {
	cn, err := net.Dial("tcp", config.Host+config.Port)
	if err != nil {
		panic("Failed to dial host")
	}

	var message *ssnr.Notification
	sndr := config.Name
	if sndr == "" {
		message = ssnr.NewAnonymousNotification(recv, *content)
	} else {
		message = ssnr.NewNotification(recv, sndr, *content)
	}
	cn.Write(message.Encode())
	fmt.Println("Message sent!")
}

func requestUsers(config *ssnr.Config) {
	cn, err := net.Dial("tcp", config.Host+config.Port)
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
