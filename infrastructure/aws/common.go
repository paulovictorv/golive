package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
)

var sess, _ = session.NewSession(&aws.Config{
	Region: aws.String("us-east-1")})

var s3Client = s3.New(sess)
var front = cloudfront.New(sess)
var route53Client = route53.New(sess)

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
