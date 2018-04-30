package infrastructure

import (
	"testing"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"fmt"
	"time"
	"strconv"
)

func TestCreateCdn(t *testing.T) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")})

	front := cloudfront.New(sess)

	itoa := strconv.FormatInt(time.Now().UnixNano() / int64(time.Millisecond), 10);

	output, e := front.CreateDistribution(&cloudfront.CreateDistributionInput{
		DistributionConfig: &cloudfront.DistributionConfig{
			Enabled: aws.Bool(true),
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
				TargetOriginId: aws.String("golive-test-bucket-1"),
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
						Id:         aws.String("golive-test-bucket-1"),
						DomainName: aws.String("golive-test-bucket-1.s3.amazonaws.com"),
						S3OriginConfig: &cloudfront.S3OriginConfig{
							OriginAccessIdentity: aws.String(""),
						},
					},
				},

				Quantity: aws.Int64(1),
			},
		},
	})

	if e != nil {
		if aerr, ok := e.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(e.Error())
		}
	}

	fmt.Println(output.Distribution.Id)
}
