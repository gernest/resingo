package resingo

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

//Key is a user public key on resin
type Key struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	PublicKey string `json:"public_key"`
	User      struct {
		ID       int64 `json:"__id"`
		Deferred struct {
			URI string `json:"uri"`
		} `json:"__deferred"`
	} `json:"user"`
	Metadata struct {
		URI  string `json:"uri"`
		Type string `json:"type"`
	} `json:"__metadata"`
	CreatedAt time.Time `json:"created_at"`
}

//KeyGetAll retrives all key for the user who authenticated ctx.
func KeyGetAll(ctx *Context) ([]*Key, error) {
	h := authHeader(ctx.Config.AuthToken)
	uri := ctx.Config.APIEndpoint("user__has__public_key")
	b, err := doJSON(ctx, "GET", uri, h, nil, nil)
	if err != nil {
		return nil, err
	}
	res := struct {
		D []*Key `json:"d"`
	}{}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}
	return res.D, nil
}

//KeyGetByID retrives public key with the given id
func KeyGetByID(ctx *Context, id int64) (*Key, error) {
	h := authHeader(ctx.Config.AuthToken)
	s := fmt.Sprintf("user__has__public_key(%d)", id)
	uri := ctx.Config.APIEndpoint(s)
	b, err := doJSON(ctx, "GET", uri, h, nil, nil)
	if err != nil {
		return nil, err
	}
	res := struct {
		D []*Key `json:"d"`
	}{}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}
	if len(res.D) > 0 {
		return res.D[0], nil
	}
	return nil, errors.New("key not found")
}

//KeyCreate creates a public key for the user with given userID
func KeyCreate(ctx *Context, userID int64, key, title string) (*Key, error) {
	h := authHeader(ctx.Config.AuthToken)
	uri := ctx.Config.APIEndpoint("user__has__public_key")
	data := make(map[string]interface{})
	data["user"] = userID
	data["public_key"] = key
	data["title"] = title
	body, err := marhsalReader(data)
	if err != nil {
		return nil, err
		//return nil
	}
	b, err := doJSON(ctx, "POST", uri, h, nil, body)
	if err != nil {
		return nil, err
		//return nil
	}
	//fmt.Println(string(b))
	e := &Key{}
	err = json.Unmarshal(b, e)
	if err != nil {
		return nil, err
	}
	return e, nil
	//return nil
}

//KeyRemove removes the public key with the given id
func KeyRemove(ctx *Context, id int64) error {
	h := authHeader(ctx.Config.AuthToken)
	s := fmt.Sprintf("user__has__public_key(%d)", id)
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
