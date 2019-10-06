package aws

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"strconv"
	"time"
)

func createCdn(bucketName, domainName, acmArn string, autoEnable bool) (string, error) {
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
			CustomErrorResponses: &cloudfront.CustomErrorResponses{
				Items: []*cloudfront.CustomErrorResponse{
					{
						ErrorCachingMinTTL: aws.Int64(0),
						ErrorCode:          aws.Int64(403),
						ResponseCode:       aws.String("200"),
						ResponsePagePath:   aws.String("/index.html"),
					},
					{
						ErrorCachingMinTTL: aws.Int64(0),
						ErrorCode:          aws.Int64(404),
						ResponseCode:       aws.String("200"),
						ResponsePagePath:   aws.String("/index.html"),
					},
				},
				Quantity: aws.Int64(2),
			},
			Enabled: aws.Bool(autoEnable),
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
				ACMCertificateArn:      aws.String(acmArn),
				MinimumProtocolVersion: aws.String("TLSv1.1_2016"),
				SSLSupportMethod:       aws.String("sni-only"),
			},
		},
	})

	if cdnError != nil {
		return "", errors.New(parseAwsError(cdnError))
	}

	return *output.Distribution.Id, nil
}
