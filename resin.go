package resingo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//AuthType is the authentication type that is used to authenticate with a
//resin.io api.
type AuthType int

// supported authentication types
const (
	Credentials AuthType = iota
	AuthToken
)

const (
	pineEndpoint           = "https://api.resin.io/ewa"
	apiEndpoint            = "https://api.resin.io"
	tokenRefreshInterval   = 3600000
	imageCacheTime         = 86400000
	applicationEndpoint    = "/application"
	deviceEndpoint         = "/device"
	keysEndpoint           = "/user__has__public_key"
	applicationEnvEndpoint = "/environment_variable"
	deviceEnvEndpoint      = "/device_environment_variable"
)

//APIVersion is the version of resin API
type APIVersion int

// supported resin API versions
const (
	VersionOne APIVersion = iota
	VersionTwo
	VersionThree
)

func (v APIVersion) String() string {
	switch v {
	case VersionOne:
		return "v1"
	case VersionTwo:
		return "v2"
	case VersionThree:
		return "v3"
	}
	return ""
}

//ErrUnkownAuthType error returned when the type of authentication is not
//supported.
var ErrUnkownAuthType = errors.New("resingo: unknown authentication type")

//ErrMissingCredentials error returned when either username or password is
//missing
var ErrMissingCredentials = errors.New("resingo: missing credentials( username or password)")

//ErrBadToken error returned when the resin session token is bad.
var ErrBadToken = errors.New("resingo: bad session token")

//HTTPClient is an interface for a http clinet that is used to communicate with
//the resin API
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
	Post(url string, bodyTyp string, body io.Reader) (*http.Response, error)
}

//Context holds information necessary to make a call to the resin API
type Context struct {
	Client HTTPClient
	Config *Config
}

//Config is the configuration object for the Client
type Config struct {
	AuthToken     string
	Username      string
	Password      string
	APIKey        string
	tokenClain    *TokenClain
	ResinEndpoint string
	ResinVersion  APIVersion
}

//TokenClain are the values that are encoded into a session token from resin.io.
//
// It embeds jst.StandardClaims, so as to help with Verification of expired
// data. Resin doens't do claim verification :(.
type TokenClain struct {
	Username string `json:"username"`
	UserID   int64  `json:"id"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

// formats a proper url forthe API call. The format is
// /<base_url>/<api_version>/<api_endpoin. The endpoint can be en empty string.
//
//This assumes that the endpoint doesn't start with /
// TODO: handle endpoint that starts with / and base url that ends with /
func apiURL(base string, version APIVersion, endpoint string) string {
	return fmt.Sprintf("%s/%s/%s", base, version, endpoint)
}

//APIEndpoint returns a url that points to the given endpoint. This adds the
//resin.io api host and version.
func (c *Config) APIEndpoint(endpoint string) string {
	return apiURL(c.ResinEndpoint, c.ResinVersion, endpoint)
}

//IsValidToken return true if the token tok is a valid resin session token.
//
// This method ecodes the token. A token that can't be doced is bad token. Any
// token that has expired is also a bad token.
func (c *Config) IsValidToken(tok string) bool {
	tk, err := ParseToken(tok)
	if err != nil {
		return false
	}
	return tk.StandardClaims.ExpiresAt > time.Now().Unix()
}

//ValidToken return true if tok is avalid token
func ValidToken(tok string) bool {
	tk, err := ParseToken(tok)
	if err != nil {
		return false
	}
	return tk.StandardClaims.ExpiresAt > time.Now().Unix()
}

//UserID returns the user id.
func (c *Config) UserID() int64 {
	return c.tokenClain.UserID
}

func authHeader(token string) http.Header {
	h := make(http.Header)
	h.Add("Authorization", "Bearer "+token)
	return h
}

//SaveToken saves token, to the current Configuration object.
func (c *Config) SaveToken(tok string) error {
	tk, err := ParseToken(tok)
	if err != nil {
		return err
	}
	c.tokenClain = tk
	c.AuthToken = tok
	return nil
}

//ParseToken parses the given token and extracts the claims emcode into it. This
//function uses JWT method to parse the token, with verification of claims
//turned off.
func ParseToken(tok string) (*TokenClain, error) {
	p := jwt.Parser{
		SkipClaimsValidation: true,
	}
	tk, _ := p.ParseWithClaims(tok, &TokenClain{}, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	claims, ok := tk.Claims.(*TokenClain)
	if ok {
		return claims, nil
	}
	return nil, ErrBadToken
}

//Authenticate authenticates the client and returns the Auth token. See Login if
//you want to save the token in the client. This function does not save the
//authentication token and user detals.
func Authenticate(ctx *Context, typ AuthType, authToken ...string) (string, error) {
	loginURL := apiEndpoint + "/login_"
	switch typ {
	case Credentials:
		// Absence of either username or password result in missing creadentials
		// error.
		if ctx.Config.Username == "" || ctx.Config.Password == "" {
			return "", ErrMissingCredentials
		}
		form := url.Values{}
		form.Add("username", ctx.Config.Username)
		form.Add("password", ctx.Config.Password)
		res, err := ctx.Client.Post(loginURL,
			"application/x-www-form-urlencoded",
			strings.NewReader(form.Encode()))
		if err != nil {
			return "", err
		}
		defer func() {
			_ = res.Body.Close()
		}()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		return string(data), nil
	case AuthToken:
		if len(authToken) > 0 {
			tk := authToken[0]
			if ctx.Config.IsValidToken(tk) {
				return tk, nil
			}
			return "", ErrBadToken
		}
		return "", errors.New("resingo: Failed to authenticate missing authToken")
	}
	return "", ErrUnkownAuthType
}

//Login authenticates the contextand stores the session token. This function
//checks the validity of the session token before saving it.
//
// The call to ctx.IsLoged() should return true if the returned error is nil.
func Login(ctx *Context, authTyp AuthType, authToken ...string) error {
	tok, err := Authenticate(ctx, authTyp, authToken...)
	if err != nil {
		return err
	}
	if ctx.Config.IsValidToken(tok) {
		return ctx.Config.SaveToken(tok)
	}
	return errors.New("resingo: Failed to login")
}

//Encode encode properly the request params for use with resin API.
//
// Encode tartegts the filter param, which for some reasom(based on OData) is
// supposed to be $filter and not filter. The value specified by the eq param
// key is combined with the value from the fileter key to produce the $filter
// value string.
//
// Any other url params are encoded by the default encoder from
// url.Values.Encoder.
//TODO: check a better way to encode OData url params.
func Encode(q url.Values) string {
	if q == nil {
		return ""
	}
	var buf bytes.Buffer
	var keys []string
	for k := range q {
		keys = append(keys, k)
	}
	for _, k := range keys {
		switch k {
		case "filter":
			if buf.Len() != 0 {
				_, _ = buf.WriteRune('&')
			}
			v := q.Get("filter")
			_, _ = buf.WriteString("$filter=" + v)
			for _, fk := range keys {
				switch fk {
				case "eq":
					fv := "%20" + fk + "%20" + quote(q.Get(fk))
					_, _ = buf.WriteString(fv)
					q.Del(fk)
				}
			}
			q.Del(k)
		case "expand":
			if buf.Len() != 0 {
				_, _ = buf.WriteRune('&')
			}
			v := q.Get("expand")
			_, _ = buf.WriteString("$expand=" + v)
			q.Del(k)
		}
	}
	e := q.Encode()
	if e != "" {
		if buf.Len() != 0 {
			_, _ = buf.WriteRune('&')
		}
		_, _ = buf.WriteString(e)
	}
	return buf.String()
}

func quote(v string) string {
	ok, _ := strconv.ParseBool(v)
	if ok {
		return v
	}
	_, err := strconv.Atoi(v)
	if err == nil {
		return v
	}
	_, err = strconv.ParseFloat(v, 64)
	if err == nil {
		return v
	}
	return "'" + v + "'"
}
