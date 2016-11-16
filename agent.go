package resingo

import (
	"encoding/json"
	"errors"
)

//AgentReboot reboots the device
func AgentReboot(ctx *Context, devID, appID int64, force bool) error {
	h := authHeader(ctx.Config.AuthToken)
	uri := apiEndpoint + "/supervisor/v1/reboot"
	data := make(map[string]interface{})
	data["deviceId"] = devID
	data["appId"] = appID
	data["force"] = force
	body, err := marhsalReader(data)
	if err != nil {
		return err
	}
	b, err := doJSON(ctx, "POST", uri, h, nil, body)
	if err != nil {
		return err
	}
	var res = struct {
		Data  string
		Error string
	}{}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return err
	}
	if res.Data != "OK" {
		return errors.New("bad response :" + res.Error)
	}
	return nil

}
