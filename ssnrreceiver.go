package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	notify "github.com/TheCreeper/go-notify"
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
			Value: "localhost",
			Usage: "Host's address",
		},
		cli.StringFlag{
			Name:  "name, n",
			Usage: "Sender's name",
		},
		cli.IntFlag{
			Name:  "code, c",
			Value: 0,
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

		reader, err := handleLogin(
			conn,
			uint16(c.Int("code")),
			config.Name)
		if err != nil {
			return err
		}

		for {
			code, err := reader.Peek(1)
			if err == io.EOF {
				log.Println("Server was closed")
				return err
			}
			if err != nil {
				return err
			}

			switch code[0] {
			case ssnr.NotificationCode:
				_, ntf, err := ssnr.ReadNotification(reader)
				if err != nil {
					return err
				}
				err = display(ntf)
				if err != nil {
					return err
				}
			case ssnr.PingCode:
				data := make([]byte, 2)
				n, err := reader.Read(data)
				if err != nil {
					return err
				}
				if n == 1 {
					_, err := reader.Read(data[1:])
					if err != nil {
						return err
					}
				}
				log.Println("Received ping from host")
			default:
				data := make([]byte, 1)
				v, err := reader.Read(data)
				if err != nil {
					return err
				}
				return errors.New("Recived invalid code:" + string(v))
			}
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Printf("Error of type: %T", err)
		log.Fatal(err)
	}
}

func display(n *ssnr.Notification) error {
	log.Println("Notification Received",
		"from: \""+n.GetEmitter()+"\"",
		n.GetMessage())
	ntf := notify.NewNotification(
		"SSNR Notification",
		fmt.Sprintf("[%s] - %s\n%s",
			n.GetEmitter(),
			n.GetTimeString(),
			n.GetMessage()))
	_, err := ntf.Show()
	return err
}

func handleLogin(cn net.Conn, code uint16, name string) (
	*bufio.Reader, error) {
	request := ssnr.NewRegister(code, name)
	log.Println("Requesting register")
	_, err := cn.Write(request.Encode())
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(cn)
	_, request, err = ssnr.ReadRegister(r)

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

func handleUnknown(rd *bufio.Reader) error {
	tmp := make([]byte, 1)
	_, err := rd.Read(tmp)
	if err != nil {
		return err
	}
	log.Printf("Invalid message received!\nCode: %d\n", tmp[0])
	return nil
}
