package infrastructure

import (
	"goclip.com.br/golive/app/env"
	"goclip.com.br/golive/infrastructure/aws"
)

type Infrastructure interface {
	ProvisionEnv(env *env.Env, status chan string, complete chan int)
	DeployEnv(env env.Env)
}

type Provider int

const (
	AWS  Provider = 0
	STUB Provider = 1
)

func CreateInfra(provider Provider) Infrastructure {
	switch provider {
	case AWS:
		return aws.AmazonInfrastructure{}
	default:
		return nil
	}
}
