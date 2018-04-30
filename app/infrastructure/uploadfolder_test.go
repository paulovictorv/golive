package infrastructure

import (
	"testing"
)

func TestUploadFolder(t *testing.T) {
	UploadDir("./", "golive-test-bucket-1")
}