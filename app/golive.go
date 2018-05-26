package golive

import (
	"gopkg.in/yaml.v2"
	"os"
	"log"
	"errors"
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
	InvalidationPaths []string `yaml:"invalidationPaths"`
}

func check(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func CreateEnvs(envsNames []string) []*Env {
	var envs []*Env
	for _, envName := range envsNames {
		envs = append(envs, &Env{Name: envName})
	}
	return envs
}

func CreateApp(app App) (App, error) {
	//for _, e := range app.Envs {
	//	e.CdnId  = infrastructure.CreateEnv(e.Bucket)
	//}

	_, e := saveFile(&app)

	if e != nil {
		return App{}, e
	} else {
		return app, nil
	}

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