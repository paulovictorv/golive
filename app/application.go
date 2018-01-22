package application

import "fmt"

type Env struct {
	Name   string
	CdnId  string
	Bucket string
}

type App struct {
	Name string
	Envs []Env
}

func CreateApp(appName string) App {
	app := App{
		Name: appName,
		Envs: []Env{
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

	return app
}
