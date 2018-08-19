package infrastructure

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"strconv"
	"time"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"path/filepath"
	"os"
	"strings"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/briandowns/spinner"
	"mime"
)

var sess, _ = session.NewSession(&aws.Config{
Region: aws.String("us-west-2")})

var s3Client = s3.New(sess)
var front = cloudfront.New(sess)

func UploadDir(rootPath string, bucketName string) {
	var files []s3manager.BatchUploadObject

	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		key := strings.TrimPrefix(path, rootPath)

		if !strings.HasPrefix(key, ".") && !info.IsDir() {
			body, err := os.Open(path)

			if err != nil {

			}

			uploadObject := s3manager.BatchUploadObject{
				Object: &s3manager.UploadInput{
					Bucket: aws.String(bucketName),
					Key:    aws.String(key),
					ContentType: aws.String(mime.TypeByExtension(filepath.Ext(key))),
					Body:   body,
				},
				After: func() error {
					fmt.Printf("%s uploaded\n", key)
					return nil
				},
			}

			files = append(files, uploadObject)
		}

		return nil
	})

	uploader := s3manager.NewUploader(sess)

	iterator := s3manager.UploadObjectsIterator{
		Objects: files,
	}

	uploader.UploadWithIterator(aws.BackgroundContext(), &iterator)
}

func InvalidateFiles(id string, files []string) (string, error) {
	itoa := strconv.FormatInt(time.Now().UnixNano() / int64(time.Millisecond), 10);

	output, e := front.CreateInvalidation(&cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(id),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: &itoa,
			Paths: &cloudfront.Paths{
				Quantity: aws.Int64(int64(len(files))),
				Items:    aws.StringSlice(files),
			},
		},
	})

	if e != nil {
		return "", e
	}

	return *output.Invalidation.Id, nil
}

func CreateEnv(bucketName, domainName string) string {
	s := spinner.New(spinner.CharSets[26], 100 * time.Millisecond)

	s.Prefix = "Creating Bucket"
	s.Start()
	createBucket(bucketName)
	s.FinalMSG = "Done!"
	s.Stop()

	fmt.Println()

	s.Prefix = "Creating Cloudfront Distribution"
	s.Start()
	cdnId := createCdn(bucketName, domainName)
	s.FinalMSG = "Done!"
	s.Stop()

	fmt.Println()

	return cdnId
}

func createBucket(bucketName string) {
	_, bucketError := s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	parseAwsError(bucketError)

	publicPolicy := fmt.Sprintf("{\"Version\":\"2012-10-17\",\"Statement\": [{\"Sid\": \"AddPerm\",\"Effect\": \"Allow\"," +
		"\"Principal\": \"*\"," +
		"\"Action\": \"s3:GetObject\"," +
		"\"Resource\": \"arn:aws:s3:::%s/*\"}]}", bucketName)

	_, policyError := s3Client.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(publicPolicy),
	})

	parseAwsError(policyError)
}

func createCdn(bucketName, domainName string) string {
	itoa := strconv.FormatInt(time.Now().UnixNano() / int64(time.Millisecond), 10);

	output, cdnError := front.CreateDistribution(&cloudfront.CreateDistributionInput{
		DistributionConfig: &cloudfront.DistributionConfig{
			Enabled: aws.Bool(true),
			Aliases: &cloudfront.Aliases{
				Items: aws.StringSlice([]string {domainName}),
				Quantity: aws.Int64(1),
			},
			DefaultRootObject: aws.String("index.html"),
			DefaultCacheBehavior: &cloudfront.DefaultCacheBehavior{
				ForwardedValues: &cloudfront.ForwardedValues{
					Cookies: &cloudfront.CookiePreference{
						Forward: aws.String("none"),
					},
					QueryString: aws.Bool(false),
				},
				MinTTL: aws.Int64(0),
				MaxTTL: aws.Int64(31536000),
				DefaultTTL: aws.Int64(86400),
				TargetOriginId: aws.String(bucketName),
				TrustedSigners: &cloudfront.TrustedSigners{
					Enabled: aws.Bool(false),
					Quantity: aws.Int64(0),
				},
				ViewerProtocolPolicy: aws.String(cloudfront.ViewerProtocolPolicyAllowAll),
				AllowedMethods: &cloudfront.AllowedMethods{
					Items: aws.StringSlice([]string{"GET", "HEAD"}),
					Quantity: aws.Int64(2),
				},
			},
			CallerReference: &itoa,
			Comment: aws.String("Distribution created by GoLive"),
			Origins: &cloudfront.Origins{
				Items: []*cloudfront.Origin{
					{
						Id:         aws.String(bucketName),
						DomainName: aws.String(fmt.Sprintf("%s.s3.amazonaws.com", bucketName)),
						S3OriginConfig: &cloudfront.S3OriginConfig{
							OriginAccessIdentity: aws.String(""),
						},
					},
				},

				Quantity: aws.Int64(1),
			},
		},
	})

	parseAwsError(cdnError)

	return *output.Distribution.Id
}

func parseAwsError(err error) {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
}
