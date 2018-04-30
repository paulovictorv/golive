package infrastructure

import (
	"testing"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

func TestCreateBucket(t *testing.T) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")})

	s3Client := s3.New(sess)

	bucketName := "golive-test-bucket-1"

	output, error := s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if error != nil {
		if aerr, ok := error.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(error.Error())
		}
	}

	publicPolicy := fmt.Sprintf("{\"Version\":\"2012-10-17\",\"Statement\": [{\"Sid\": \"AddPerm\",\"Effect\": \"Allow\"," +
		"\"Principal\": \"*\"," +
		"\"Action\": \"s3:GetObject\"," +
		"\"Resource\": \"arn:aws:s3:::%s/*\"}]}", bucketName)

	policyOutput, err := s3Client.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(publicPolicy),
	})

	if err != nil {
		if aerr, ok := error.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(error.Error())
		}
	}

	fmt.Println(output.Location)
	fmt.Println(policyOutput.String())
}