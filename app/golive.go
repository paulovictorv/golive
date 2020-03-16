package golive

import (
	"fmt"
	"goclip.com.br/golive/app/env"
	"goclip.com.br/golive/infrastructure"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"sync"
)

type App struct {
	Name              string                  `yaml:"name"`
	Provider          infrastructure.Provider `yaml:"provider"`
	Envs              []*env.Env              `yaml:"envs"`
	OriginFolder      string                  `yaml:"originFolder"`
	DestinationFolder string                  `yaml:"destinationFolder"`
	InvalidationPaths []string                `yaml:"invalidationPaths"`
}

func createApp(appName string) App {
	return App{
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
}

func InitApp(appName string) (App, error) {
	app := createApp(appName)

	if err := saveFile(&app); err != nil {
		return App{}, err
	}

	return app, nil
}

func LoadApp(filePath string) (App, error) {
	bytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		return App{}, err
	}

	app := &App{}

	if err := yaml.Unmarshal(bytes, app); err != nil {
		return App{}, err
	}

	return *app, nil
}

func ProvisionApp(app App) error {
	var waitGroup sync.WaitGroup

	for _, appEnv := range app.Envs {
		progress := make(chan string)
		complete := make(chan int)

		waitGroup.Add(1)

		infra := infrastructure.CreateInfra(app.Provider)

		go infra.ProvisionEnv(appEnv, progress, complete)
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
	}

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
