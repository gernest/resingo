package resingo

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestKeys(t *testing.T) {
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
	sample := []struct {
		title, file string
		key         *Key
	}{
		{"testKey", "fixture/id_rsa.pub", nil},
	}
	t.Run("Create", func(ts *testing.T) {
		for i, v := range sample {
			d, err := ioutil.ReadFile(v.file)
			if err != nil {
				ts.Fatal(err)
			}
			k, err := KeyCreate(ctx, ctx.Config.UserID(), string(d), v.title)
			if err != nil {
				ts.Fatal(err)
			}
			sample[i].key = k
		}
	})
	t.Run("GetByID", func(ts *testing.T) {
		for _, v := range sample {
			k, err := KeyGetByID(ctx, v.key.ID)
			if err != nil {
				ts.Fatal(err)
			}
			if k.Title != v.title {
				ts.Errorf("expected %s got %s", v.title, k.Title)
			}
		}
	})
	t.Run("GetAll", func(ts *testing.T) {
		keys, err := KeyGetAll(ctx)
		if err != nil {
			ts.Fatal(err)
		}
		if len(keys) < 1 {
			ts.Error("expected at least one key")
		}
	})
	t.Run("Remove", func(ts *testing.T) {
		for _, v := range sample {
			err := KeyRemove(ctx, v.key.ID)
			if err != nil {
				ts.Fatal(err)
			}
		}
	})
}
