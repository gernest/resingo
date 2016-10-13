package resingo

import (
	"fmt"
	"net/http"
	"testing"
)

func TestApplication(t *testing.T) {
	config := &Config{
		Username:      ENV.Username,
		Password:      ENV.Password,
		ResinEndpoint: apiEndpoint,
	}
	client := &http.Client{}
	ctx := &Context{
		Client: client,
		Config: config,
	}
	err := Login(ctx, Credentials)
	if err != nil {
		t.Fatal(err)
	}
	applications := []struct {
		name string
		app  *Application
		typ  DeviceType
	}{
		{"resingo", nil, RaspberryPi3},
		{"algorithm_zero", nil, RaspberryPi2},
	}
	for i, a := range applications {
		_, _ = AppDelete(ctx, a.name)
		app, err := AppCreate(ctx, a.name, a.typ)
		if err != nil {
			t.Fatal(err)
		}
		applications[i].app = app
	}
	defer func() {
		for _, a := range applications {
			_, _ = AppDelete(ctx, a.name)
		}
	}()

	t.Run("AppGetAll", func(ts *testing.T) {
		testAppGetAll(ctx, ts)
	})
	t.Run("AppGetByName", func(ts *testing.T) {
		for _, a := range applications {
			testApGetByName(ctx, ts, a.name)
		}
	})
	t.Run("AppGetByID", func(ts *testing.T) {
		for _, a := range applications {
			testApGetByID(ctx, ts, a.app.ID)
		}
	})
	t.Run("AppCreate", func(ts *testing.T) {
		testAppCreate(ctx, ts, "resingo_test", RaspberryPi3)
	})
	t.Run("GetApiKey", func(ts *testing.T) {
		for _, a := range applications {
			testAppAPIKey(ctx, ts, a.name)
		}
	})
}

func testAppGetAll(ctx *Context, t *testing.T) {
	apps, err := AppGetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) == 0 {
		t.Fatal("expected at least more than one application")
	}
}

func testApGetByName(ctx *Context, t *testing.T, name string) {
	app, err := AppGetByName(ctx, name)
	if err != nil {
		t.Fatal(err)
	}
	if app.Name != name {
		t.Errorf("expected %s got %s %v", name, app.Name, *app)
	}
}
func testApGetByID(ctx *Context, t *testing.T, id int64) {
	app, err := AppGetByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if app.ID != id {
		t.Errorf("expected %d got %d", id, app.ID)
	}
}

func testAppCreate(ctx *Context, t *testing.T, name string, typ DeviceType) {
	app, err := AppCreate(ctx, name, typ)
	if err != nil {
		t.Fatal(err)
	}
	if app.Name != name {
		t.Fatalf("expected %s got %s", name, app.Name)
	}
	t.Run("Delete", func(ts *testing.T) {
		testAppDelete(ctx, ts, name)
	})
}

func testAppDelete(ctx *Context, t *testing.T, name string) {
	_, err := AppDelete(ctx, name)
	if err != nil {
		t.Fatal(err)
	}
	_, err = AppGetByName(ctx, name)
	if err == nil {
		t.Error("expected devcice not found error")
	}
}
func testAppAPIKey(ctx *Context, t *testing.T, name string) {
	b, err := AppGetAPIKey(ctx, name)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}
