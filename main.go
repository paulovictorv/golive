package main

import (
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
				paths := askInvalidationPaths()

				goliveApp := golive.App{
					Name:              appName,
					InvalidationPaths: paths,
					Envs:              initEnvs,
				}

				golive.CreateApp(goliveApp)

				return nil
			},
		},
		{
			Name: "deploy",
			Aliases: []string{"d"},
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func askInvalidationPaths() []string {
	tm.Print(tm.Bold("GoLive needs to submit a cache invalidation request when it deploys your files."))
	tm.Flush()
	invalidationPaths := ""
	survey.AskOne(&survey.Input{
		Message: "Please provide a list of files (comma separated) that GoLive needs to look for in order to invalidate",
	}, &invalidationPaths, survey.Required)

	return util.SplitComma(invalidationPaths);
}

func askDomainNames(envs []*golive.Env) {
	tm.Println(tm.Bold(tm.Color("Almost there!", tm.BLUE)))
	tm.Print(tm.Bold("Now, for each environment you will need to provide a domain name."))
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