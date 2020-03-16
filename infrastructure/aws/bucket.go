package aws

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func createBucket(bucketName string) (string, error) {
	bucketNameReference := aws.String(bucketName)

	_, bucketError := s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: bucketNameReference,
	})

	if bucketError != nil {
		return "", errors.New(parseAwsError(bucketError))
	}

	publicPolicy := fmt.Sprintf("{\"Version\":\"2012-10-17\",\"Statement\": [{\"Sid\": \"AddPerm\",\"Effect\": \"Allow\","+
		"\"Principal\": \"*\","+
		"\"Action\": \"s3:GetObject\","+
		"\"Resource\": \"arn:aws:s3:::%s/*\"}]}", bucketName)

	_, policyError := s3Client.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: bucketNameReference,
		Policy: aws.String(publicPolicy),
	})

	if policyError != nil {
		return "", errors.New(parseAwsError(policyError))
	}

	_, webSiteError := s3Client.PutBucketWebsite(&s3.PutBucketWebsiteInput{
		Bucket: bucketNameReference,
		WebsiteConfiguration: &s3.WebsiteConfiguration{
			IndexDocument: &s3.IndexDocument{Suffix: aws.String("index.html")},
		},
	})

	if webSiteError != nil {
		return "", errors.New(parseAwsError(webSiteError))
	}

	return fmt.Sprintf("%s.s3-website-%s.amazonaws.com", bucketName, *sess.Config.Region), nil
}
