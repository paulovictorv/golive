package infrastructure

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"strconv"
	"time"
)

func createCdn(bucketName, domainName string) string {
	itoa := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

	output, cdnError := front.CreateDistribution(&cloudfront.CreateDistributionInput{
		DistributionConfig: &cloudfront.DistributionConfig{
			Aliases: &cloudfront.Aliases{
				Items:    aws.StringSlice([]string{domainName}),
				Quantity: aws.Int64(1),
			},
			CallerReference: &itoa,
			Comment:         aws.String("Distribution created by GoLive"),
			DefaultCacheBehavior: &cloudfront.DefaultCacheBehavior{
				ForwardedValues: &cloudfront.ForwardedValues{
					Cookies: &cloudfront.CookiePreference{
						Forward: aws.String("none"),
					},
					QueryString: aws.Bool(false),
				},
				MinTTL:         aws.Int64(0),
				MaxTTL:         aws.Int64(31536000),
				DefaultTTL:     aws.Int64(86400),
				TargetOriginId: aws.String(bucketName),
				TrustedSigners: &cloudfront.TrustedSigners{
					Enabled:  aws.Bool(false),
					Quantity: aws.Int64(0),
				},
				ViewerProtocolPolicy: aws.String(cloudfront.ViewerProtocolPolicyAllowAll),
				AllowedMethods: &cloudfront.AllowedMethods{
					Items:    aws.StringSlice([]string{"GET", "HEAD"}),
					Quantity: aws.Int64(2),
				},
			},
			DefaultRootObject: aws.String("index.html"),
			Enabled:           aws.Bool(true),
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
			ViewerCertificate: &cloudfront.ViewerCertificate{
				ACMCertificateArn:            nil,
				CloudFrontDefaultCertificate: nil,
				IAMCertificateId:             nil,
				MinimumProtocolVersion:       aws.String("TLSv1.1_2016"),
				SSLSupportMethod:             aws.String("sni-only"),
			},
		},
	})

	parseAwsError(cdnError)

	return *output.Distribution.Id
}
