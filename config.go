package resingo

import "encoding/json"

//ResinConfig resin configuration
type ResinConfig struct {
	SocialProviders    []string `json:"supportedSocialProviders"`
	SignupCodeRequired bool     `json:"signupCodeRequired"`
	MixPanelToken      string   `json:"mixpanelToken"`
	KeenProjectID      string   `json:"keenProjectId"`
	KeenReadKey        string   `json:"keenReadKey"`
	DeviceURLBase      string   `json:"deviceUrlsBase"`
	GitServerURL       string   `json:"gitServerUrl"`
	ImageMakerURL      string   `json:"imgMakerUrl"`
	AdminURL           string   `json:"adminUrl"`
	DebugEnabled       bool     `json:"debugEnabled"`
	PubNub             struct {
		PubKey string `json:"publish_key"`
		SubKey string `json:"subscribe_key"`
	} `json:"pubnub"`
	GoogleAnalytics struct {
		ID   string `kson:"id"`
		Site string `json:"site"`
	} `json:"ga"`

	DeviceTypes []struct {
		Slug              string   `json:"slug"`
		Version           int      `json:"version"`
		Aliases           []string `json:"aliases"`
		Name              string   `json:"name"`
		Arch              string   `json:"arch"`
		State             string   `json:"state"`
		StateInstructions struct {
			PostProvisioning []string `json:"postProvisioning"`
		} `json:"stateInstructions"`
		//Instructions []string `json:"instructions"`
		SupportsBlink bool `json:"supportsBlink"`
		Yocto         struct {
			Machine       string `json:"machine"`
			Image         string `json:"image"`
			FSType        string `json:"fstype"`
			Version       string `json:"version"`
			DeployArtfact string `json:"deployArtfact"`
			Compressed    bool   `json:"compressed"`
		} `json:"yocto"`
		Options []struct {
			IsGroup bool   `json:"isGroup"`
			Name    string `json:"name"`
			Message string `json:"message"`
			Options []struct {
				Name    string   `json:"name"`
				Message string   `json:"message"`
				Type    string   `json:"type"`
				CHoices []string `json:"choices"`
			} `json:"options"`
		} `json:"options"`
		Configuration struct {
			Config struct {
				Partition struct {
					Primary int `json:"primary"`
				} `json:"partition"`
			} `json:"config"`
		} `json:"configuration"`
		Initialization struct {
			Options []struct {
				Name    string `json:"name"`
				Message string `json:"message"`
				Type    string `json:"type"`
			} `json:"options"`
			Operations []struct {
				Command string `json:"command"`
			} `json:"operations"`
		} `json:"initialization"`
		BuildID string `json:"buildId"`
	} `json:"deviceTypes"`
}

//ConfigGetAll return resin congiguration
func ConfigGetAll(ctx *Context) (*ResinConfig, error) {
	h := authHeader(ctx.Config.AuthToken)
	//uri := ctx.Config.APIEndpoint("config")
	uri := apiEndpoint + "/config"
	b, err := doJSON(ctx, "GET", uri, h, nil, nil)
	if err != nil {
		return nil, err
	}
	cfg := &ResinConfig{}
	err = json.Unmarshal(b, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
