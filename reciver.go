package main

import (
	"bufio"
	"errors"
	"log"
	"net"
	"os"

	"github.com/urfave/cli"

	ssnr "github.com/Jonathas-Conceicao/ssnrgo"
)

func main() {
	app := cli.NewApp()
	app.Name = "SSNR desktop reciver APP"
	app.Usage = "Recive distributed notifications over SSNR protocol"
	app.Version = "0.1.0"

	cli.HelpFlag = cli.BoolFlag{
		Name:  "help",
		Usage: "show this dialog",
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port, p",
			Value: ":30106",
			Usage: "Application's port",
		},
		cli.StringFlag{
			Name:  "host, h",
			Value: "196.168.0.1",
			Usage: "Host's address",
		},
		cli.StringFlag{
			Name:  "name, n",
			Usage: "Sender's name",
		},
		cli.IntFlag{
			Name:  "code, c",
			Usage: "Desired register ID (16 bits unsigned integer)",
		},
	}

	app.Action = func(c *cli.Context) error {
		config, err := ssnr.NewConfig(
			c.String("host"),
			c.String("port"),
			c.String("name"))
		if err != nil {
			return err
		}

		conn, err := net.Dial("tcp", config.Host+config.Port)
		if err != nil {
			return err
		}

		reader, err := handleLogin(conn, uint16(c.Int("code")))
		if err != nil {
			return err
		}
		for {
			n, err := readNextNotification(reader)
			if err != nil {
				return err
			}
			err = display(n)
			if err != nil {
				return err
			}
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func display(n *ssnr.Notification) error { return nil }

func handleLogin(cn net.Conn, code uint16) (*bufio.Reader, error) {
	request := ssnr.NewRegister(code, "Jonathas")
	log.Println("Requesting register")
	_, err := cn.Write(request.Encode())
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(cn)
	tmp := make([]byte, 500)
	_, err = r.Read(tmp)
	request, err = ssnr.DecodeRegister(tmp)
	if err != nil {
		return r, err
	}

	switch request.GetReturn() {
	case ssnr.ConnAccepted:
		return r, nil
	case ssnr.ConnNewAddres:
		log.Println("Connected to id: ", request.GetReceptor())
		return r, nil
	case ssnr.RefServerFull:
		return r, errors.New("Connection refused, server is full")
	case ssnr.RefBlackList:
		return r, errors.New("Connection refused, blacklist")
	case ssnr.RefUnknowEror:
		return r, errors.New("Connection refused, error not informed")
	default:
		return r, errors.New("Invalid error message returned")
	}
}

func readNextNotification(rd *bufio.Reader) (*ssnr.Notification, error) {
	return nil, errors.New("TO DO")
}

// func handleNotification(data []byte) error {
// 	message := ssnr.DecodeNotification(data)
// 	log.Println("Message Received:\n" + message.String())
// 	log.Println("Current list of users: ", users.Length())
// 	log.Print(users)
// 	return nil
// }

func handleUnknown(rd *bufio.Reader) error {
	tmp := make([]byte, 1)
	_, err := rd.Read(tmp)
	if err != nil {
		return err
	}
	log.Printf("Invalid message received!\nCode: %d\n", tmp[0])
	return nil
}
