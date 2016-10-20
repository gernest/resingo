package resingo

import (
	"net/http"
	"os"
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
		_, _ = AppDelete(ctx, app.ID)
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
	t.Run("GetByUUID", func(ts *testing.T) {
		for i, d := range devices {
			devices[i].dev = testDevGetByUUID(ctx, ts, d.uuid)
		}
	})
	t.Run("GetByName", func(ts *testing.T) {
		for _, d := range devices {
			testDevGetByName(ctx, ts, d.dev.Name)
		}
	})
	t.Run("GetAllByApp", func(ts *testing.T) {
		testDevGetAllByApp(ctx, ts, app.ID, maxDevices)
	})
	t.Run("GetAll", func(ts *testing.T) {
		testDevGetAll(ctx, ts, appName, maxDevices)
	})
	t.Run("Rename", func(ts *testing.T) {
		testDevRename(ctx, ts, "avocado", devices[0].uuid)
	})
	t.Run("GetApp", func(ts *testing.T) {
		testDevGetApp(ctx, ts, devices[0].uuid, appName)
	})
	t.Run("EnableURL", func(ts *testing.T) {
		uuid := os.Getenv("RESINTEST_REALDEVICE_UUID")
		if uuid == "" {
			ts.Skip("missing RESINTEST_REALDEVICE_UUID")
		}
		testDevEnableURL(ctx, ts, devices[0].uuid)
	})
	t.Run("DisableURL", func(ts *testing.T) {
		uuid := os.Getenv("RESINTEST_REALDEVICE_UUID")
		if uuid == "" {
			ts.Skip("missing RESINTEST_REALDEVICE_UUID")
		}
		testDevDisableURL(ctx, ts, devices[0].uuid)
	})
	t.Run("Delete", func(ts *testing.T) {
		u, _ := GenerateUUID()
		d, err := DevRegister(ctx, appName, u)
		if err != nil {
			ts.Fatal(err)
		}
		err = DevDelete(ctx, d.ID)
		if err != nil {
			ts.Fatal(err)
		}
		_, err = DevGetByUUID(ctx, u)
		if err != ErrDeviceNotFound {
			t.Errorf("expected %s got %s", ErrDeviceNotFound.Error(), err.Error())
		}
	})
	t.Run("Note", func(ts *testing.T) {
		note := "hello,world"
		err := DevNote(ctx, devices[0].dev.ID, note)
		if err != nil {
			ts.Fatal(err)
		}
	})
	env := []struct {
		key, value string
	}{
		{"Mad", "Scientist"},
		{"MONIKER", "IOT"},
	}
	t.Run("CreateEnv", func(ts *testing.T) {
		for _, v := range devices {
			for _, e := range env {
				en, err := EnvDevCreate(ctx, v.dev.ID, e.key, e.value)
				if err != nil {
					ts.Error(err)
				}
				if en.Name != e.key {
					t.Errorf("expected %s got %s", e.key, en.Name)
				}
			}
		}
	})
	t.Run("EnvGetAll", func(ts *testing.T) {
		for _, v := range devices {
			envs, err := EnvDevGetAll(ctx, v.dev.ID)
			if err != nil {
				ts.Error(err)
			}
			if len(envs) != len(env) {
				t.Errorf("expected %d got %d", len(env), len(envs))
			}
		}
	})
}

func testDevGetAll(ctx *Context, t *testing.T, appName string, expect int) {
	dev, err := DevGetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(dev) < expect {
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

func testDevGetByUUID(ctx *Context, t *testing.T, uuid string) *Device {
	dev, err := DevGetByUUID(ctx, uuid)
	if err != nil {
		t.Fatal(err)
	}
	if dev.UUID != uuid {
		t.Fatalf("expected %s got %s", uuid, dev.UUID)
	}
	return dev
}
func testDevGetByName(ctx *Context, t *testing.T, name string) {
	dev, err := DevGetByName(ctx, name)
	if err != nil {
		t.Fatal(err)
	}
	if dev.Name != name {
		t.Errorf("expected %s got %s", name, dev.Name)
	}
}
func testDevGetAllByApp(ctx *Context, t *testing.T, appID int64, expect int) {
	dev, err := DevGetAllByApp(ctx, appID)
	if err != nil {
		t.Fatal(err)
	}
	if len(dev) != expect {
		t.Errorf("expected %d devies got %d", expect, len(dev))
	}
}

func testDevRename(ctx *Context, t *testing.T, newName string, uuid string) {
	err := DevRename(ctx, uuid, newName)
	if err != nil {
		t.Fatal(err)
	}
	dev, err := DevGetByUUID(ctx, uuid)
	if err != nil {
		t.Fatal(err)
	}
	if dev.Name != newName {
		t.Errorf("expected %s got %s", newName, dev.Name)
	}
}

func testDevGetApp(ctx *Context, t *testing.T, uuid, appName string) {
	app, err := DevGetApp(ctx, uuid)
	if err != nil {
		t.Fatal(err)
	}
	if app.Name != appName {
		t.Errorf("expected %s got %s", appName, app.Name)
	}

}

func testDevEnableURL(ctx *Context, t *testing.T, uuid string) {
	err := DevEnableURL(ctx, uuid)
	if err != nil {
		t.Fatal(err)
	}
	dev, err := DevGetByUUID(ctx, uuid)
	if err != nil {
		t.Fatal(err)
	}
	if !dev.WebAccessible {
		t.Error("the device should be web accessible")
	}
}
func testDevDisableURL(ctx *Context, t *testing.T, uuid string) {
	err := DevDisableURL(ctx, uuid)
	if err != nil {
		t.Fatal(err)
	}
	dev, err := DevGetByUUID(ctx, uuid)
	if err != nil {
		t.Fatal(err)
	}
	if dev.WebAccessible {
		t.Error("the device should not be web accessible")
	}
}
