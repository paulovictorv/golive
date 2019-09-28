package infrastructure

import (
	"goclip.com.br/golive/app/env"
)

type Infrastructure interface {
	ProvisionEnv(env *env.Env, status chan string, complete chan int)
	DeployEnv(env env.Env)
}

type Provider int

const (
	AWS Provider = 0
)

func CreateInfra(provider Provider) Infrastructure {
	switch provider {
	case AWS:
		return AmazonInfrastructure{}
	default:
		return nil
	}
}
