package resingo

import "encoding/json"

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
