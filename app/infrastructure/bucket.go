package infrastructure

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func createBucket(bucketName string) {
	_, bucketError := s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	parseAwsError(bucketError)

	publicPolicy := fmt.Sprintf("{\"Version\":\"2012-10-17\",\"Statement\": [{\"Sid\": \"AddPerm\",\"Effect\": \"Allow\","+
		"\"Principal\": \"*\","+
		"\"Action\": \"s3:GetObject\","+
		"\"Resource\": \"arn:aws:s3:::%s/*\"}]}", bucketName)

	_, policyError := s3Client.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(publicPolicy),
	})

	parseAwsError(policyError)
}
