package golive

import "testing"

func TestCreateApp(t *testing.T) {
	app, e := CreateApp("taskingio-3")

	if e != nil {
		t.Fail()
	} else if len(app.Name) == 0 {
		t.Errorf("expecting %s to equal %s", app.Name, "taskingio-2")
	}

}