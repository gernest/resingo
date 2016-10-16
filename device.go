package resingo

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/guregu/null"
)

//Device represent the information about a resin device
type Device struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	WebAccessible bool   `json:"is-web_accessible"`
	Type          string `json:"device_type"`
	Application   struct {
		ID       int64 `json:"__id"`
		Metadata struct {
			URI string `json:"uri"`
		} `json:"__deferred"`
	} `json:"application"`
	UUID                  string    `json:"uuid"`
	User                  User      `json:"user"`
	Actor                 int64     `json:"actor"`
	IsOnline              bool      `json:"is_online"`
	Commit                string    `json:"commit"`
	Status                string    `json:"status"`
	LastConnectivityEvent null.Time `json:"last_connectivity_event"`
	IP                    string    `json:"ip_address"`
	VPNAddr               string    `json:"vpn_address"`
	PublicAddr            string    `json:"public_address"`
	SuprevisorVersion     string    `json:"supervisor_version"`
	Note                  string    `json:"note"`
	OsVersion             string    `json:"os_version"`
	Location              string    `json:"location"`
	Longitude             string    `json:"longitude"`
	Latitude              string    `json:"latitude"`
}

//DevGetAll returns all devices that belong to the user who authorized the
//context ctx.
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

//GenerateUUID generates uuid suitable for resin devices
func GenerateUUID() (string, error) {
	src := make([]byte, 31)
	_, err := rand.Read(src)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(src), nil
}

//DevRegister registers the device with uuid to the the application with name
//appName.
func DevRegister(ctx *Context, appName, uuid string) (*Device, error) {
	app, err := AppGetByName(ctx, appName)
	if err != nil {
		return nil, err
	}
	appKey, err := AppGetAPIKey(ctx, appName)
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

//DevGetByUUID returns the device with the given uuid.
func DevGetByUUID(ctx *Context, uuid string) (*Device, error) {
	h := authHeader(ctx.Config.AuthToken)
	uri := ctx.Config.APIEndpoint("device")
	params := make(url.Values)
	params.Set("filter", "uuid")
	params.Set("eq", uuid)
	b, err := doJSON(ctx, "GET", uri, h, params, nil)
	if err != nil {
		return nil, err
	}
	var devRes = struct {
		D []*Device `json:"d"`
	}{}
	//fmt.Println(string(b))
	err = json.Unmarshal(b, &devRes)
	if err != nil {
		return nil, err
	}
	if len(devRes.D) > 0 {
		return devRes.D[0], nil
	}
	return nil, errors.New("device not found")
}

//DevGetByName returns the device with the given name
func DevGetByName(ctx *Context, name string) (*Device, error) {
	h := authHeader(ctx.Config.AuthToken)
	uri := ctx.Config.APIEndpoint("device")
	params := make(url.Values)
	params.Set("filter", "name")
	params.Set("eq", name)
	b, err := doJSON(ctx, "GET", uri, h, params, nil)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(b))
	var devRes = struct {
		D []*Device `json:"d"`
	}{}
	err = json.Unmarshal(b, &devRes)
	if err != nil {
		return nil, err
	}
	if len(devRes.D) > 0 {
		//fmt.Println(*devRes.D[0])
		return devRes.D[0], nil
	}
	return nil, errors.New("device not found")
}

//DevIsOnline return true if the device with uuid is online and false otherwise.
//Any errors encountered is returned too.
func DevIsOnline(ctx *Context, uuid string) (bool, error) {
	dev, err := DevGetByUUID(ctx, uuid)
	if err != nil {
		return false, err
	}
	return dev.IsOnline, nil
}

//DevGetAllByApp returns all devices that are registered to the application with
//the given appName.
func DevGetAllByApp(ctx *Context, appName string) ([]*Device, error) {
	app, err := AppGetByName(ctx, appName)
	if err != nil {
		return nil, err
	}
	h := authHeader(ctx.Config.AuthToken)
	uri := ctx.Config.APIEndpoint("device")
	params := make(url.Values)
	params.Set("filter", "application")
	params.Set("eq", fmt.Sprint(app.ID))
	b, err := doJSON(ctx, "GET", uri, h, params, nil)
	if err != nil {
		return nil, err
	}
	var devRes = struct {
		D []*Device `json:"d"`
	}{}
	//fmt.Println(string(b))
	err = json.Unmarshal(b, &devRes)
	if err != nil {
		return nil, err
	}
	return devRes.D, nil
}

//DevRename renames the device with uuid to nwName
func DevRename(ctx *Context, uuid, newName string) error {
	_, err := DevGetByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	h := authHeader(ctx.Config.AuthToken)
	uri := ctx.Config.APIEndpoint("device")
	params := make(url.Values)
	params.Set("filter", "uuid")
	params.Set("eq", uuid)
	data := make(map[string]interface{})
	data["name"] = newName
	body, err := marhsalReader(data)
	if err != nil {
		return err
	}
	b, err := doJSON(ctx, "PATCH", uri, h, params, body)
	if err != nil {
		return err
	}
	if string(b) != "OK" {
		return errors.New("bad response")
	}
	return nil
}

//DevGetApp returns the application in which the device belongs to. This
//function is convenient only when you are interested on other information about
//the application.
//
// If your intention is only to retrieve the applicayion id, then just use this
// instead.
//
//	dev,err:=DevGetByUUID(ctx,<uuid goes here>)
//	if err!=nil{
//		//handle error error
//	}
//	// you can now access the application id like this
//	fmt.Println(dev.Application.ID
func DevGetApp(ctx *Context, uuid string) (*Application, error) {
	dev, err := DevGetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return AppGetByID(ctx, dev.Application.ID)
}
