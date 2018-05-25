package golive

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"log"
	"errors"
	"github.com/paulovictorv/golive/app/infrastructure"
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
}


func check(err error) {
	if (err != nil) {
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

func CreateApp(appName string) (App, error) {
	app := App{
		Name: appName,
		Envs: []*Env{
			{
				Name:   "staging",
				Bucket: fmt.Sprintf("%s-%d-%s", appName, 1, "staging"),
				CdnId:  "aa",
			},
			{
				Name:   "production",
				Bucket: fmt.Sprintf("%s-%d-%s", appName, 1, "production"),
				CdnId:  "aa",
			},
		},
	}

	for _, e := range app.Envs {
		e.CdnId  = infrastructure.CreateEnv(e.Bucket)
	}

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
	check(e);

	defer file.Close();

	n, e := file.Write(out)

	if n > 0 {
		return n, nil;
	} else {
		return -1, errors.New("error while creating app")
	}
}