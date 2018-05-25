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

	app.Commands = []cli.Command{
		{
			Name: "generate",
			Aliases: []string{"g"},
			Action: func(c *cli.Context) error {
				tm.Println(tm.Bold("Ok! Let's get it started, shall we?"))
				tm.Flush()

				appName := ""
				survey.AskOne(&survey.Input{Message:"What's the name of your app?"}, &appName, survey.Required)

				envs := []string {"production", "staging"}

				if askChangeEnvs() {
					envs = askNewEnvs()
				}

				initEnvs := golive.CreateEnvs(envs)

				askDomainNames(initEnvs)


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

func askDomainNames(envs []*golive.Env) {
	tm.Println(tm.Bold("Now, for each environment you will need to provide a domain name."))
	tm.Flush()

	for _, env := range envs {
		domainName := ""
		survey.AskOne(&survey.Input{
			Message: fmt.Sprintf("Specify domain name for %s environment", env.Name),
		}, &domainName, nil)
		env.Domain = domainName
	}
}

func askChangeEnvs() bool {
	changeDefEnvQ := &survey.Confirm{
		Message: "By default, GoLive creates two environments for you: staging & production. Do you want to change that?",
	}

	changeEnv := false
	survey.AskOne(changeDefEnvQ, &changeEnv, nil)
	return changeEnv
}

func askNewEnvs() []string {
	commaString := "staging,production"
	survey.AskOne(&survey.Input{
		Message: "Type a list of environments (comma separated) that you wish to create",
	}, &commaString, nil)
	return util.SplitComma(commaString)
}