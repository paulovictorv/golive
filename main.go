package main

import (
	"br.com.pmelo/golive/app"
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	reader := bufio.NewReader(os.Stdin)

	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Action: func(c *cli.Context) error {
				fmt.Println("Enter app name:")
				appName, _ := reader.ReadString('\n')
				application.CreateApp(appName)
				return nil
			},
		},
	}

	app.Run(os.Args)
}
