package infrastructure

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/s3"
	"time"
)

var sess, _ = session.NewSession(&aws.Config{
	Region: aws.String("us-west-2")})

var s3Client = s3.New(sess)
var front = cloudfront.New(sess)

func CreateEnv(bucketName, domainName string, status chan string, complete chan int) string {
	status <- "BUCKET_START"
	//createBucket(bucketName)
	time.Sleep(1000)
	status <- "BUCKET_COMPLETE"

	status <- "CDN_START"
	//cdnId := createCdn(bucketName, domainName)
	time.Sleep(1000)
	status <- "CDN_COMPLETE"

	complete <- 1

	return "cdnId"
}

func parseAwsError(err error) string {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return aerr.Error()
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return err.Error()
		}
	}
	return ""
}
