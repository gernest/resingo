package resingo

import (
	"net/http"
	"testing"
)

func TestResinConfig(t *testing.T) {
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
	_, err = ConfigGetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
