package golive

import (
	"goclip.com.br/golive/infrastructure"
	"testing"
)

func TestProvisionApp(t *testing.T) {
	app := createApp("stub")
	app.Provider = infrastructure.STUB

	ProvisionApp(app)
}
