package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"os"
	"github.com/paulovictorv/golive/app"
	tm "github.com/buger/goterm"
	"github.com/paulovictorv/golive/app/util"
)

func main() {
	app := cli.NewApp()
	reader := bufio.NewReader(os.Stdin)

	qs := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{Message:"What's the name of your app?"},
			Validate: survey.Required,
		},
	}




	app.Commands = []cli.Command{
		{
			Name: "generate",
			Aliases: []string{"g"},
			Action: func(c *cli.Context) error {
				tm.Println("Ok! Let's get it started, shall we?")
				tm.Print("This line should disappear")
				tm.ResetLine("Lol")

				i := &golive.App{}
				err := survey.Ask(qs, i)

				if err != nil {
					fmt.Println(err.Error())
				} else {
					if changeEnv() {
						fmt.Println(getEnvs())
					}

				}
				return nil
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

func changeEnv() bool {
	changeDefEnvQ := &survey.Confirm{
		Message: "By default, GoLive creates two environments for you: staging & production. Do you want to change that?",
	}

	changeEnv := false
	survey.AskOne(changeDefEnvQ, &changeEnv, nil)
	return changeEnv
}

func getEnvs() []string {
	commaString := "staging,production"
	survey.AskOne(&survey.Input{
		Message: "Type a list of environments (comma separated) that you wish to create",
	}, &commaString, nil)
	return util.SplitComma(commaString)
}