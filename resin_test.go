package resingo

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

var ENV *EnvVars

type EnvVars struct {
	Username, Email, Password string
	ID                        int64
	Register                  struct {
		Username, Email, Password string
	}
}

func init() {
	ENV = &EnvVars{
		Username: os.Getenv("RESINTEST_USERNAME"),
		Password: os.Getenv("RESINTEST_PASSWORD"),
		Email:    os.Getenv("RESINTEST_EMAIL"),
	}
	ENV.Register.Username = os.Getenv("RESINTEST_REGISTER_USERNAME")
	ENV.Register.Password = os.Getenv("RESINTEST_REGISTER_PASSWORD")
	ENV.Register.Email = os.Getenv("RESINTEST_REGISTER_EMAIL")
}

func TestResin(t *testing.T) {
	config := &Config{
		Username:      ENV.Username,
		Password:      ENV.Password,
		ResinEndpoint: apiEndpoint,
	}
	client := &http.Client{}
	t.Run("Authenticate", func(ts *testing.T) {
		cfg := *config
		ctx := &Context{
			Client: client,
			Config: &cfg,
		}
		testAuthenticate(ctx, ts)
	})
	t.Run("Login", func(ts *testing.T) {
		cfg := *config
		ctx := &Context{
			Client: client,
			Config: &cfg,
		}
		testLogin(ctx, ts)
	})
}

func testAuthenticate(ctx *Context, t *testing.T) {
	token, err := Authenticate(ctx, Credentials)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("ParseToken", func(ts *testing.T) {
		claims, err := ParseToken(token)
		if err != nil {
			ts.Fatal(err)
		}
		if claims.Username != ctx.Config.Username {
			ts.Errorf("expected username %s got %s", ctx.Config.Username, claims.Username)
		}

	})
}

func testLogin(ctx *Context, t *testing.T) {
	err := Login(ctx, Credentials)
	if err != nil {
		t.Fatal(err)
	}
	if ctx.Config.tokenClain == nil {
		t.Error("expected the token to be saved")
	}
}

func TestEncode(t *testing.T) {
	sample := []struct {
		params []string
		expect string
	}{
		{[]string{"filter,Name", "eq,Milk"}, "$filter=Name%20eq%20'Milk'"},
		{[]string{"expand,device"}, "$expand=device"},
	}

	for _, v := range sample {
		param := make(url.Values)
		for _, p := range v.params {
			s := strings.Split(p, ",")
			param.Set(s[0], s[1])
		}
		e := Encode(param)
		if e != v.expect {
			t.Errorf("expectes %s got %s", v.expect, e)
		}
	}
}
