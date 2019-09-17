package golive

import (
	"errors"
	"fmt"
	"goclip.com.br/golive/app/infrastructure"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type Env struct {
	Name   string `yaml:"name"`
	Domain string `yaml:"domain"`
	CdnId  string `yaml:"cdnId"`
	Bucket string `yaml:"bucket"`
}

type App struct {
	Name              string   `yaml:"name"`
	Envs              []*Env   `yaml:"envs"`
	OriginFolder      string   `yaml:"originFolder"`
	DestinationFolder string   `yaml:"destinationFolder"`
	InvalidationPaths []string `yaml:"invalidationPaths"`
}

func InitApp(appName string) error {
	app := App{
		Name: appName,
		Envs: []*Env{
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

	for _, env := range app.Envs {
		progress := make(chan string)
		complete := make(chan int)

		waitGroup.Add(1)

		go infrastructure.CreateEnv(env.Bucket, env.Domain, progress, complete)
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

func DeployApp(envName string) {
	bytes, e := ioutil.ReadFile(".golive.yml")

	if e != nil {
		panic(".golive.yml file not found")
	}

	app := &App{}

	if err := yaml.Unmarshal(bytes, app); err != nil {
		return
	}

	env, err := pickEnv(app, envName)

	if err != nil {
		panic(err)
	}

	infrastructure.UploadDir(app.OriginFolder, env.Bucket)
	_, invErr := infrastructure.InvalidateFiles(env.CdnId, app.InvalidationPaths)

	if invErr != nil {
		panic(invErr)
	}

}

func pickEnv(app *App, envName string) (*Env, error) {
	for _, env := range app.Envs {
		if strings.Compare(envName, env.Name) == 0 {
			return env, nil
		}
	}

	return nil, errors.New("Supplied name does not match to any env")
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
