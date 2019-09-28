package main

import (
	"fmt"
	"goclip.com.br/golive/infrastructure"
	"sync"
)

func main() {

	//app := &cli.App{
	//	Name:  "golive",
	//	Usage: "Automatic provisioning of static sites deployments to AWS using CDN, ACM, S3 and Route53",
	//	Commands: [] *cli.Command{
	//		{
	//			Name:    "init",
	//			Aliases: []string{"i"},
	//			Usage:   "scaffolds the initial file for a project",
	//			Action: func(context *cli.Context) error {
	//				//name
	//				if err := golive.InitApp(context.Args().Get(0)); err != nil {
	//					log.Fatalf("error: %v", err)
	//				}
	//
	//				return nil
	//			},
	//			Flags: []cli.Flag{
	//
	//			},
	//		},
	//		{
	//			Name:    "provision",
	//			Aliases: []string{"p"},
	//			Usage:   "provisions the full environment validating the input file for the minimal requirements",
	//			Action: func(context *cli.Context) error {
	//				filePath := context.Args().Get(0)
	//
	//				var err error
	//
	//				if len(filePath) == 0 {
	//					err = golive.ProvisionApp(".golive.yml")
	//				} else {
	//					err = golive.ProvisionApp(filePath)
	//				}
	//
	//				if err != nil {
	//					log.Fatalf("error: %v", err)
	//				}
	//
	//				return nil
	//			},
	//			Flags: []cli.Flag{
	//
	//			},
	//		},
	//	},
	//}
	//
	//_ = app.Run(os.Args)
	var waitGroup sync.WaitGroup
	progress := make(chan string)
	complete := make(chan int)

	waitGroup.Add(1)

	go infrastructure.CreateEnv("env.Bucket", "env.Domain", progress, complete)
	go func() {
		for {
			select {
			case val := <-progress:
				fmt.Printf(val)
			case <-complete:
				waitGroup.Done()
				break
			}
		}
	}()

	waitGroup.Wait()
}
