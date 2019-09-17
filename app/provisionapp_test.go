package golive

import (
	"fmt"
	"testing"
)

func TestProvisionApp(t *testing.T) {
	error := ProvisionApp(".golive.yml")

	fmt.Println(error)
}
