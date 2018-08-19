package golive

import (
	"gopkg.in/yaml.v2"
	"os"
	"log"
	"errors"
	"github.com/paulovictorv/golive/app/infrastructure"
	"fmt"
	"io/ioutil"
	"strings"
)

type Env struct {
	Name   string `yaml:"name"`
	Domain string `yaml:"domain"`
	CdnId  string `yaml:"cdnId"`
	Bucket string `yaml:"bucket"`
}

type App struct {
	Name string `yaml:"name"`
	Envs []*Env `yaml:"envs"`
	OriginFolder string `yaml:"originFolder"`
	DestinationFolder string `yaml:"destinationFolder"`
	InvalidationPaths []string `yaml:"invalidationPaths"`
}

func check(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func CreateEnvs(appName string, envsNames []string) []*Env {
	var envs []*Env
	for _, envName := range envsNames {

		env := Env{
			Name: envName,
			Bucket: fmt.Sprintf("%s-%d-%s", appName, 1, envName),
		}

		envs = append(envs, &env)
	}
	return envs
}

func CreateApp(app App) (App, error) {
	for _, e := range app.Envs {
		e.CdnId = infrastructure.CreateEnv(e.Bucket, e.Domain)
	}

	_, e := saveFile(&app)

	if e != nil {
		return App{}, e
	} else {
		return app, nil
	}

}

func DeployApp(envName string) {
	bytes, e := ioutil.ReadFile(".golive.yml")

	if e != nil {
		panic(".golive.yml file not found")
	}

	app := &App{}

	yaml.Unmarshal(bytes, app)

	env, err := pickEnv(app, envName)

	if err != nil {
		panic(err)
	}

	infrastructure.UploadDir(app.OriginFolder, env.Bucket)
	infrastructure.InvalidateFiles(env.CdnId, app.InvalidationPaths)
}

func pickEnv(app *App, envName string) (*Env, error) {
	for _, env := range app.Envs {
		if strings.Compare(envName, env.Name) == 0 {
			return env, nil
		}
	}

	return nil, errors.New("Supplied name does not match to any env")
}

func saveFile(app *App) (int, error) {
	out, err := yaml.Marshal(app)
	check(err)

	file, e := os.Create(".golive.yml")
	check(e)

	defer file.Close()

	n, e := file.Write(out)

	if n > 0 {
		return n, nil
	} else {
		return -1, errors.New("error while creating app")
	}
}