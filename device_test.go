package resingo

import (
	"net/http"
	"testing"
)

func TestDevice(t *testing.T) {
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
	appName := "device_test"
	app, err := AppCreate(ctx, appName, RaspberryPi3)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = AppDelete(ctx, app.Name)
	}()
	maxDevices := 4
	devices := make([]struct {
		uuid string
		dev  *Device
	}, maxDevices)
	for i := 0; i < maxDevices; i++ {
		uid, err := GenerateUUID()
		if err != nil {
			t.Fatal(err)
		}
		devices[i].uuid = uid
	}
	t.Run("Register", func(ts *testing.T) {
		for _, d := range devices {
			testDevRegister(ctx, ts, appName, d.uuid)
		}
	})
}

func testDevGetAll(ctx *Context, t *testing.T, appName string, expect int) {
	dev, err := DevGetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(dev) != expect {
		t.Errorf("expected %d devices got %d ", expect, len(dev))
	}
}

func testDevRegister(ctx *Context, t *testing.T, appname, uuid string) {
	dev, err := DevRegister(ctx, appname, uuid)
	if err != nil {
		t.Error(err)
	}
	if dev != nil {
		if dev.UUID != uuid {
			t.Error("device uuid mismatch")
		}
	}
}
