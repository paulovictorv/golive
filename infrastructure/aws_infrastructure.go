package infrastructure

import (
	"goclip.com.br/golive/app/env"
	"time"
)

type AmazonInfrastructure struct {
}

func (i AmazonInfrastructure) ProvisionEnv(env *env.Env, status chan string, complete chan int) {
	//1. Check if domain has an certificate provisioned
	//1.1 if it does, get the ARN for the certificate [X]
	//1.2 if it doesn't, request a new certificate. Use DNS validation by default. [X]
	//1.2.1 if domain has a hosted zone on Route53, register DNS entries and wait for provisioning [X]
	//1.2.2 if not, print out DNS entries before continuing
	//2. create a bucket (PUBLIC_READ, hosting static website). Store the bucket website address
	//3. create a CDN
	// 3.1 point origin to the bucket website hosting URL
	// 3.2 use certificate obtained from step 1
	// 3.3 setup default error pages for 403, 401 and 404
	//4. prepare Route53 configuration
	// 4.1 if it's a root domain, register as alias
	// 4.2 if it's not, register CNAME register
	//env is provisioned

	status <- "BUCKET_START"
	//createBucket(bucketName)
	time.Sleep(1000)
	status <- "BUCKET_COMPLETE"

	status <- "CDN_START"
	//cdnId := createCdn(bucketName, domainName)
	time.Sleep(1000)
	status <- "CDN_COMPLETE"

	complete <- 1

	env.CdnId = ""
}

func (i AmazonInfrastructure) DeployEnv(env env.Env) {
	panic("implement me")
}
