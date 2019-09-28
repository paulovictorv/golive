package golive

import (
	"fmt"
	"goclip.com.br/golive/app/env"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"sync"
)

type App struct {
	Name              string     `yaml:"name"`
	Envs              []*env.Env `yaml:"envs"`
	OriginFolder      string     `yaml:"originFolder"`
	DestinationFolder string     `yaml:"destinationFolder"`
	InvalidationPaths []string   `yaml:"invalidationPaths"`
}

func InitApp(appName string) error {
	app := App{
		Name: appName,
		Envs: []*env.Env{
			{
				Name:   "production",
				Domain: "",
				CdnId:  "",
				Bucket: fmt.Sprintf("%s-%d-%s", appName, 1, "production"),
			},
			{
				Name:   "staging",
				Domain: "",
				CdnId:  "",
				Bucket: fmt.Sprintf("%s-%d-%s", appName, 1, "staging"),
			},
		},
		OriginFolder:      "dist/",
		DestinationFolder: "/",
		InvalidationPaths: []string{"/index.html"},
	}

	if err := saveFile(&app); err != nil {
		return err
	}

	return nil
}

func ProvisionApp(filePath string) error {
	bytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		return err
	}

	app := &App{}

	if err := yaml.Unmarshal(bytes, app); err != nil {
		return err
	}

	var waitGroup sync.WaitGroup

	//for _, env := range app.Envs {
	//	progress := make(chan string)
	//	complete := make(chan int)
	//
	//	waitGroup.Add(1)
	//
	//	infra := infrastructure.CreateInfra(infrastructure.AWS)
	//
	//	go infra.ProvisionEnv(env, progress, complete)
	//	go func() {
	//		for {
	//			select {
	//			case val := <-progress:
	//				fmt.Printf(val)
	//			case <-complete:
	//				waitGroup.Done()
	//				break
	//			}
	//		}
	//	}()
	//}

	waitGroup.Wait()
	return nil
}

func saveFile(app *App) error {
	out, err := yaml.Marshal(app)

	if err != nil {
		return err
	}

	file, err := os.Create(".golive.yml")

	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.Write(out); err != nil {
		return err
	} else {
		return nil
	}
}
