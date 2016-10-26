package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/gernest/resingo"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

//ResinDir the data directory for resin
var ResinDir string

//ResinURL the resin api endpoint
var ResinURL string

func init() {
	ResinDir = os.Getenv("RESINRC_DATA_DIRECTORY")
	if ResinDir == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		ResinDir = filepath.Join(usr.HomeDir, ".resin")
		_ = os.MkdirAll(ResinDir, 0600)
	}
	ResinURL = os.Getenv("RESINRC_RESIN_URL")
	if ResinURL == "" {
		ResinURL = "https://api.resin.io"
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "resingo"
	app.Usage = "manage resin devices on the commandline"
	app.Commands = []cli.Command{
		{
			Name:   "login",
			Usage:  "logs into resin service",
			Action: login,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "username",
					Usage: "resin account username",
				},
				cli.StringFlag{
					Name:  "password",
					Usage: "resin account password",
				},
			},
		},
		{
			Name:   "up",
			Usage:  "uses docker-compose.yml to start your services",
			Action: compose,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "projectName",
					Usage: "application name",
				},
				cli.StringFlag{
					Name:  "host",
					Usage: "url to the device",
				},
				cli.StringFlag{
					Name:  "ip",
					Usage: "url to the device",
				},
			},
		},
		{
			Name:   "devices",
			Usage:  "list all devices",
			Action: devices,
		},
	}
	_ = app.Run(os.Args)
}

func login(c *cli.Context) error {
	tk, err := readToken()
	if err != nil {
		if os.IsNotExist(err) {
			return loginWithCredentials(c)
		}
		fmt.Println("HERE")
		return err
	}
	if !resingo.ValidToken(tk) {
		return loginWithCredentials(c)
	}
	fmt.Println("login successful here")
	return nil
}
func loginWithCredentials(c *cli.Context) error {
	uname := c.String("username")
	password := c.String("username")
	config := &resingo.Config{
		Username:      uname,
		Password:      password,
		ResinEndpoint: ResinURL,
	}
	client := &http.Client{}
	ctx := &resingo.Context{
		Client: client,
		Config: config,
	}
	err := resingo.Login(ctx, resingo.Credentials)
	if err != nil {
		return err
	}
	err = writeToken(ctx.Config.AuthToken)
	if err != nil {
		return err
	}
	fmt.Println("login successful")
	return nil
}

func readToken() (string, error) {
	b, err := ioutil.ReadFile(filepath.Join(ResinDir, "token"))
	if err != nil {
		return "", err
	}
	return string(b), err
}

func writeToken(tok string) error {
	return ioutil.WriteFile(filepath.Join(ResinDir, "token"), []byte(tok), 0600)
}

func devices(c *cli.Context) error {
	ctx, err := getCOntext()
	if err != nil {
		return err
	}
	dev, err := resingo.DevGetAll(ctx)
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name"})
	for _, d := range dev {
		table.Append([]string{fmt.Sprint(d.ID), d.Name})
	}
	table.Render()
	return nil
}

func getCOntext() (*resingo.Context, error) {
	tk, err := readToken()
	if err != nil {
		return nil, err
	}
	if !resingo.ValidToken(tk) {
		return nil, errors.New("not logged in " + tk)
	}
	return &resingo.Context{
		Client: &http.Client{},
		Config: &resingo.Config{APIKey: tk, ResinEndpoint: ResinURL},
	}, nil
}

func compose(ctx *cli.Context) error {
	cfg := &resingo.Compose{
		Projectname:  ctx.String("projectName"),
		Host:         ctx.String("host"),
		ComposeFiles: ctx.Args(),
		DeviceIP:     ctx.String("ip"),
	}
	return resingo.Up(cfg)
}
