package infrastructure

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

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
					Bucket:      aws.String(bucketName),
					Key:         aws.String(key),
					ContentType: aws.String(mime.TypeByExtension(filepath.Ext(key))),
					Body:        body,
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
	itoa := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

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
		return "", errors.New(parseAwsError(e))
	}

	return *output.Invalidation.Id, nil
}
