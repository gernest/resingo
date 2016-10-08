package resingo

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

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

type ApiVersion int

const (
	VersionOne ApiVersion = iota
	VersionTwo
	VersionThree
)

func (v ApiVersion) String() string {
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

var ErrUnkownAuthType = errors.New("resingo: unknown authentication type")
var ErrMissingCredentials = errors.New("resingo: missing credentials( username or password)")
var ErrBadToken = errors.New("resingo: bad session token")

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
	Post(url string, bodyTyp string, body io.Reader) (*http.Response, error)
}

type Context struct {
	Client HTTPClient
	Config *Config
}

//Config is the configuration object for the Client
type Config struct {
	AuthToken     string
	Username      string
	Password      string
	ApiKey        string
	tokenClain    *TokenClain
	ResinEndpoint string
	ResinVersion  ApiVersion
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

func apiURL(base string, version ApiVersion, endpoint string) string {
	return fmt.Sprintf("%s/%s/%s", base, version, endpoint)
}

//ApiEndpoint returns a url that points to the given endpoint. This adds the
//resin.io api host and version.
func (c *Config) APIEndpoint(endpoint string) string {
	return apiURL(c.ResinEndpoint, c.ResinVersion, endpoint)
}

func (c *Config) IsValidToken(tok string) bool {
	return true
}

func authHeader(token string) http.Header {
	h := make(http.Header)
	h.Add("Authorization", "Bearer "+token)
	return h
}

//ParseToken parses and saves the token into the *Config instance.
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
//you want to save the token in the client. This function doens not save the
//authentication token and user detals.
func Authenticate(ctx *Context, typ AuthType, authToken ...string) (string, error) {
	loginURL := apiEndpoint + "/login_"
	switch typ {
	case Credentials:
		// authenticate using credentials
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
		defer res.Body.Close()
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

func Login(ctx *Context, authTyp AuthType, authToken ...string) error {
	tok, err := Authenticate(ctx, authTyp, authToken...)
	if err != nil {
		return err
	}
	if ctx.Config.IsValidToken(tok) {
		ctx.Config.SaveToken(tok)
		return nil
	}
	return errors.New("resingo: Failed to login")
}

func Encode(q url.Values) string {
	if q == nil {
		return ""
	}
	if filter := q.Get("filter"); filter != "" {
		eq := q.Get("eq")
		f := "$filter=" + filter + "%20eq%20" + fmt.Sprintf("'%s'", eq)
		q.Del("eq")
		q.Del("filter")
		s := q.Encode()
		if s != "" {
			return fmt.Sprintf("%s&%s", f, s)
		}
		return f
	}
	return q.Encode()
}
