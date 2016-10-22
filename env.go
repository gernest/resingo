package resingo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

//Env contains the response for device environment variable
type Env struct {
	ID     int64  `json:"id"`
	Name   string `json:"env_var_name"`
	Value  string `json:"value"`
	Device struct {
		ID       int64 `json:"__id"`
		Deferred struct {
			URI string `json:"uri"`
		} `json:"__deferred"`
	} `json:"device"`
	Metadata struct {
		URI  string `json:"uri"`
		Type string `json:"type"`
	} `json:"__metadata"`
}

//AppEnv application environment variable
type AppEnv struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Application struct {
		ID       int64 `json:"__id"`
		Deferred struct {
			URI string `json:"uri"`
		} `json:"__deferred"`
	} `json:"application"`
	Metadata struct {
		URI  string `json:"uri"`
		Type string `json:"type"`
	} `json:"__metadata"`
}

//EnvDevCreate creates environment variable for the device
func EnvDevCreate(ctx *Context, id int64, key, value string) (*Env, error) {
	h := authHeader(ctx.Config.AuthToken)
	uri := ctx.Config.APIEndpoint("device_environment_variable")
	data := make(map[string]interface{})
	data["device"] = id
	data["env_var_name"] = key
	data["value"] = value
	body, err := marhsalReader(data)
	if err != nil {
		return nil, err
	}
	b, err := doJSON(ctx, "POST", uri, h, nil, body)
	if err != nil {
		return nil, err
	}
	e := &Env{}
	err = json.Unmarshal(b, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

//EnvDevGetAll get all environment variables for the device
func EnvDevGetAll(ctx *Context, id int64) ([]*Env, error) {
	h := authHeader(ctx.Config.AuthToken)
	uri := ctx.Config.APIEndpoint("device_environment_variable")
	param := make(url.Values)
	param.Set("filter", "device")
	param.Set("eq", fmt.Sprint(id))
	b, err := doJSON(ctx, "GET", uri, h, param, nil)
	if err != nil {
		return nil, err
	}
	res := struct {
		D []*Env `json:"d"`
	}{}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}
	return res.D, nil
}

//EnvDevUpdate updates environment variable for device. The id is for  the
//environmant variable.
func EnvDevUpdate(ctx *Context, id int64, value string) error {
	h := authHeader(ctx.Config.AuthToken)
	s := fmt.Sprintf("device_environment_variable(%d)", id)
	uri := ctx.Config.APIEndpoint(s)
	data := make(map[string]interface{})
	data["value"] = value
	body, err := marhsalReader(data)
	if err != nil {
		return err
	}
	b, err := doJSON(ctx, "PATCH", uri, h, nil, body)
	if err != nil {
		return err
	}
	if string(b) != "OK" {
		return errors.New("bad response")
	}
	return nil
}

//EnvDevDelete deketes device environment variable
func EnvDevDelete(ctx *Context, id int64) error {
	h := authHeader(ctx.Config.AuthToken)
	s := fmt.Sprintf("device_environment_variable(%d)", id)
	uri := ctx.Config.APIEndpoint(s)
	b, err := doJSON(ctx, "DELETE", uri, h, nil, nil)
	if err != nil {
		return err
	}
	if string(b) != "OK" {
		return errors.New("bad response")
	}
	return nil
}

//EnvAppGetAll retruns all environment variables for application
func EnvAppGetAll(ctx *Context, id int64) ([]*AppEnv, error) {
	h := authHeader(ctx.Config.AuthToken)
	uri := ctx.Config.APIEndpoint("environment_variable")
	param := make(url.Values)
	param.Set("filter", "application")
	param.Set("eq", fmt.Sprint(id))
	b, err := doJSON(ctx, "GET", uri, h, param, nil)
	if err != nil {
		return nil, err
	}
	res := struct {
		D []*AppEnv `json:"d"`
	}{}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}
	return res.D, nil
}

//EnvAppCreate creates a newapplication environment variable
func EnvAppCreate(ctx *Context, id int64, key, value string) (*AppEnv, error) {
	h := authHeader(ctx.Config.AuthToken)
	uri := ctx.Config.APIEndpoint("environment_variable")
	data := make(map[string]interface{})
	data["application"] = id
	data["name"] = key
	data["value"] = value
	body, err := marhsalReader(data)
	if err != nil {
		return nil, err
	}
	b, err := doJSON(ctx, "POST", uri, h, nil, body)
	if err != nil {
		return nil, err
	}
	e := &AppEnv{}
	err = json.Unmarshal(b, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

//EnvAppUpdate updates an existing application environmant variable
func EnvAppUpdate(ctx *Context, id int64, value string) error {
	h := authHeader(ctx.Config.AuthToken)
	s := fmt.Sprintf("environment_variable(%d)", id)
	uri := ctx.Config.APIEndpoint(s)
	data := make(map[string]interface{})
	data["value"] = value
	body, err := marhsalReader(data)
	if err != nil {
		return err
	}
	b, err := doJSON(ctx, "PATCH", uri, h, nil, body)
	if err != nil {
		return err
	}
	if string(b) != "OK" {
		return errors.New("bad response")
	}
	return nil
}

//EnvAppDelete deletes application environment variable
func EnvAppDelete(ctx *Context, id int64) error {
	h := authHeader(ctx.Config.AuthToken)
	s := fmt.Sprintf("environment_variable(%d)", id)
	uri := ctx.Config.APIEndpoint(s)
	b, err := doJSON(ctx, "DELETE", uri, h, nil, nil)
	if err != nil {
		return err
	}
	if string(b) != "OK" {
		return errors.New("bad response")
	}
	return nil
}
