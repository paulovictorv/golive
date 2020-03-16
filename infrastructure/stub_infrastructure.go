package infrastructure

import (
	"goclip.com.br/golive/app/env"
	"time"
)

type StubInfrastructure struct {
}

func (s StubInfrastructure) ProvisionEnv(env *env.Env, status chan string, complete chan int) {
	status <- "BUCKET_START"
	//createBucket(bucketName)
	time.Sleep(10 * time.Second)
	status <- "BUCKET_COMPLETE"

	status <- "CDN_START"
	//cdnId := createCdn(bucketName, domainName)
	time.Sleep(1000)
	status <- "CDN_COMPLETE"

	complete <- 1

	env.CdnId = ""
}

func (s StubInfrastructure) DeployEnv(env env.Env) {
	panic("implement me")
}
