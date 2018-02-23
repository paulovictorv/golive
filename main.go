package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"github.com/paulovictorv/golive/app"
)

func main() {
	app := cli.NewApp()
	reader := bufio.NewReader(os.Stdin)

	app.Commands = []cli.Command{
		{
			Name: "generate",
			Aliases: []string{"g"},
			Action: func(c *cli.Context) error {
				fmt.Println("Enter app name:")
				appName, _ := reader.ReadString('\n')
				createApp, err := golive.CreateApp(appName)

				if err != nil {
					return err
				} else {
					fmt.Printf("Created app with name %s", createApp.Name)
					return nil
				}
			},
		},
		{
			Name: "deploy",
			Aliases: []string{"d"},
			Action: func(c *cli.Context) error {
				fmt.Println("Enter app name:")
				appName, _ := reader.ReadString('\n')
				createApp, err := golive.CreateApp(appName)

				if err != nil {
					return err
				} else {
					fmt.Printf("Created app with name %s", createApp.Name)
					return nil
				}
			},
		},
	}

	app.Run(os.Args)
}
