package golive

import (
	"testing"
	"path/filepath"
	"os"
	"strings"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
)

func TestWalkFolder(t *testing.T) {
	s := "../"

	var sess, _ = session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")})

	uploader := s3manager.NewUploader(sess)

	uploader.UploadWithIterator(aws.BackgroundContext(), )

	filepath.Walk(s, func(path string, info os.FileInfo, err error) error {
		effectiveName := strings.TrimPrefix(path, s)

		files := []string{}

		//ignoring "dotfiles"
		if !(strings.HasPrefix(effectiveName, ".")) && !(info.IsDir()) {
			files = append(files, effectiveName)
		}

		for _, file := range files {
			fmt.Println(file)
		}

		return nil
	})
}