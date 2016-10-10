package resingo

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"time"
)

type Device struct {
	ID            int64       `json:"id"`
	Name          string      `json:"name"`
	WebAccessible bool        `json:"is-web_accessible"`
	Type          string      `json:"device_type"`
	Application   Application `json:"application"`
	UUID          string      `json:"uuid"`
	User          User        `json:"user"`
}

func DevGetAll(ctx *Context) ([]*Device, error) {
	h := authHeader(ctx.Config.AuthToken)
	uri := ctx.Config.APIEndpoint("device")
	b, err := doJSON(ctx, "GET", uri, h, nil, nil)
	if err != nil {
		return nil, err
	}
	var devRes = struct {
		D []*Device `json:"d"`
	}{}
	err = json.Unmarshal(b, &devRes)
	if err != nil {
		return nil, err
	}
	return devRes.D, nil
}

func GenerateUUID() (string, error) {
	src := make([]byte, 31)
	_, err := rand.Read(src)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(src), nil
}

func DevRegister(ctx *Context, appName, uuid string) (*Device, error) {
	app, err := AppGetByName(ctx, appName)
	if err != nil {
		return nil, err
	}
	appKey, err := AppGetApiKey(ctx, appName)
	if err != nil {
		return nil, err
	}
	h := authHeader(ctx.Config.AuthToken)
	uri := ctx.Config.APIEndpoint("device")
	data := make(map[string]interface{})
	data["user"] = ctx.Config.UserID()
	data["device_type"] = app.DeviceType
	data["application"] = app.ID
	data["registered_at"] = time.Now().Unix()
	data["uuid"] = uuid
	data["apikey"] = appKey
	body, err := marhsalReader(data)
	if err != nil {
		return nil, err
	}
	b, err := doJSON(ctx, "POST", uri, h, nil, body)
	//fmt.Println(string(b))
	if err != nil {
		return nil, err
	}
	rst := &Device{}
	err = json.Unmarshal(b, rst)
	if err != nil {
		return nil, err
	}
	return rst, nil

}
