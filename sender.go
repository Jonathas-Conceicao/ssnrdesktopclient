package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/urfave/cli"

	ssnr "github.com/Jonathas-Conceicao/ssnrgo"
)

func main() {
	app := cli.NewApp()
	app.Name = "SSNR desktop sender CLI"
	app.Usage = "Send distributed notifications over SSNR protocol"
	app.Version = "0.1.0"

	cli.HelpFlag = cli.BoolFlag{
		Name:  "help",
		Usage: "show this dialog",
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port, p",
			Value: ":30106",
			Usage: "Server port",
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
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "Request server the list of avaliable users",
			Action: func(c *cli.Context) error {
				config, err := newConfig(c)
				if err != nil {
					return err
				}
				return requestUsers(config)
			},
		},

		cli.Command{
			Name:    "send",
			Aliases: []string{"s"},
			Usage:   "Send the `user` a `notification`",
			Action: func(c *cli.Context) error {
				config, err := newConfig(c)
				if err != nil {
					return err
				}
				whom, err := strconv.Atoi(c.Args().Get(0))
				if err != nil {
					return err
				}
				return sendMessage(config, uint16(whom), c.Args().Get(1))
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error of type: %T", err)
		panic(err)
	}
}

func newConfig(c *cli.Context) (*ssnr.Config, error) {
	return ssnr.NewConfig(
		c.Parent().String("host"),
		c.Parent().String("port"),
		c.Parent().String("name"))
}

func sendMessage(config *ssnr.Config, recv uint16, content string) error {
	cn, err := net.Dial("tcp", config.Host+config.Port)
	if err != nil {
		return err
	}
	message := ssnr.NewNotification(recv, config.Name, content)
	cn.Write(message.Encode())
	return nil
}

func requestUsers(config *ssnr.Config) error {
	cn, err := net.Dial("tcp", config.Host+config.Port)
	if err != nil {
		return err
	}

	listing := ssnr.NewListingRequestAll()
	cn.Write(listing.Encode())

	_, listing, err = ssnr.ReadListing(bufio.NewReader(cn), false)
	if err != nil {
		return err
	}
	fmt.Println(listing)
	return nil
}
